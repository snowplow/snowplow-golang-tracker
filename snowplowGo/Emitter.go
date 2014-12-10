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
	"fmt"
	"net/url"
)

// Constants for sending a payload
const (
	DEFAULT_REQ_TYPE    = "POST"
	DEFAULT_PROTOCOL    = "http"
	DEFAULT_BUFFER_SIZE = 10

	SCHEMA_VENDOR       = "com.snowplowanalytics.snowplow"
	SCHEMA_FORMAT       = "jsonschema"
)

type Emitter struct {
	PostRequestSchema string
	ReqType           string
	Protocol          string
	CollectorUrl      url.URL
	BufferSize        int
	Buffer            []string
	RequestsResult    []string
}

// Intitialize emitter to send event data to a collector
func (e *Emitter) InitEmitter(collectorUri string, reqType string, protocol string, bufferSize int) Emitter {

	emitter.PostRequestSchema = fmt.Sprintf("iglu:%s/payload_data/%s/1-0-2", SCHEMA_VENDOR, SCHEMA_FORMAT)
	if reqType != "" {
		emitter.ReqType = reqType
	} else {
		emitter.ReqType = DEFAULT_REQ_TYPE
	}
	if protocol != "" {
		emitter.Protocol = protocol
	} else {
		emitter.Protocol = DEFAULT_PROTOCOL
	}
	emitter.CollectorUrl = e.ReturnCollectorUrl(collectorUri)
	if bufferSize == nil {
		if emitter.ReqType == "POST" {
			emitter.bufferSize = DEFAULT_BUFFER_SIZE
		} else {
			emitter.bufferSize = 1
		}
	} else {
		emitter.bufferSize = (int)(BufferSize)
	}
	emitter.bufferSize = nil // TODO(alexanderdean): won't this overwrite the assignments above?
	emitter.RequestsResult = nil
	return emitter
}

// Returns the collector URL based on: request type, protocol and host given
// If a bad type is given in emitter creation, returns nil.
func (e *Emitter) ReturnCollectorUrl( host string) url.URL {
	switch e.reqType {
	case "POST":
		url = e.protocol + "://" + host + "/com.snowplowanalytics.snowplow/tp2"
		urlEncoded = url.Parse(url)
		return urlEncoded
	case "GET":
		url = e.protocol + "://" + host + "/i?"
		urlEncoded = url.Parse(url)
		return urlEncoded
	default:
		return nil
	}

}

// Pushes the event payload into the emitter buffer.
// When buffer is full it flushes the buffer.
func (e *Emitter) SendEvent( finalPayload []string) {
	Extend(e.Buffer, finalPayload)
	if len(e.Buffer) >= e.BufferSize {
		Flush()
	}

}

// Flushes the event buffer of the emitter
// Checks which send type the emitter is using and forwards data accordingly
// Resets the buffer to nothing after flushing
func (e *Emitter) Flush(emitter Emitter) {
	if len(emitter.Buffer) != 0 {
		if emitter.ReqType == "POST" {
			data := emitter.ReturnPostRequest
			e.PostRequest(data)
		} else if emitter.ReqType == "GET" {
			for _, value := range emitter.Buffer {
				e.GetRequest(data)
			}
		}
		emitter.Buffer = nil
	}
}

// Sends the payload to the collector via a GET request
func (e *Emitter) GetRequest(data []string) {
	r := url.Get(HttpBuildQuery(data))
	e.StoreRequestResults(r)
}

// Sends the payload to the collector via a POST request
func (e *Emitter) PostRequest(data []string) {
	m = make(map[string]string)
	m["Content-Type"] = "application/json; charset=utf-8"
	//post method to be made properly here
	r := url.Post(e.CollectorUrl)
	//
	e.StoreRequestResults(r)
}

// Returns an array formatted to be ready for a POST request
func (e *Emitter) ReturnPostRequest( ) []string {
	dataPostRequest := make(map[string][]map[string]string)
	dataPostRequest["schema"] = e.PostRequestSchema
	for _, element := range e.Buffer {
		append(dataPostRequest["data"], element)
	}
	return dataPostRequest
}

// Stores all of the parameters of the request's response
// into a dynamic array for use in unit testing
// TODO(alexanderdean): is there a cleaner way of doing this?
func (e *Emitter) StoreRequestResults(r RequestsResponse) {
	storeArray = make(map[string]string)
	storeArray["url"] = r.url
	storeArray["code"] = r.StatusCode
	storeArray["headers"] = r.headers
	storeArray["body"] = r.body
	storeArray["raw"] = r.raw
	append(emitter.RequestsResult, storeArray)
}
