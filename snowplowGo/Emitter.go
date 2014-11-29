package main

import (
	"net/url"
)

const (
	DEFAULT_REQ_TYPE    = "POST"
	DEFAULT_PROTOCOL    = "http"
	DEFAULT_BUFFER_SIZE = 10
	BASE_SCHEMA_PATH    = "iglu:com.snowplowanalytics.snowplow"
	SCHEMA_TAG          = "jsonschema"
)

var emitter Emitter

type Emitter struct {
	PostRequestSchema string
	ReqType           string
	Protocol          string
	CollectorUrl      url.URL
	BufferSize        int
	Buffer            []string
	RequestsResult    []string
}

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

func (e *Emitter) ReturnCollectorUrl(emitter Emitter, host string) url.URL {
	switch emitter.reqType {
	case "POST":
		url = emitter.protocol + "://" + host + "/com.snowplowanalytics.snowplow/tp2"
		urlEncoded = url.Parse(url)
		return urlEncoded
	case "GET":
		url = emitter.protocol + "://" + host + "/i?"
		urlEncoded = url.Parse(url)
		return urlEncoded
	default:
		return nil
	}

}

func (e *Emitter) SendEvent(emitter Emitter, finalPayload string) {
	Extend(emitter.Buffer, finalPayload)
	if len(emitter.Buffer) >= emitter.BufferSize {
		Flush()
	}

}

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

//Need to complete this function
func (e *Emitter) GetRequest(data) {
	r := url.Get(HttpBuildQuery(data))
	e.StoreRequestResults(r)
}

func (e *Emitter) PostRequest(data, emitter Emitter) {
	m = make(map[string]string)
	m["Content-Type"] = "application/json; charset=utf-8"
	//post method to be made properly here
	r := url.Post(emitter.CollectorUrl)
	//
	e.StoreRequestResults(r)
}

func (e *Emitter) ReturnPostRequest(emitter Emitter) []string {
	dataPostRequest := make(map[string][]map[string]string)
	dataPostRequest["schema"] = emitter.PostRequestSchema
	for _, element := range emitter.Buffer {
		append(dataPostRequest["data"], element)
	}
	return dataPostRequest
}

func (e *Emitter) StoreRequestResults(r RequestsResponse) {
	storeArray = make(map[string]string)
	storeArray["url"] = r.url
	storeArray["code"] = r.StatusCode
	storeArray["headers"] = r.headers
	storeArray["body"] = r.body
	storeArray["raw"] = r.raw
	append(emitter.RequestsResult, storeArray)
}
