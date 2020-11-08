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
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	sgrpc "github.com/snowplow/snowplow-golang-tracker/v2/grpc"
	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

func TestEmitterInit(t *testing.T) {
	assert := assert.New(t)
	emitter := InitEmitter(
		RequireCollectorURI("com.acme"),
		OptionSendLimit(1000),
		OptionStreamLimit(1000),
		OptionTLSFilePath("/var/fake/tls.pem"),
		OptionDbName("test.db"),
		OptionCallback(func(g []gt.CallbackResult, b []gt.CallbackResult) {
			log.Println("Successes: " + gt.IntToString(len(g)))
			log.Println("Failures: " + gt.IntToString(len(b)))
		}),
		OptionStorage(*gt.InitStorageMemory()),
	)

	// Assert the option builders
	assert.NotNil(emitter)
	assert.Equal("com.acme", emitter.CollectorURI)
	assert.Equal(1000, emitter.SendLimit)
	assert.Equal(1000, emitter.StreamLimit)
	assert.Equal("test.db", emitter.DbName)
	assert.Equal("/var/fake/tls.pem", emitter.TLSFilePath)
	assert.NotNil(emitter.Storage)
	assert.NotNil(emitter.GetStorage())
	assert.Nil(emitter.SendChannel)
	assert.Nil(emitter.GetSendChannel())
	assert.NotNil(emitter.Callback)
	assert.Nil(emitter.GrpcConn)
	assert.Nil(emitter.GrpcClient)
	assert.NotNil(emitter.Storage)
	assert.Equal("tracker.StorageMemory", reflect.TypeOf(emitter.Storage).String())
	assert.True(emitter.IsSending())

	// Assert defaults
	emitter = InitEmitter(RequireCollectorURI("com.acme"), OptionDbName("test.db"))
	assert.NotNil(emitter)
	assert.Equal("com.acme", emitter.CollectorURI)
	assert.Equal(500, emitter.SendLimit)
	assert.Equal(500, emitter.StreamLimit)
	assert.Equal("test.db", emitter.DbName)
	assert.Equal("", emitter.TLSFilePath)
	assert.NotNil(emitter.Storage)
	assert.NotNil(emitter.GetStorage())
	assert.Nil(emitter.SendChannel)
	assert.Nil(emitter.GetSendChannel())
	assert.Nil(emitter.Callback)
	assert.Nil(emitter.GrpcConn)
	assert.Nil(emitter.GrpcClient)
	assert.NotNil(emitter.Storage)
	assert.Equal("tracker.StorageSQLite3", reflect.TypeOf(emitter.Storage).String())
	assert.True(emitter.IsSending())
}

func TestEmitterNoURI(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		if err := recover(); err != nil {
			assert.Equal("FATAL: CollectorURI cannot be empty.", err)
		}
	}()

	emitter := InitEmitter()
	assert.Nil(emitter)
}

func TestEmitterInterface(t *testing.T) {
	assert := assert.New(t)

	successes := 0
	failures := 0

	emitter := InitEmitter(
		RequireCollectorURI("localhost:50051"),
		OptionStorage(*gt.InitStorageMemory()),
		OptionSendLimit(1),
		OptionStreamLimit(1),
		OptionCallback(func(s []gt.CallbackResult, f []gt.CallbackResult) {
			successes = successes + len(s)
			failures = failures + len(f)
		}),
	)

	payload := *gt.InitPayload()
	payload.Add("e", gt.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))

	for n := 0; n < 5; n++ {
		emitter.Add(payload)
		emitter.Stop()
	}
	emitter.Flush()
	emitter.Stop()

	assert.Equal(0, successes)
	assert.Equal(6, failures)
}

func TestPayloadToTrackPayloadRequest_Valid(t *testing.T) {
	assert := assert.New(t)

	payload := *gt.InitPayload()
	payload.Add("e", gt.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))

	res, err := payloadToTrackPayloadRequest(payload)
	assert.Nil(err)
	assert.NotNil(res)
	assert.Equal("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", res.E)
}

func TestDoSend_InvalidURI_Insecure(t *testing.T) {
	assert := assert.New(t)

	emitter := InitEmitter(
		RequireCollectorURI("localhost:50051"),
		OptionStorage(*gt.InitStorageMemory()),
		OptionSendLimit(1),
		OptionStreamLimit(1),
	)

	conn, _ := dial(emitter.CollectorURI, "")
	defer conn.Close()
	emitter.GrpcConn = conn
	emitter.GrpcClient = sgrpc.NewCollectorServiceClient(emitter.GrpcConn)

	payload := *gt.InitPayload()
	payload.Add("e", gt.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))

	eventRows := []gt.EventRow{{Id: -1, Event: payload}, {Id: -1, Event: payload}, {Id: -1, Event: payload}, {Id: -1, Event: payload}, {Id: -1, Event: payload}}
	results := emitter.doSend(eventRows)

	assert.NotNil(results)
	assert.Equal(5, len(results))
	assert.Equal(-1, results[0].Ids[0])
	assert.Equal(-1, results[0].Status)
}
