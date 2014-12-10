/*
Copyright (c) 2014-2015 Snowplow Analytics Ltd. All rights reserved.

This program is licensed to you under the Apache License Version 2.0,
and you may not use this file except in compliance with the Apache License
Version 2.0. You may obtain a copy of the Apache License Version 2.0 at

    http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing,
software distributed under the Apache License Version 2.0 is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the Apache License Version 2.0 for the specific
language governing permissions and limitations there under.
*/
package snowplowGo

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const (
	DEFAULT_BASE_64  = true
	TRACKER_VERSION  = "golang-0.1.0"

	SCHEMA_VENDOR    = "com.snowplowanalytics.snowplow"
	SCHEMA_FORMAT    = "jsonschema"
)

type Tracker struct {
	Emitter      Emitter
	subject      Subject
	EncodeBase64 bool
}

type JsonSchema struct {
	ContextSchema       string
	UnstructEventSchema string
	ScreenViewSchema    string
}

// TODO(alexanderdean): why are these variables global to the package?

// TODO(alexanderdean): typo, should be schema
var scehma JsonSchema

// TODO(alexanderdean): why does this start with capital 'S'?
var StdNvPairs map[string]string

var s Tracker

// Initializes a new tracker instance with emitter(s) and a subject.
func InitTracker(emitterTracker map[string]string, subject Subject, namespace string, AppId string, EncodeBase64 string) {
	if len(emitterTracker) > 0 {
		s.emitter = emitterTracker
	} else {
		s.emitter = emitterTracker
	}
	s.subject = subject
	if s.EncodeBase64 != nil {
		s.EncodeBase64 = StringToBool(s.EncodeBase64)
	} else {
		s.EncodeBase64 = DEFAULT_BASE_64
	}

	StdNvPairs["tv"] = TRACKER_VERSION
	StdNvPairs["tna"] = namespace
	StdNvPairs["aid"] = AppId

	schema.ContextSchema = "iglu:" + SCHEMA_VENDOR + "/contexts/" + SCHEMA_FORMAT + "/1-0-1"
	schema.UnstructEventSchema = "iglu:" + SCHEMA_VENDOR + "/unstruct_event/" + SCHEMA_FORMAT + "/1-0-0"
	schema.ScreenViewSchema = "iglu:" + SCHEMA_VENDOR + "/screen_view/" + SCHEMA_FORMAT + "/1-0-0"
}

// Updates the tracker with a new subject
func UpdateSubject(subject Subject) {
	s.Subject = subject
}

// Appends another emitter to the tracker
func AddEmitter(emitter Emitter) {
	append(s.Emitter, emitter)
}

// Sends the Payload to the emitter for processing
func SendRequest(payload Payload) {
	finalPayload = ReturnArrayStringify("strval", payload)
	for _, element := range s.Emitter {
		element.SendEvent(finalPayload)
	}
}

// Will force-send all events in the emitter(s) buffers.
// This happens irrespective of whether or not buffer limit has been reached
func FlushEmitters() {
	for _, element = range s.Emitter {
		element.Flush()
	}
}

// Takes a Payload object as a parameter and appends all necessary event data to it

func (t *Tracker) ReturnCompletePayload(payload Payload, context string) Payload{
	var contextEnvelope map[string]string
	if context != nil {
		contextEnvelope["schema"] = t.CONTEXT_SCHEMA 
		contextEnvelope["data"] = context
		payload.AddJson(contextEnvelope, t.EncodeBase64,"cx","co")
	}
	payload.AddDict(t.StdNvPairs)
	payload.AddDict(t.s.GetSubject())
	payload.Add("eid", payload.GenerateUuid())
	return payload
}

// Returns a UUID for a time stamp in Nanoseconds
// TODO(alexanderdean): replace with a UUID library
func (t *Tracker) GenerateUuid() [Size]byte{
	convert := time.Nanoseconds()
	return md5.Sum(convert)
}

// Takes a Payload and a Context and forwards the finalised payload
// map [string]string to the sendRequest function.
func (t *Tracker) Track(payload Payload, context string) {
	payload = t.ReturnCompletePayload(payload, context)
	t.SendRequest(payload.Get())	
}

// Tracks a page view.
func (t *Tracker) TrackPageView(pageUrl string, pageTitle string, referrer string, context string, tstamp string) {
	var payloadEp Payload
	payloadEp.InitPayload(tstamp)
	payloadEp.Add("e", "pv")
	payloadEp.Add("url", pageUrl)
	payloadEp.Add("page", pageTitle)
	payloadEp.Add("refr", referrer)
	t.Track(payloadEp, context)
}
