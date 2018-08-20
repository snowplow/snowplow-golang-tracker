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
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSelfDescribingJson asserts creating a new SDJ.
func TestSelfDescribingJson(t *testing.T) {
	assert := assert.New(t)
	sdj := *InitSelfDescribingJson("iglu:com.acme", map[string]interface{}{"hello": "world"})
	assert.NotNil(sdj)
	assert.Equal("{\"data\":{\"hello\":\"world\"},\"schema\":\"iglu:com.acme\"}", sdj.String())

	// Assert setters
	sdj.SetDataWithMap(map[string]interface{}{"snowplow": 1})
	assert.Equal("{\"data\":{\"snowplow\":1},\"schema\":\"iglu:com.acme\"}", sdj.String())

	payload := *InitPayload()
	payload.Add("e", NewString("pv"))
	sdj.SetDataWithPayload(payload)
	assert.Equal("{\"data\":{\"e\":\"pv\"},\"schema\":\"iglu:com.acme\"}", sdj.String())

	sdj.SetDataWithSelfDescribingJson(sdj)
	assert.Equal("{\"data\":{\"data\":{\"e\":\"pv\"},\"schema\":\"iglu:com.acme\"},\"schema\":\"iglu:com.acme\"}", sdj.String())
}
