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
  "log"
  "net/url"
  "net/http"
  "errors"
  "bytes"
  "time"
)

const (
  DEFAULT_REQ_TYPE        = "POST"
  DEFAULT_PROTOCOL        = "http"
  DEFAULT_SEND_LIMIT      = 500
  DEFAULT_BYTE_LIMIT_GET  = 40000
  DEFAULT_BYTE_LIMIT_POST = 40000
  DEFAULT_DB_NAME         = "events.db"
  POST_WRAPPER_BYTES      = 88 // "schema":"iglu:com.snowplowanalytics.snowplow/payload_data/jsonschema/1-0-3","data":[]
  POST_STM_BYTES          = 22 // "stm":"1443452851000"
)

type SendResult struct {
  ids    []int
  status int
}

type CallbackResult struct {
  Count  int
  Status int
}

type Emitter struct {
  CollectorUri  string
  CollectorUrl  url.URL
  RequestType   string
  Protocol      string
  SendLimit     int
  ByteLimitGet  int
  ByteLimitPost int
  DbName        string
  Storage       Storage
  SendChannel   chan bool
  Callback      func(successCount []CallbackResult, failureCount []CallbackResult)
  HttpClient    http.Client
}

// InitEmitter creates a new Emitter object which handles
// storing and sending Snowplow Events.
func InitEmitter(options ...func(*Emitter)) *Emitter {
  e := &Emitter{}

  // Set Defaults
  e.RequestType = DEFAULT_REQ_TYPE
  e.Protocol = DEFAULT_PROTOCOL
  e.SendLimit = DEFAULT_SEND_LIMIT
  e.ByteLimitGet = DEFAULT_BYTE_LIMIT_GET
  e.ByteLimitPost = DEFAULT_BYTE_LIMIT_POST
  e.DbName = DEFAULT_DB_NAME

  // Option parameters
  for _, op := range options { op(e) }

  // Check collector URI is not empty
  if e.CollectorUri == "" {
    panic("FATAL: CollectorUri cannot be empty.")
  } else {
    collectorUrl, err := returnCollectorUrl(e.RequestType, e.Protocol, e.CollectorUri)
    if err != nil {
      panic(err.Error())
    } else {
      e.CollectorUrl = *collectorUrl
    }
  }

  // Setup Event Storage
  e.Storage = *InitStorage(e.DbName)

  // Setup HttpClient
  timeout := time.Duration(5 * time.Second)
  e.HttpClient = http.Client{
    Timeout: timeout,
  }

  return e
}

// --- Require

// RequireCollectorUri sets the Emitters collector URI.
func RequireCollectorUri(collectorUri string) func(e *Emitter) {
  return func(e *Emitter) { e.CollectorUri = collectorUri }
}

// --- Option

// OptionRequestType sets the request type to use (GET or POST).
func OptionRequestType(requestType string) func(e *Emitter) {
  return func(e *Emitter) { e.RequestType = requestType }
}

// OptionProtocol sets the protocol type to use (http or https).
func OptionProtocol(protocol string) func(e *Emitter) {
  return func(e *Emitter) { e.Protocol = protocol }
}

// OptionSendLimit sets the send limit for the emitter.
func OptionSendLimit(sendLimit int) func(e *Emitter) {
  return func(e *Emitter) { e.SendLimit = sendLimit }
}

// OptionByteLimitGet sets the byte limit for GET requests.
func OptionByteLimitGet(byteLimitGet int) func(e *Emitter) {
  return func(e *Emitter) { e.ByteLimitGet = byteLimitGet }
}

// OptionByteLimitPost sets the byte limit for POST requests.
func OptionByteLimitPost(byteLimitPost int) func(e *Emitter) {
  return func(e *Emitter) { e.ByteLimitPost = byteLimitPost }
}

// OptionDbName sets the name of the storage database.
func OptionDbName(dbName string) func(e *Emitter) {
  return func(e *Emitter) { e.DbName = dbName }
}

