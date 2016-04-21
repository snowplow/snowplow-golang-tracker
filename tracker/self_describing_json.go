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

type SelfDescribingJson struct {
  schema string
  data interface{}
}

// InitSelfDescribingJson creates a new SelfDescribingJson object.
func InitSelfDescribingJson(schema string, data interface{}) *SelfDescribingJson {
  return &SelfDescribingJson{ schema: schema, data: data }
}

// SetDataWithMap updates the structs data to a new map.
func (s *SelfDescribingJson) SetDataWithMap(data map[string]interface{}) {
  s.data = data
}

// SetDataWithPayload updates the structs data to the contents of a Payload object.
func (s *SelfDescribingJson) SetDataWithPayload(data Payload) {
  s.data = data.Get()
}

// SetDataWithSelfDescribingJson updates the structs data to a JSON.
// Used for nesting SelfDescribingJsons.
func (s *SelfDescribingJson) SetDataWithSelfDescribingJson(data SelfDescribingJson) {
  s.data = data.Get()
}

// Get wraps the schema and data into a JSON.
func (s SelfDescribingJson) Get() map[string]interface{} {
  return map[string]interface{}{
    SCHEMA: s.schema,
    DATA: s.data,
  }
}

// String returns the JSON as a String.
func (s SelfDescribingJson) String() string {
  return MapToJson(s.Get())
}
