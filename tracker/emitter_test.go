//
// Copyright (c) 2016-2023 Snowplow Analytics Ltd. All rights reserved.
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
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/memory"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/sqlite3"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/storageiface"
)

func TestEmitterInit(t *testing.T) {
	assert := assert.New(t)
	emitter := InitEmitter(
		RequireCollectorUri("com.acme"),
		OptionRequestType("GET"),
		OptionProtocol("https"),
		OptionSendLimit(1000),
		OptionByteLimitGet(53000),
		OptionByteLimitPost(200000),
		OptionCallback(func(g []CallbackResult, b []CallbackResult) {
			log.Println("Successes: " + common.IntToString(len(g)))
			log.Println("Failures: " + common.IntToString(len(b)))
		}),
		RequireStorage(*memory.Init()),
	)

	// Assert the option builders
	assert.NotNil(emitter)
	assert.Equal("https://com.acme/i", emitter.GetCollectorUrl())
	assert.Equal("com.acme", emitter.CollectorUri)
	assert.Equal("GET", emitter.RequestType)
	assert.Equal("https", emitter.Protocol)
	assert.Equal(1000, emitter.SendLimit)
	assert.Equal(53000, emitter.ByteLimitGet)
	assert.Equal(200000, emitter.ByteLimitPost)
	assert.NotNil(emitter.Storage)
	assert.Nil(emitter.SendChannel)
	assert.NotNil(emitter.Callback)
	assert.NotNil(emitter.HttpClient)
	assert.NotNil(emitter.Storage)
	assert.Equal("memory.StorageMemory", reflect.TypeOf(emitter.Storage).String())

	// Assert defaults
	emitter = InitEmitter(RequireCollectorUri("com.acme"), RequireStorage(*sqlite3.Init("test.db")))
	assert.NotNil(emitter)
	assert.Equal("http://com.acme/com.snowplowanalytics.snowplow/tp2", emitter.GetCollectorUrl())
	assert.Equal("com.acme", emitter.CollectorUri)
	assert.Equal("POST", emitter.RequestType)
	assert.Equal("http", emitter.Protocol)
	assert.Equal(500, emitter.SendLimit)
	assert.Equal(40000, emitter.ByteLimitGet)
	assert.Equal(40000, emitter.ByteLimitPost)
	assert.NotNil(emitter.Storage)
	assert.Nil(emitter.SendChannel)
	assert.Nil(emitter.Callback)
	assert.NotNil(emitter.HttpClient)
	assert.NotNil(emitter.Storage)
	assert.Equal("sqlite3.StorageSQLite3", reflect.TypeOf(emitter.Storage).String())

	// Assert the set functions
	emitter.SetCollectorUri("com.snplow")
	assert.Equal("http://com.snplow/com.snowplowanalytics.snowplow/tp2", emitter.GetCollectorUrl())
	assert.Equal("com.snplow", emitter.CollectorUri)
	emitter.SetRequestType("GET")
	assert.Equal("http://com.snplow/i", emitter.GetCollectorUrl())
	assert.Equal("GET", emitter.RequestType)
	emitter.SetProtocol("https")
	assert.Equal("https://com.snplow/i", emitter.GetCollectorUrl())
	assert.Equal("https", emitter.Protocol)
}

func TestEmitterNoUri(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		if err := recover(); err != nil {
			assert.Equal("FATAL: CollectorUri cannot be empty.", err)
		}
	}()

	emitter := InitEmitter()
	assert.Nil(emitter)
}

func TestEmitterNoStorage(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		if err := recover(); err != nil {
			assert.Equal("FATAL: Storage must be defined.", err)
		}
	}()

	emitter := InitEmitter(RequireCollectorUri("com.acme"))
	assert.Nil(emitter)
}

func TestEmitterBadRequestType(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		if err := recover(); err != nil {
			assert.Equal("FATAL: RequestType did not match either POST or GET.", err)
		}
	}()

	emitter := InitEmitter(
		RequireCollectorUri("com.acme"),
		OptionRequestType("OOPS"),
	)
	assert.Nil(emitter)
}

func TestSingleRowOversize(t *testing.T) {
	assert := assert.New(t)
	emitter := InitEmitter(
		RequireCollectorUri("localhost"),
		OptionRequestType("POST"),
		OptionByteLimitPost(1),
		RequireStorage(*memory.Init()),
	)

	// Single Row > than byte limit POST
	payload0 := *payload.Init()
	payload0.Add("e", common.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
	eventRows := []storageiface.EventRow{{Id: -1, Event: payload0}}
	results := emitter.doSend(eventRows)
	assert.NotNil(results)
	assert.True(len(results) == 1)
	assert.Equal(-1, results[0].ids[0])
	assert.Equal(200, results[0].status)

	emitter = InitEmitter(
		RequireCollectorUri("localhost"),
		OptionRequestType("GET"),
		OptionByteLimitGet(1),
		RequireStorage(*memory.Init()),
	)

	// Single Row > than byte limit GET
	payload1 := *payload.Init()
	payload1.Add("e", common.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
	eventRows2 := []storageiface.EventRow{{Id: -1, Event: payload1}}
	results = emitter.doSend(eventRows2)
	assert.NotNil(results)
	assert.True(len(results) == 1)
	assert.Equal(-1, results[0].ids[0])
	assert.Equal(200, results[0].status)
}

func TestThreeRowsOversize(t *testing.T) {
	assert := assert.New(t)
	emitter := InitEmitter(
		RequireCollectorUri("localhost"),
		OptionRequestType("POST"),
		OptionByteLimitPost(500),
		RequireStorage(*memory.Init()),
	)

	// Three rows > than byte limit POST
	payload1 := *payload.Init()
	payload1.Add("e", common.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
	payload2 := *payload.Init()
	payload2.Add("e", common.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
	payload3 := *payload.Init()
	payload3.Add("e", common.NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))

	eventRows := []storageiface.EventRow{{Id: -1, Event: payload1}, {Id: -1, Event: payload2}, {Id: -1, Event: payload3}}
	results := emitter.doSend(eventRows)
	assert.NotNil(results)
	assert.True(len(results) == 2)
	for _, val := range results {
		for _, id := range val.ids {
			assert.Equal(-1, id)
		}
		assert.Equal(-1, val.status)
	}
}

func TestBadInputToGET(t *testing.T) {
	assert := assert.New(t)
	emitter := InitEmitter(
		RequireCollectorUri("localhost"),
		OptionRequestType("GET"),
		RequireStorage(*memory.Init()),
	)

	// Bad URL
	result := <-emitter.sendGetRequest("", []int{}, false)
	assert.NotNil(result)
	assert.Equal(-1, result.status)

	// Non-Active Collector
	result = <-emitter.sendGetRequest("http://localhost/", []int{}, false)
	assert.NotNil(result)
	assert.Equal(-1, result.status)
}

func TestBadInputToPOST(t *testing.T) {
	assert := assert.New(t)
	emitter := InitEmitter(
		RequireCollectorUri("localhost"),
		OptionRequestType("POST"),
		RequireStorage(*memory.Init()),
	)

	// Bad URL
	result := <-emitter.sendPostRequest("", []int{}, nil, false)
	assert.NotNil(result)
	assert.Equal(-1, result.status)

	// Non-Active Collector
	result = <-emitter.sendPostRequest("http://localhost/", []int{}, []payload.Payload{}, false)
	assert.NotNil(result)
	assert.Equal(-1, result.status)
}
