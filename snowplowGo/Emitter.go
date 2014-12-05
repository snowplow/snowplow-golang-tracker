/*
Emitter.go
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
package main

import (
	"net/url"
)

//Basic Schema Constants for specfying basic schemas for Events
const (
	DEFAULT_REQ_TYPE    = "POST"
	DEFAULT_PROTOCOL    = "http"
	DEFAULT_BUFFER_SIZE = 10
	BASE_SCHEMA_PATH    = "iglu:com.snowplowanalytics.snowplow"
	SCHEMA_TAG          = "jsonschema"
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

/**
* Intitialize emitter to send event data to a collector
*
* @param string collectorUri - URI of the collector
* @param string|null reqType - Request Type (POST or GET)
* @param string|null protocol - Protocol to be used (HTTP or HTTPS)
* @param int|null bufferSize - Buffer Size, amount of events to be stored before sending
*/

func (e *Emitter) InitEmitter(collectorUri string, reqType string, protocol string, bufferSize int) Emitter {

	emitter.PostRequestSchema = BASE_SCHEMA_PATH + "/payload_data/" + SCHEMA_TAG + "/1-0-0"
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
	emitter.bufferSize = nil
	emitter.RequestsResult = nil
	return emitter
}

 // Make Functions
/**
* Returns the collector URL based on: request type, protocol and host given
* IF a bad type is given in emitter creation returns NULL
*
* @param string host - The Collector URI to be used for tracking
* @return string|null collector_url - Returns the Collector URL
*/

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

 /**
* Pushes the event payload into the emitter buffer.
* When buffer is full it flushes the buffer.
*
* @param array finalPayload - Takes the Trackers Payload as a parameter
*/

func (e *Emitter) SendEvent( finalPayload []string) {
	Extend(e.Buffer, finalPayload)
	if len(e.Buffer) >= e.BufferSize {
		Flush()
	}

}

 /**
* Flushes the event buffer of the emitter
* Checks which send type the emitter is using and forwards data accordingly
* Resets the buffer to nothing after flushing
*/

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

 // Send Functions
/**
* Using a GET Request sends the data to a collector
*
* @param array data - The array which is going to be sent in the GET Request
*/
func (e *Emitter) GetRequest(data []string) {
	r := url.Get(HttpBuildQuery(data))
	e.StoreRequestResults(r)
}


/**
* Using a POST Request sends the data to a collector
*
* @param array data - Is the array which is going to be sent in the POST Request
*/

func (e *Emitter) PostRequest(data []string) {
	m = make(map[string]string)
	m["Content-Type"] = "application/json; charset=utf-8"
	//post method to be made properly here
	r := url.Post(e.CollectorUrl)
	//
	e.StoreRequestResults(r)
}

/**
* Returns an array which has been formatted to be ready for a POST Request
**/

func (e *Emitter) ReturnPostRequest( ) []string {
	dataPostRequest := make(map[string][]map[string]string)
	dataPostRequest["schema"] = e.PostRequestSchema
	for _, element := range e.Buffer {
		append(dataPostRequest["data"], element)
	}
	return dataPostRequest
}

 /**
* Stores all of the parameters of the Request Response into a Dynamic Array for use in unit testing.
* @param Requests_Response r - Is the response from a GET or POST request
**/

func (e *Emitter) StoreRequestResults(r RequestsResponse) {
	storeArray = make(map[string]string)
	storeArray["url"] = r.url
	storeArray["code"] = r.StatusCode
	storeArray["headers"] = r.headers
	storeArray["body"] = r.body
	storeArray["raw"] = r.raw
	append(emitter.RequestsResult, storeArray)
}
