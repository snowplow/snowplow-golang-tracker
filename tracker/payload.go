//
// Copyright (c) 2016 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package tracker

import (
  "encoding/base64"
  "encoding/json"
)

type Payload struct {
  pairs map[string]string
}

// InitPayload returns a new payload object.
func InitPayload() *Payload {
  return &Payload{ pairs: map[string]string{} }
}

// Add pushes a key value pair to the payload.
func (p Payload) Add(key string, value *string) {
  if key != "" && value != nil && *value != "" {
    p.pairs[key] = *value
  }
}

// AddDict pushes an array of key-value pairs to the payload.
func (p Payload) AddDict(dict map[string]string) {
  for name, element := range dict {
    p.Add(name, NewString(element))
  }
}

// AddJson pushes a JSON formatted array to the payload.
// Json encodes the array first (turns it into a string) and then will encode (or not) the string in base64.
func (p Payload) AddJson(instance map[string]interface{}, isBase64 bool, keyEncoded string, keyNotEncoded string) {
  if instance != nil {
    b, err := json.Marshal(instance)
    if err == nil {
      if isBase64 {
        p.Add(keyEncoded, NewString(base64.StdEncoding.EncodeToString(b)))
      } else {
        p.Add(keyNotEncoded, NewString(string(b)))
      }
    }
  }
}

// Get returns the payload as a map[string]string.
func (p Payload) Get() map[string]string {
  return p.pairs
}

// String returns a JSON representation of the internal Map.
func (p Payload) String() string {
  return MapToJson(p.pairs)
}
