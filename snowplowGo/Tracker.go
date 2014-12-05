/*
Tracker.go
Copyright (c) 2014 Snowplow Analytics Ltd. All rights reserved.
This program is licensed to you under the Apache License Version 2.0,
and you may not use this file except in compliance with the Apache License
Version 2.0. You may obtain a copy of the Apache License Version 2.0 at
http://www.apache.org/licenses/LICENSE-2.0.
Unless required by applicable law or agreed to in writing,
software distributed under the Apache License Version 2.0 is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
express or implied. See the Apache License Version 2.0 for the specific
language governing permissions and limitations there under.
Authors: Aalekh Nigam
Copyright: Copyright (c) 2014 Snowplow Analytics Ltd
License: Apache License Version 2.0
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
	BASE_SCHEMA_PATH = "iglu:com.snowplowanalytics.snowplow"
	SCHEMA_TAG       = "jsonschema"
)

type Tracker struct {
	Emitter      Emitter
	subject      Subject
	EncodeBase64 bool
}

// Building Json Schema
type JsonSchema struct {
	ContextSchema       string
	UnstructEventSchema string
	ScreenViewSchema    string
}

var scehma JsonSchema

var StdNvPairs map[string]string

var s Tracker

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

	schema.ContextSchema = BASE_SCHEMA_PATH + "/contexts/" + SCHEMA_TAG + "/1-0-0"
	schema.UnstructEventSchema = BASE_SCHEMA_PATH + "/unstruct_event/" + SCHEMA_TAG + "/1-0-0"
	schema.ScreenViewSchema = BASE_SCHEMA_PATH + "/screen_view/" + SCHEMA_TAG + "/1-0-0"
}

func UpdateSubject(subject Subject) {
	s.Subject = subject
}

func AddEmitter(emitter Emitter) {
	append(s.Emitter, emitter)
}

func SendRequest(payload Payload) {
	finalPayload = ReturnArrayStringify("strval", payload)
	for _, element := range s.Emitter {
		element.SendEvent(finalPayload)
	}
}

func FlushEmitters() {
	for _, element = range s.Emitter {
		element.Flush()
	}
}
