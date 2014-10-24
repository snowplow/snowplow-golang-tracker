package snowplowGo

import (
	"net/url"
	)

constant(

	DEFAULT_REQ_TYPE = "POST"
    DEFAULT_PROTOCOL = "http"
    DEFAULT_BUFFER_SIZE = 10
    BASE_SCHEMA_PATH = "iglu:com.snowplowanalytics.snowplow"
    SCHEMA_TAG = "jsonschema"
	
)

type ConstructEmitter struct{
	postRequestSchema string
	reqType string
	protocol string
	collectorUrl url.URL
	bufferSize int
	buffer []string
	requestsResult []string
}

func InitEmitter(collectorUri string, reqType string, protocol string, bufferSize int) ConstructEmitter{
	
	var s ConstructEmitter
	s.postRequestSchema = BASE_SCHEMA_PATH+"/payload_data/"+SCHEMA_TAG+"/1-0-0"
	if reqType == nil {
		s.reqType = reqType
	}
	else{
		s.reqType = DEFAULT_REQ_TYPE
	}
	if protocol == nil {
		s.protocol = protocol
	}
	else{
		s.protocol = DEFAULT_PROTOCOL
	}
	collectorUrl = ReturnCollectorUrl(collectorUri)
	if (bufferSize == nil) {
            if (s.reqType == "POST") {
                s.bufferSize = DEFAULT_BUFFER_SIZE;
            }
            else {
                s.bufferSize = 1;
            }
    }

}

func ReturnCollectorUrl(host string, s *ConstructEmitter) url.URL{
	switch s.reqType {
    case "POST":  url = s.protocol+"://"+host+"/com.snowplowanalytics.snowplow/tp2"
    			  urlEncoded = url.Parse(url)
    			  return urlEncoded
    case "GET":	url = s.protocol."://".host."/i?"
    		   	urlEncoded = url.Parse(url)
    		   	return urlEncoded
    default: return nil
    }
	
}

