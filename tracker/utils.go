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
  "encoding/gob"
  "bytes"
  "time"
  "github.com/twinj/uuid"
  "strconv"
  "net/url"
  "encoding/json"
)

// NewString returns a pointer to a string.
func NewString(val string) *string {
  return &val
}

// NewInt64 returns a pointer to an int64.
func NewInt64(val int64) *int64 {
  return &val
}

// NewFloat64 returns a pointer to a float64.
func NewFloat64(val float64) *float64 {
  return &val
}

// GetTimestamp returns the current unix timestamp in milliseconds
func GetTimestamp() int64 {
  return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// GetTimestampString returns the current unix timestamp in milliseconds
func GetTimestampString() string {
  return strconv.FormatInt(GetTimestamp(), 10)
}

// GetUUID generates a Version 4 UUID string.
func GetUUID() string {
  return uuid.NewV4().String()
}

// IntToString converts an Integer to a String.
func IntToString(value int) string {
  return strconv.Itoa(value)
}

// Int64ToString converts an Integer of 64 bits to a String.
func Int64ToString(value *int64) string {
  if value != nil {
    return strconv.FormatInt(*value, 10)
  } else {
    return ""
  }
}

// Float64ToString does conversion of floats to string values.
func Float64ToString(value *float64, places int) string {
  if value != nil {
    return strconv.FormatFloat(*value, 'f', places, 64)
  } else {
    return ""
  }
}

// IntArrayToString converts an array of integers to a string delimited
// by a string of your choice.
func IntArrayToString(values []int, delimiter string) string {
  var strBuffer bytes.Buffer
  for index, val := range values {
    strBuffer.WriteString(IntToString(val))
    if index < (len(values) - 1) {
      strBuffer.WriteString(delimiter)
    }
  }
  return strBuffer.String()
}

// CountBytesInString takes a string and gets a byte count.
func CountBytesInString(str string) int {
  return len(str)
}

// MapToQueryParams takes a map of string keys and values and builds an encoded query string.
func MapToQueryParams(m map[string]string) url.Values {
  params := url.Values{}
  for key, value := range m {
    params.Add(key, value)
  }
  return params
}

// MapToString takes a generic and converts it to a JSON representation.
func MapToJson(m interface{}) string {
  b, err := json.Marshal(m)
  if err == nil {
    return string(b)
  } else {
    return ""
  }
}

// SerializeMap takes a map and attempts to convert it to a byte buffer.
func SerializeMap(m map[string]string) []byte {
  b := new(bytes.Buffer)
  e := gob.NewEncoder(b)
  e.Encode(m)
  return b.Bytes()
}

// DeserializeMap takes a byte buffer and attempts to convert it back to a map.
func DeserializeMap(b []byte) (map[string]string, error) {
  var decodedMap map[string]string
  d := gob.NewDecoder(bytes.NewBuffer(b))

  err := d.Decode(&decodedMap)
  if err != nil {
    return nil, err
  } else {
    return decodedMap, nil
  }
}