// OptionCallback sets a custom callback for the emitter loop.
func OptionCallback(callback func(successCount []CallbackResult, failureCount []CallbackResult)) func(e *Emitter) {
  return func(e *Emitter) { e.Callback = callback }
}

// --- Event Handlers

// Add will push an event to the database and will then initiate a sending loop.
func (e *Emitter) Add(payload Payload) {
  e.Storage.AddEventRow(payload)
  e.start()
}

// Flush will attempt to start the send loop regardless of an event coming in.
func (e *Emitter) Flush() {
  e.start()
}

// Stop waits for the send channel to have a value and then resets it to nil.
func (e *Emitter) Stop() {
  <-e.SendChannel
  e.SendChannel = nil
}

// start will begin the sending loop.
func (e *Emitter) start() {
  if e.SendChannel == nil || !e.IsSending() {
    e.SendChannel = make(chan bool, 1)
    go func() {
      var done bool
      defer func() {
        e.SendChannel <- done
      }()

      for {
        eventRows := e.Storage.GetEventRowsWithinRange(e.SendLimit)

        // If there are no events in the database exit
        if len(eventRows) == 0 {
          break
        }
        results := e.doSend(eventRows)

        // Process results
        ids := []int{}
        successes := []CallbackResult{}
        failures := []CallbackResult{}

        for _, res := range results {

          count := len(res.ids)
          status := res.status

          if status >= 200 && status < 400 {
            ids = append(ids, res.ids...)
            successes = append(successes, CallbackResult{count, status})
          } else {
            failures = append(failures, CallbackResult{count, status})
          }
        }

        if e.Callback != nil {
          e.Callback(successes, failures)
        }

        // If all the events failed to be sent exit
        if len(successes) == 0 && len(failures) > 0 {
          break
        }

        e.Storage.DeleteEventRows(ids)
      }
      done = true
    }()
  }
}

// doSend will send all of the eventsRows it is given.
func (e *Emitter) doSend(eventRows []EventRow) []SendResult {
  futures := []<-chan SendResult{}
  url := e.GetCollectorUrl()

  if e.RequestType == "POST" {
    ids := []int{}
    payloads := []Payload{}
    totalByteSize := 0

    for _, val := range eventRows {
      byteSize := CountBytesInString(val.event.String()) + POST_STM_BYTES
      if byteSize + POST_WRAPPER_BYTES > e.ByteLimitPost {
        // A single payload has exceeded the Byte Limit
        futures = append(futures, e.sendPostRequest(url, []int{val.id}, []Payload{val.event}, true))
      } else if (totalByteSize + byteSize + POST_WRAPPER_BYTES + (len(payloads) - 1)) > e.ByteLimitPost {
        // Byte limit reached
        futures = append(futures, e.sendPostRequest(url, ids, payloads, false))

        // Reset accumulators
        ids = []int{val.id}
        payloads = []Payload{val.event}
        totalByteSize = byteSize
      } else {
        ids = append(ids, val.id)
        payloads = append(payloads, val.event)
        totalByteSize += byteSize
      }
    }
    if len(payloads) > 0 {
      futures = append(futures, e.sendPostRequest(url, ids, payloads, false))
    }
  } else if e.RequestType == "GET" {
    for _, val := range eventRows {
      val.event.Add(SENT_TIMESTAMP, NewString(GetTimestampString()))
      queryString := MapToQueryParams(val.event.Get()).Encode()
      oversize := CountBytesInString(queryString) > e.ByteLimitGet
      futures = append(futures, e.sendGetRequest(url + "?" + queryString, []int{val.id}, oversize))
    }
  }

  // Wait for all Futures to complete
  results := []SendResult{}
  for _, future := range futures {
    results = append(results, <-future)
  }

  return results
}

