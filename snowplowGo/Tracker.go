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
	//Tracker Constants
	DEFAULT_BASE_64  = true
	TRACKER_VERSION  = "golang-0.1.0"
	// Schema Constants
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


/**
* Initialize a new tracker instance with emitter(s) and a subject.
*
* @param map[string]string emitterTracker - Emitter object, used for sending event payloads to for processing
* @param Subject Subject - Subject object, contains extra information which is parcelled with the event
* @param string namespace
* @param string AppId
* @param string EncodeBase64 - Boolean stating whether or not to encode certain values as base64
*/
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

/**
* Updates the subject of the tracker with a new subject
*
* @param Subject subject
*/
func UpdateSubject(subject Subject) {
	s.Subject = subject
}

 /**
* Appends another emitter to the tracker
*
* @param Emitter emitter
*/
func AddEmitter(emitter Emitter) {
	append(s.Emitter, emitter)
}

 // Emitter Send Functions
/**
* Sends the Payload to the emitter for processing
*
* @param Payload payload
*/
func SendRequest(payload Payload) {
	finalPayload = ReturnArrayStringify("strval", payload)
	for _, element := range s.Emitter {
		element.SendEvent(finalPayload)
	}
}

 /**
* Will force send all events in the emitter(s) buffers
* This happens irrespective of whether or not buffer limit has been reached
*/
func FlushEmitters() {
	for _, element = range s.Emitter {
		element.Flush()
	}
}
