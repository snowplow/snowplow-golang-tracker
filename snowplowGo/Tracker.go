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
	emitter map[string]string
	subject Subject
	EncodeBase64 bool
} 

func InitTracker(emitterTracker map[string]string, subject Subject, namespace string = nil, AppId string = nil, EncodeBase64 string = nil, s Tracker) {
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
}