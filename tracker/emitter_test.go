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
  "log"
  "reflect"
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
    OptionDbName("/home/vagrant/test.db"),
    OptionCallback(func(g []CallbackResult, b []CallbackResult) {
      log.Println("Successes: " + IntToString(len(g)))
      log.Println("Failures: " + IntToString(len(b)))
    }),
    OptionStorage(*InitStorageMemory()),
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
  assert.Equal("/home/vagrant/test.db", emitter.DbName)
  assert.NotNil(emitter.Storage)
  assert.Nil(emitter.SendChannel)
  assert.NotNil(emitter.Callback)
  assert.NotNil(emitter.HttpClient)
  assert.NotNil(emitter.Storage)
  assert.Equal("tracker.StorageMemory", reflect.TypeOf(emitter.Storage).String())

  // Assert defaults
  emitter = InitEmitter(RequireCollectorUri("com.acme"), OptionDbName("/home/vagrant/test.db"))
  assert.NotNil(emitter)
  assert.Equal("http://com.acme/com.snowplowanalytics.snowplow/tp2", emitter.GetCollectorUrl())
  assert.Equal("com.acme", emitter.CollectorUri)
  assert.Equal("POST", emitter.RequestType)
  assert.Equal("http", emitter.Protocol)
  assert.Equal(500, emitter.SendLimit)
  assert.Equal(40000, emitter.ByteLimitGet)
  assert.Equal(40000, emitter.ByteLimitPost)
  assert.Equal("/home/vagrant/test.db", emitter.DbName)
  assert.NotNil(emitter.Storage)
  assert.Nil(emitter.SendChannel)
  assert.Nil(emitter.Callback)
  assert.NotNil(emitter.HttpClient)
  assert.NotNil(emitter.Storage)
  assert.Equal("tracker.StorageSQLite3", reflect.TypeOf(emitter.Storage).String())

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
    OptionDbName("/home/vagrant/test.db"),
  )

  // Single Row > than byte limit POST
  payload := *InitPayload()
  payload.Add("e", NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
  eventRows := []EventRow{EventRow{ id: -1, event: payload }}
  results := emitter.doSend(eventRows)
  assert.NotNil(results)
  assert.True(len(results) == 1)
  assert.Equal(-1, results[0].ids[0])
  assert.Equal(200, results[0].status)

  emitter = InitEmitter(
    RequireCollectorUri("localhost"),
    OptionRequestType("GET"),
    OptionByteLimitGet(1),
    OptionDbName("/home/vagrant/test.db"),
  )

  // Single Row > than byte limit GET
  payload1 := *InitPayload()
  payload1.Add("e", NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
  eventRows2 := []EventRow{EventRow{ id: -1, event: payload1 }}
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
    OptionDbName("/home/vagrant/test.db"),
  )

  // Three rows > than byte limit POST
  payload1 := *InitPayload()
  payload1.Add("e", NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
  payload2 := *InitPayload()
  payload2.Add("e", NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))
  payload3 := *InitPayload()
  payload3.Add("e", NewString("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"))

  eventRows := []EventRow{EventRow{ id: -1, event: payload1 }, EventRow{ id: -1, event: payload2 }, EventRow{ id: -1, event: payload3 }}
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
    OptionDbName("/home/vagrant/test.db"),
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
    OptionDbName("/home/vagrant/test.db"),
  )

  // Bad URL
  result := <-emitter.sendPostRequest("", []int{}, nil, false)
  assert.NotNil(result)
  assert.Equal(-1, result.status)

  // Non-Active Collector
  result = <-emitter.sendPostRequest("http://localhost/", []int{}, []Payload{}, false)
  assert.NotNil(result)
  assert.Equal(-1, result.status)
}
