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
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// TODO(alexanderdean): why is this variable out here? What is meant to set it?
var TimeStamp string

// TODO(alexanderdean): why are the values integers? Shouldn't they be strings?
var NameValuePair map[string]int64
// TODO(alexanderdean): why does paraValue need to be global? Why not a temp in InitPayload?
// TODO(alexanderdean): ./Payload.go:33: syntax error: unexpected {
var paraValue {}int64

// Initializes a Payload instance, containing an array in which event
// parameters are stored (type int64)
// TODO(alexanderdean): why int64 for event parameters?
func InitPayload() {
	if TimeStamp != nil {
		paraValue := (int64)(TimeStamp)
	} else {
		paraValue = ((int64)(time.Now()) - http.Server.ReadTimeout) * 1000
	}
	Add("dtm", paraValue)
}

// Adds a single map (string[int64]) to the payload
// TODO(alexanderdean): why int64 for event parameters?
func Add(name string, value int64) {
	if value != nil && value != "" {
		NameValuePair[name] = value
	}
}

// Adds an array of name-value pairs to the payload
// TODO(alexanderdean): how does this work without a type for dict?
func AddDict(dict) {
	for name, element := range dict {
		Add(name, element)
	}
}

// Adds a JSON formatted array to the payload.
// Json encodes the array first (turns it into a string) and then will encode (or not) the string in base64
func AddJson(json map[string]string, Base64 bool, NameEncoded string, NameNotEncode string) {
	if json != nil {
		if Base64 {
			Add(NameEncoded, b64.StdEncoding.EncodeToString(json.Marshal(json)))
		}
	} else {
		Add(NameNotEncode, json.Marshal(json))

	}

}
