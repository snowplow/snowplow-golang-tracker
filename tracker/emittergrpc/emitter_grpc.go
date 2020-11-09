//
// Copyright (c) 2016-2020 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package emittergrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	sgrpc "github.com/snowplow/snowplow-golang-tracker/v2/grpc"
	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

// EmitterGRPC contains the attributes of the GRPC Snowplow Emitter
type EmitterGRPC struct {
	CollectorURI string
	SendLimit    int
	StreamLimit  int
	DbName       string
	Storage      gt.Storage
	SendChannel  chan bool
	Callback     func(successCount []gt.CallbackResult, failureCount []gt.CallbackResult)
	GrpcConn     *grpc.ClientConn
	GrpcClient   sgrpc.CollectorServiceClient
	TLSFilePath  string
}

// InitEmitter creates a new GRPC Emitter object which handles
// storing and sending Snowplow Events.
func InitEmitter(options ...func(*EmitterGRPC)) *EmitterGRPC {
	e := &EmitterGRPC{}

	// Set Defaults
	e.SendLimit = gt.DEFAULT_SEND_LIMIT
	e.StreamLimit = 500
	e.DbName = gt.DEFAULT_DB_NAME
	e.TLSFilePath = ""

	// Option parameters
	for _, op := range options {
		op(e)
	}

	// Check collector URI is not empty
	if e.CollectorURI == "" {
		panic("FATAL: CollectorURI cannot be empty.")
	}

	// Setup default event storage
	if e.Storage == nil {
		e.Storage = *gt.InitStorageSQLite3(e.DbName)
	}

	return e
}

// --- Require

// RequireCollectorURI sets the Emitters collector URI.
func RequireCollectorURI(collectorURI string) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.CollectorURI = collectorURI }
}

// --- Option

// OptionSendLimit sets the send limit for the emitter.
func OptionSendLimit(sendLimit int) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.SendLimit = sendLimit }
}

// OptionStreamLimit sets the amount of events to be sent per stream.
func OptionStreamLimit(streamLimit int) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.StreamLimit = streamLimit }
}

// OptionTLSFilePath configures the GRPC clients PEM keys to use in authenticating
// a secure connection with the server.
//
// Note: If this value is empty the client defaults to insecure
func OptionTLSFilePath(tlsFilePath string) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.TLSFilePath = tlsFilePath }
}

// OptionDbName overrides the default name of the storage database.
func OptionDbName(dbName string) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.DbName = dbName }
}

// OptionStorage sets a custom event Storage target which implements the Storage interface
//
// Note: If this option is used OptionDbName will be ignored
func OptionStorage(storage gt.Storage) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.Storage = storage }
}

// OptionCallback sets a custom callback for the emitter loop.
func OptionCallback(callback func(successCount []gt.CallbackResult, failureCount []gt.CallbackResult)) func(e *EmitterGRPC) {
	return func(e *EmitterGRPC) { e.Callback = callback }
}

// --- GRPC Handlers

// dial establishes the connection with the GRPC server.
func dial(collectorURI string, tlsFilePath string) (*grpc.ClientConn, error) {
	if tlsFilePath == "" {
		return grpc.Dial(collectorURI, grpc.WithInsecure())
	}

	creds, err := credentials.NewClientTLSFromFile(tlsFilePath, "")
	if err != nil {
		return nil, err
	}
	return grpc.Dial(collectorURI, grpc.WithTransportCredentials(creds))
}

// --- Event Handlers

// Add will push an event to the database and will then initiate a sending loop.
func (e *EmitterGRPC) Add(payload gt.Payload) {
	e.Storage.AddEventRow(payload)
	e.start()
}

// Flush will attempt to start the send loop regardless of an event coming in.
func (e *EmitterGRPC) Flush() {
	e.start()
}

// Stop waits for the send channel to have a value and then resets it to nil.
func (e *EmitterGRPC) Stop() {
	<-e.SendChannel
	e.SendChannel = nil

	if e.GrpcConn != nil {
		e.GrpcConn.Close()
	}
}

