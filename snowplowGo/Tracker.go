package snowplowGo

import (
	"net/url"
	"net/http"
	"time"
	"encoding/base64"
	"encoding/json"
	)

const(
	DEFAULT_BASE_64 = true
	TRACKER_VERSION = "golang-0.1.0"
	BASE_SCHEMA_PATH = "iglu:com.snowplowanalytics.snowplow"
	SCHEMA_TAG = "jsonschema"
)

type Tracker struct{
	Emitter ConstructEmitter 
	subject Subject
	EncodeBase64 bool
} 

// Building Json Schema
type JsonSchema struct{
	ContextSchema string
	UnstructEventSchema string
	ScreenViewSchema string
}

var scehma JsonSchema

var StdNvPairs map[string]string

var s Tracker

func InitTracker(emitterTracker map[string]string, subject Subject, namespace string = nil, AppId string = nil, EncodeBase64 string = nil) {
	if len(emitterTracker) > 0 {
		s.emitter = emitterTracker
	}else{
		s.emitter = emitterTracker
	}	
	s.subject = subject
	if s.EncodeBase64 != nil {
		s.EncodeBase64 = StringToBool(s.EncodeBase64)	
	}else{
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
	finalPayload = ReturnArrayStringify('strval', payload)
	for _,element := range s.Emitter{
		element.SendEvent(finalPayload)
	}
}

func FlushEmitters() {
	for _,element = range s.Emitter{
		element.Flush()
	}
}

