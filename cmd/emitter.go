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

package cmd

import (
	"log"
	"time"

	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
	egrpc "github.com/snowplow/snowplow-golang-tracker/v2/tracker/emittergrpc"
	"github.com/spf13/cobra"
)

var emitterType string
var sendLimit int
var events int
var streamLimit int
var tlsFilePath string
var requestType string
var protocol string

var emitterCmd = &cobra.Command{
	Use:   "emitter",
	Short: "Stress test different emitter configurations",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize objects
		storage := gt.InitStorageMemory()
		emitter := buildEmitter(storage, emitterType)
		subject := gt.InitSubject()
		subject.SetUserId("new-user-id")
		subject.SetScreenResolution(1920, 1080)
		subject.SetViewPort(1080, 1080)
		subject.SetColorDepth(1080)
		subject.SetTimeZone("ACT")
		subject.SetLanguage("EN")
		subject.SetUseragent("golang-tracker")
		subject.SetDomainUserId("5339cd71-d0c4-40b9-88e3-4f3db87bcfad")
		subject.SetNetworkUserId("c0818e23-8e07-40ea-8442-192f79ab82c6")
		tracker := gt.InitTracker(
			gt.RequireEmitter(emitter),
			gt.OptionSubject(subject),
		)

		data := map[string]interface{}{
			"programmingLanguage": "GOLANG",
			"message":             "Only Kiddin'",
			"lineNumber":          12,
		}
		sdj := gt.InitSelfDescribingJson("iglu:com.snowplowanalytics.snowplow/application_error/jsonschema/1-0-2", data)

		contextArray := []gt.SelfDescribingJson{
			*gt.InitSelfDescribingJson(
				"iglu:com.snowplowanalytics.snowplow/ad_click/jsonschema/1-0-0",
				map[string]interface{}{
					"targetUrl":  "https://snowplowanalytics.com/",
					"campaignId": "12345",
				},
			),
		}

		// Perform test!
		log.Printf("Sending %v events to '%s' now ...", events, collectorURI)
		start := time.Now()

		for n := 0; n < events; n++ {
			tracker.TrackSelfDescribingEvent(gt.SelfDescribingEvent{
				Event:    sdj,
				Contexts: contextArray,
			})
		}

		elapsed := time.Since(start)
		log.Printf("Pushing %v events to the tracker took '%s'; waiting for all events to be sent ...", events, elapsed)
		failedRows := tracker.BlockingFlush(10, 50)

		elapsed = time.Since(start)
		if failedRows == 0 {
			log.Printf("Successfully sent '%v' events to the Collector in '%s'", events, elapsed)
		} else {
			log.Printf("Successfully sent '%v' events and failed to send '%v' to the Collector in '%s'", events-failedRows, failedRows, elapsed)
		}
	},
}

func init() {
	RootCmd.AddCommand(emitterCmd)
	emitterCmd.Flags().StringVarP(&emitterType, "type", "t", "GRPC", "The type of emitter to use (HTTP or GRPC)")
	emitterCmd.Flags().IntVarP(&sendLimit, "sendLimit", "", 500, "The number of events to emit from the database at a time")
	emitterCmd.Flags().IntVarP(&events, "events", "e", 10, "The number of events during this test")
	// GRPC flags
	emitterCmd.Flags().IntVarP(&streamLimit, "streamLimit", "", 100, "The number of events to emit per GRPC stream")
	emitterCmd.Flags().StringVarP(&tlsFilePath, "tlsFilePath", "", "", "The absolute path to the root certificates file for securing the GRPC server connection")
	// HTTP flags
	emitterCmd.Flags().StringVarP(&requestType, "requestType", "r", "POST", "The request type to use for HTTP sending")
	emitterCmd.Flags().StringVarP(&protocol, "protocol", "p", "HTTPS", "The protocol to use for HTTP sending")
}

// --- Helpers

func buildEmitter(s gt.Storage, et string) gt.Emitter {
	if et == "GRPC" {
		return egrpc.InitEmitter(
			egrpc.RequireCollectorURI(collectorURI),
			egrpc.OptionStorage(s),
			egrpc.OptionSendLimit(sendLimit),
			egrpc.OptionStreamLimit(streamLimit),
			egrpc.OptionTLSFilePath(tlsFilePath),
		)
	} else if et == "HTTP" {
		return gt.InitEmitter(
			gt.RequireCollectorUri(collectorURI),
			gt.OptionStorage(s),
			gt.OptionRequestType(requestType),
			gt.OptionProtocol(protocol),
			gt.OptionSendLimit(sendLimit),
		)
	} else {
		log.Fatalf("FATAL: Illegal emitter type supplied.  Expected one of HTTP or GRPC and got '%s'!", et)
		return nil
	}
}
