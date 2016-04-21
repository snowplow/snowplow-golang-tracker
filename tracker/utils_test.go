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
  "testing"
  "github.com/stretchr/testify/assert"
  "regexp"
)

// TestGetTimestamp asserts that the GetTimestamp function returns a correct length timestamp.
func TestGetTimestamp(t *testing.T) {
  assert := assert.New(t)
  assert.NotNil(GetTimestamp())
}

// TestGetTimestampString asserts that the GetTimestampString function returns a
// NonNill and correct length timestamp.
func TestGetTimestampString(t *testing.T) {
  assert := assert.New(t)
  assert.NotNil(GetTimestampString())
  assert.Equal(13, len(GetTimestampString()))
}

// TestGetUUID asserts that the UUID is version 4 compliant.
func TestGetUUID(t *testing.T) {
  assert := assert.New(t)
  uuid := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}`)
  assert.NotEqual("", GetUUID())
  assert.Equal(true, uuid.MatchString(GetUUID()))
}

// TestIntToString asserts that Integer to String conversion works.
func TestIntToString(t *testing.T) {
  assert := assert.New(t)
  assert.Equal("5", IntToString(5))
  assert.Equal("0", IntToString(0))
}

// TestInt64ToString asserts that Integer to String conversion works.
func TestInt64ToString(t *testing.T) {
  assert := assert.New(t)
  assert.Equal("5", Int64ToString(NewInt64(5)))
  assert.Equal("0", Int64ToString(NewInt64(0)))
  assert.Equal("", Int64ToString(nil))
}

// TestFloat64ToString asserts that we can print dollar values correctly.
func TestFloat64ToString(t *testing.T) {
  assert := assert.New(t)
  assert.Equal("5050.73", Float64ToString(NewFloat64(5050.72665), 2))
  assert.Equal("", Float64ToString(nil, 2))
}

// TestIntArrayToString asserts that we can build a delimited string from an int array.
func TestIntArrayToString(t *testing.T) {
  assert := assert.New(t)
  intArr := []int{1,2,3,4,5}
  assert.Equal("1,2,3,4,5", IntArrayToString(intArr, ","))
}

// TestCountBytesInString asserts that we can correctly count bytes in a string.
func TestCountBytesInString(t *testing.T) {
  assert := assert.New(t)
  assert.Equal(52, CountBytesInString("http://com.snplow/com.snowplowanalytics.snowplow/tp2"))
}

// TestMapToQueryStr asserts our ability to build valid query strings from maps.
func TestMapToQueryParams(t *testing.T) {
  assert := assert.New(t)
  testMap := map[string]string{
    "k1": "v1",
    "k2": "s p a c e",
    "k3": "s+p+a+c+e",
  }

  assert.Equal("k1=v1&k2=s+p+a+c+e&k3=s%2Bp%2Ba%2Bc%2Be", MapToQueryParams(testMap).Encode())
}

// TestMapToJson asserts the conversion of maps to JSON strings
func TestMapToJson(t *testing.T) {
  assert := assert.New(t)
  testMap := map[string]string{"e":"pv"}
  badMap := map[int]string{1:"pv"}

  assert.Equal("{\"e\":\"pv\"}", MapToJson(testMap))
  assert.Equal("", MapToJson(badMap))
}

// TestMapSerialization asserts the ability to serialize and deserialize maps.
func TestMapSerialization(t *testing.T) {
  assert := assert.New(t)
  testMap := map[string]string{ "hello": "world" }

  byteBuffer := SerializeMap(testMap)
  assert.NotNil(byteBuffer)

  returnedMap, err2 := DeserializeMap(byteBuffer)
  assert.Nil(err2)
  assert.NotNil(returnedMap)
  assert.Equal("world", returnedMap["hello"])

  returnedMap, err2 = DeserializeMap(nil)
  assert.NotNil(err2)
  assert.Nil(returnedMap)
}