// start will begin the sending loop.
func (e *EmitterGRPC) start() {
	if e.SendChannel == nil || !e.IsSending() {
		if e.GrpcConn != nil {
			e.GrpcConn.Close()
		}
		conn, err := dial(e.CollectorURI, e.TLSFilePath)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: could not establish connection with GRPC server: %v", err))
			return
		}
		e.GrpcConn = conn
		e.GrpcClient = sgrpc.NewCollectorServiceClient(e.GrpcConn)

		e.SendChannel = make(chan bool, 1)
		go func() {
			var done bool
			defer func() {
				e.SendChannel <- done
			}()

			for {
				eventRows := e.Storage.GetEventRowsWithinRange(e.SendLimit)

				// If there are no events in the database exit
				if len(eventRows) == 0 {
					break
				}
				results := e.doSend(eventRows)

				// Process results
				ids := []int{}
				successes := []gt.CallbackResult{}
				failures := []gt.CallbackResult{}

				for _, res := range results {

					count := len(res.Ids)
					status := res.Status

					if status >= 200 && status < 400 {
						ids = append(ids, res.Ids...)
						successes = append(successes, gt.CallbackResult{count, status})
					} else {
						failures = append(failures, gt.CallbackResult{count, status})
					}
				}

				if e.Callback != nil {
					e.Callback(successes, failures)
				}

				// If all the events failed to be sent exit
				if len(successes) == 0 && len(failures) > 0 {
					break
				}

				e.Storage.DeleteEventRows(ids)
			}
			done = true
		}()
	}
}

// doSend will send all of the eventsRows it is given.
func (e *EmitterGRPC) doSend(eventRows []gt.EventRow) []gt.SendResult {
	futures := []<-chan gt.SendResult{}

	ids := []int{}
	payloads := []gt.Payload{}

	for _, val := range eventRows {
		if len(payloads) >= e.StreamLimit {
			futures = append(futures, e.sendGrpcRequest(e.GrpcClient, ids, payloads))
			ids = []int{val.Id}
			payloads = []gt.Payload{val.Event}
		} else {
			ids = append(ids, val.Id)
			payloads = append(payloads, val.Event)
		}
	}
	if len(payloads) > 0 {
		futures = append(futures, e.sendGrpcRequest(e.GrpcClient, ids, payloads))
	}

	// Wait for all Futures to complete
	results := []gt.SendResult{}
	for _, future := range futures {
		results = append(results, <-future)
	}

	return results
}

// SendPostRequest sends an array of Payloads together to the collector endpoint via POST.
func (e *EmitterGRPC) sendGrpcRequest(client sgrpc.CollectorServiceClient, ids []int, body []gt.Payload) <-chan gt.SendResult {
	c := make(chan gt.SendResult, 1)
	go func() {
		var result gt.SendResult
		defer func() {
			c <- result
		}()

		status := -1

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		stream, err := client.StreamTrackPayload(ctx)
		if err != nil {
			log.Println(fmt.Sprintf("%v.StreamTrackPayload(_) = _, %v", client, err))
			result = gt.SendResult{Ids: ids, Status: status}
			return
		}

		for _, val := range body {
			val.Add(gt.SENT_TIMESTAMP, gt.NewString(gt.GetTimestampString()))
			tpReq, _ := payloadToTrackPayloadRequest(val)
			err := stream.Send(tpReq)
			if err != nil {
				log.Println(fmt.Sprintf("%v.Send(%v) = %v", stream, tpReq, err))
				result = gt.SendResult{Ids: ids, Status: status}
				return
			}
		}

		resp, err := stream.CloseAndRecv()
		if err != nil {
			log.Println(fmt.Sprintf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil))
			result = gt.SendResult{Ids: ids, Status: status}
			return
		}

		if resp.GetSuccess() {
			result = gt.SendResult{Ids: ids, Status: 200}
		} else {
			result = gt.SendResult{Ids: ids, Status: 503}
		}
	}()
	return c
}

// --- Helpers

// IsSending checks whether the send channel has finished.
func (e EmitterGRPC) IsSending() bool {
	return len(e.SendChannel) == 0
}

// payloadToTrackEventRequest converts a Payload into a gRPC TrackPayloadRequest
//
// TODO: Remove double serialization due to field names being renamed
//       https://developers.google.com/protocol-buffers/docs/reference/go-generated#fields
func payloadToTrackPayloadRequest(payload gt.Payload) (*sgrpc.TrackPayloadRequest, error) {
	ter := &sgrpc.TrackPayloadRequest{}
	d := json.NewDecoder(strings.NewReader(payload.String()))
	d.UseNumber()
	err := d.Decode(&ter)
	if err != nil {
		return nil, err
	}
	return ter, nil
}

// --- Getters & Setters

// GetSendChannel returns the send channel
func (e EmitterGRPC) GetSendChannel() chan bool {
	return e.SendChannel
}

// GetStorage returns the send channel
func (e EmitterGRPC) GetStorage() gt.Storage {
	return e.Storage
}
