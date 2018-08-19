//
// Copyright (c) 2016-2018 Snowplow Analytics Ltd. All rights reserved.
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
  "testing"
  "github.com/stretchr/testify/assert"
)

// TestAdd asserts the behaviour of the payload Add func.
func TestAdd(t *testing.T) {
  assert := assert.New(t)
  payload := *InitPayload()

  // Add an entry
  payload.Add("e", NewString("pv"))
  assert.Equal(1, len(payload.Get()))
  assert.Equal("pv", payload.Get()["e"])
  assert.Equal("{\"e\":\"pv\"}", payload.String())

  // Adding the same key twice overrides the value
  payload.Add("e", NewString("ue"))
  assert.Equal(1, len(payload.Get()))
  assert.Equal("ue", payload.Get()["e"])
  assert.Equal("{\"e\":\"ue\"}", payload.String())

  // Empty values are not added
  payload.Add("empty", NewString(""))
  assert.Equal(1, len(payload.Get()))

  // Empty keys are not added
  payload.Add("", NewString("pv"))
  assert.Equal(1, len(payload.Get()))

  // Nil values are not added
  payload.Add("", nil)
  assert.Equal(1, len(payload.Get()))
}

// TestBadPayload checks that error handling is working.
func TestBadPayload(t *testing.T) {
  assert := assert.New(t)
  payload := Payload{}
  assert.Equal("null", payload.String())
}

// TestAddDict asserts the behaviour of the payload AddDict func.
func TestAddDict(t *testing.T) {
  assert := assert.New(t)
  payload := *InitPayload()
  dict := map[string]string{ "e": "pv", "p": "srv", "res": "", "": "1920x1080" }

  // Add a dictionary of entries
  payload.AddDict(dict)
  assert.Equal(2, len(payload.Get()))
  assert.Equal("pv", payload.Get()["e"])
  assert.Equal("srv", payload.Get()["p"])
  assert.Equal("{\"e\":\"pv\",\"p\":\"srv\"}", payload.String())
}

// TestAddJson asserts the behaviour of the payload AddJson func.
func TestAddJson(t *testing.T) {
  assert := assert.New(t)
  payload := *InitPayload()

  data := map[string]interface{}{}
  data[SCHEMA] = "iglu:com.acme/test_data/jsonschema/1-0-0"
  data[DATA] = map[string]interface{}{ "hello": "data", "world": false, "count": 5 }
  
  json := map[string]interface{}{}
  json[SCHEMA] = SCHEMA_UNSTRUCT_EVENT
  json[DATA] = data

  // Add a JSON structure
  payload.AddJson(json, false, UNSTRUCTURED_ENCODED, UNSTRUCTURED)
  assert.Equal(1, len(payload.Get()))
  assert.Equal("{\"data\":{\"data\":{\"count\":5,\"hello\":\"data\",\"world\":false},\"schema\":\"iglu:com.acme/test_data/jsonschema/1-0-0\"},\"schema\":\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\"}", payload.Get()[UNSTRUCTURED])
  payload.AddJson(json, true, UNSTRUCTURED_ENCODED, UNSTRUCTURED)
  assert.Equal(2, len(payload.Get()))
  assert.Equal("eyJkYXRhIjp7ImRhdGEiOnsiY291bnQiOjUsImhlbGxvIjoiZGF0YSIsIndvcmxkIjpmYWxzZX0sInNjaGVtYSI6ImlnbHU6Y29tLmFjbWUvdGVzdF9kYXRhL2pzb25zY2hlbWEvMS0wLTAifSwic2NoZW1hIjoiaWdsdTpjb20uc25vd3Bsb3dhbmFseXRpY3Muc25vd3Bsb3cvdW5zdHJ1Y3RfZXZlbnQvanNvbnNjaGVtYS8xLTAtMCJ9", payload.Get()[UNSTRUCTURED_ENCODED])
}
