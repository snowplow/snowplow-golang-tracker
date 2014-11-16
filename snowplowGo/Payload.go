package snowplowGo

import (
	"net/url"
	"net/http"
	"time"
	"encoding/base64"
	"encoding/json"
	)

var TimeStamp string
TimeStamp = nil

NameValuePair = make(map[string]int)
var paraValue int64

func InitPayload() {
	if TimeStamp != nil{
		paraValue = (int64)(TimeStamp)	
	} else{
		paraValue = ((int64)(time.Now()) - http.Server. ReadTimeout)*1000
	}
	Add("dtm", paraValue)	
}

func Add(name string, value int64) {
	if value != nil and value != "" {
		NameValuePair[name] = value	
	}
}

func AddDict(dict){
	for name,element := range dict{
		Add(name, element)
	}		
}

func AddJson(json map([string]string), Base64 bool, NameEncoded string, NameNotEncode string) {
	if json != nil {
		if Base64 {
			Add(NameEncoded, b64.StdEncoding.EncodeToString(json.Marshal(json)))
	}else{
		Add(NameNotEncode, json.Marshal(json))

	}

}
