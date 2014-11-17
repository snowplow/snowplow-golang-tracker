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

var emitter ConstructEmitter

type ConstructEmitter struct{
	PostRequestSchema string
	ReqType string
	Protocol string
	CollectorUrl url.URL
	BufferSize int
	Buffer []string
	RequestsResult []string
}

func InitEmitter(collectorUri string, reqType string, protocol string, bufferSize int) {
	
	emitter.postRequestSchema = BASE_SCHEMA_PATH+"/payload_data/"+SCHEMA_TAG+"/1-0-0"
	if reqType == nil {
		emitter.reqType = reqType
	}
	else{
		emitter.reqType = DEFAULT_REQ_TYPE
	}
	if protocol == nil {
		emitter.protocol = protocol
	}
	else{
		emitter.protocol = DEFAULT_PROTOCOL
	}
	collectorUrl = ReturnCollectorUrl(collectorUri)
	if (bufferSize == nil) {
            if (emitter.reqType == "POST") {
                emitter.bufferSize = DEFAULT_BUFFER_SIZE;
            }
            else {
                emitter.bufferSize = 1;
            }
    }else{
    	emitter.bufferSize = (int)(BufferSize)
    }
    emitter.bufferSize = nil
    emitter.RequestsResult = nil
}

func ReturnCollectorUrl(host string) url.URL{
	switch emitter.reqType {
    case "POST":  url = emitter.protocol+"://"+host+"/com.snowplowanalytics.snowplow/tp2"
    			  urlEncoded = url.Parse(url)
    			  return urlEncoded
    case "GET":	url = emitter.protocol + "://" + host + "/i?"
    		   	urlEncoded = url.Parse(url)
    		   	return urlEncoded
    default: return nil
    }
	
}

func SendEvent(finalPayload string) {
	Extend(emitter.Buffer, finalPayload)
	if len(emitter.Buffer) >= emitter.BufferSize{
		Flush()
	}

}

func Flush() {
	if len(emitter.Buffer) != 0 {
		if emitter.ReqType == "POST" {
			data := emitter.ReturnPostRequest
			PostRequest(data)
		}else if emitter.ReqType == "GET" {
			for _, value := range emitter.Buffer{
				GetRequest(data)
			}
		}
		emitter.Buffer = nil
	}
}


//Need to complete this function
func GetRequest(data) {
	r := url.Get(HttpBuildQuery(data))
	StoreRequestResults(r)	
}

func postRequest(data, emitter ConstructEmitter){
	m := map[string]string{
		"'Content-Type' => 'application/json; charset=utf-8'"
	}
	//post method to be made properly here
	r := url.Post(emitter.CollectorUrl)
	//
	StoreRequestResults(r)
}

func ReturnPostRequest(emitter ConstructEmitter) []string{
	dataPostRequest = make(map[string]string)
	dataPostRequest["schema"] = emitter.PostRequestSchema
	dataPostRequest["data"] = []string
	for _,element := range emitter.Buffer {
		append(dataPostRequest["data"], element)
	}
	return dataPostRequest
}

func StoreRequestResults(r RequestsResponse){
	storeArray = make(map[string]string)
	storeArray["url"] = r.url
	storeArray["code"] = r.StatusCode
	storeArray["headers"] = r.headers
	storeArray["body"] = r.body
	storeArray["raw"] = r.raw
	append(emitter.RequestsResult, storeArray)
}
