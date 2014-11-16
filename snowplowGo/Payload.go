package snowplowGo

import (
	"net/url"
	"net/http"
	"time"
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