// SendGetRequest sends a payload to the collector endpoint via GET.
func (e *Emitter) sendGetRequest(url string, ids []int, oversize bool) <-chan SendResult {
  c := make(chan SendResult, 1)
  go func() {
    var result SendResult
    defer func() {
      c <- result
    }()

    status := -1; if oversize { status = 200 }

    req, _ := http.NewRequest("GET", url, nil)
    req.Close = true

    resp, err := e.HttpClient.Do(req)
    if err != nil {
        log.Println(err.Error())
        result = SendResult{ ids: ids, status: status }
        return
    }
    resp.Body.Close()

    status = resp.StatusCode; if oversize { status = 200 }
    result = SendResult{ ids: ids, status: status }
  }()
  return c
}

// SendPostRequest sends an array of Payloads together to the collector endpoint via POST.
func (e *Emitter) sendPostRequest(url string, ids []int, body []Payload, oversize bool) <- chan SendResult {
  c := make(chan SendResult, 1)
  go func() {
    var result SendResult
    defer func() {
      c <- result
    }()

    status := -1; if oversize { status = 200 }

    postEnvelope := map[string]interface{}{
      SCHEMA: SCHEMA_PAYLOAD_DATA,
      DATA: addSentTimeToEvents(body),
    }

    req, _ := http.NewRequest("POST", url, bytes.NewBufferString(MapToJson(postEnvelope)))
    req.Close = true
    req.Header.Set("Content-Type", POST_CONTENT_TYPE)

    resp, err := e.HttpClient.Do(req)
    if err != nil {
        log.Println(err.Error())
        result = SendResult{ ids: ids, status: status }
        return
    }
    resp.Body.Close()

    status = resp.StatusCode; if oversize { status = 200 }
    result = SendResult{ ids: ids, status: status }
  }()
  return c
}

// --- Helpers

// IsSending checks whether the send channel has finished.
func (e Emitter) IsSending() bool {
  return len(e.SendChannel) == 0
}

// returnCollectorUrl builds and returns the full collector URL to be used.
func returnCollectorUrl(requestType string, protocol string, collectorUri string) (*url.URL, error) {
  var rawUrl string
  switch requestType {
    case "POST":
      rawUrl = protocol + "://" + collectorUri + "/" + POST_PROTOCOL_VENDOR + "/" + POST_PROTOCOL_VERSION
    case "GET":
      rawUrl = protocol + "://" + collectorUri + "/" + GET_PROTOCOL_PATH
    default:
      return nil, errors.New("FATAL: RequestType did not match either POST or GET.")
  }
  return url.Parse(rawUrl)
}

// addSentTimeToEvents ranges over an array of events and appends the same timestamp to them all.
func addSentTimeToEvents(events []Payload) []map[string]string {
  eventMaps := []map[string]string{}
  stm := NewString(GetTimestampString())
  for _, p := range events {
    p.Add(SENT_TIMESTAMP, stm)
    eventMaps = append(eventMaps, p.Get())
  }
  return eventMaps
}

// --- Getters & Setters

// GetCollectorUrl returns the stringified collector URL.
func (e Emitter) GetCollectorUrl() string {
  return e.CollectorUrl.String()
}

// SetCollectorUri sets a new Collector URI and updates the Collector URL.
func (e *Emitter) SetCollectorUri(collectorUri string) {
  collectorUrl, err := returnCollectorUrl(e.RequestType, e.Protocol, collectorUri)
  if err == nil {
    e.CollectorUrl = *collectorUrl
    e.CollectorUri = collectorUri
  }
}

// SetRequestType sets a new Request Type and updates the Collector URL.
func (e *Emitter) SetRequestType(requestType string) {
  collectorUrl, err := returnCollectorUrl(requestType, e.Protocol, e.CollectorUri)
  if err == nil {
    e.CollectorUrl = *collectorUrl
    e.RequestType = requestType
  }
}

// SetProtocol sets a new Protocol and updates the Collector URL.
func (e *Emitter) SetProtocol(protocol string) {
  collectorUrl, err := returnCollectorUrl(e.RequestType, protocol, e.CollectorUri)
  if err == nil {
    e.CollectorUrl = *collectorUrl
    e.Protocol = protocol
  }
}
