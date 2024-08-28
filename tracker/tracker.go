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
	"time"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
)

const (
	DEFAULT_PLATFORM = "srv"
	DEFAULT_BASE_64  = true
)

type Tracker struct {
	Emitter      *Emitter
	Subject      *Subject
	Namespace    string
	AppId        string
	Platform     string
	Base64Encode bool
}

// InitTracker creates a new tracker instance linked to an emitter and subject.
// Will assert that the Emitter is valid and not nil.
func InitTracker(options ...func(*Tracker)) *Tracker {
	t := &Tracker{}

	// Set Defaults
	t.Platform = DEFAULT_PLATFORM
	t.Base64Encode = DEFAULT_BASE_64

	// Option parameters
	for _, op := range options {
		op(t)
	}

	// Check Emitter is not nil
	if t.Emitter == nil {
		panic("FATAL: Emitter cannot be nil.")
	}

	return t
}

// --- Require

// RequireEmitter sets the Tracker Emitter
func RequireEmitter(emitter *Emitter) func(t *Tracker) {
	return func(t *Tracker) { t.Emitter = emitter }
}

// --- Option

// OptionSubject sets the Tracker Subject
func OptionSubject(subject *Subject) func(t *Tracker) {
	return func(t *Tracker) { t.Subject = subject }
}

// OptionNamespace sets the Tracker Namespace
func OptionNamespace(namespace string) func(t *Tracker) {
	return func(t *Tracker) { t.Namespace = namespace }
}

// OptionAppId sets the Tracker Application ID
func OptionAppId(appId string) func(t *Tracker) {
	return func(t *Tracker) { t.AppId = appId }
}

// OptionPlatform sets the Tracker Platform
func OptionPlatform(platform string) func(t *Tracker) {
	return func(t *Tracker) { t.Platform = platform }
}

// OptionBase64Encode sets the Tracker base64encode
func OptionBase64Encode(base64Encode bool) func(t *Tracker) {
	return func(t *Tracker) { t.Base64Encode = base64Encode }
}

// --- Utility

// FlushEmitter will force-send all events in the emitter buffer.
func (t Tracker) FlushEmitter() {
	t.Emitter.Flush()
}

func (t *Tracker) waitForEmitter(flushSleepTimeMs int) {
	for {
		if !t.Emitter.IsSending() || t.Emitter.SendChannel == nil {
			break
		}
		time.Sleep(time.Duration(flushSleepTimeMs) * time.Millisecond)
	}
}

// BlockingFlush will block the executing thread until the tracker has fired all events
// from the queue.  Useful for short-lived applications that have to wait for events to
// be fired.
func (t *Tracker) BlockingFlush(flushAttempts int, flushSleepTimeMs int) int {
	t.waitForEmitter(flushSleepTimeMs)

	rowCount := 0
	attemptCount := 0

	for {
		rowCount = len(t.Emitter.Storage.GetAllEventRows())
		if attemptCount >= flushAttempts || rowCount == 0 {
			break
		} else {
			t.FlushEmitter()
			t.waitForEmitter(flushSleepTimeMs)
			attemptCount++
		}
	}

	return rowCount
}

// --- Event Senders

// track takes the event payload and context and completes the build
// process before handing it off to the emitter.
func (t Tracker) track(payload payload.Payload, contexts []SelfDescribingJson) {

	// Add standard KV Pairs
	payload.Add(T_VERSION, common.NewString(TRACKER_VERSION))
	payload.Add(PLATFORM, common.NewString(t.Platform))
	payload.Add(APP_ID, common.NewString(t.AppId))
	payload.Add(NAMESPACE, common.NewString(t.Namespace))

	// Build the final context and add it to the payload
	if contexts != nil && len(contexts) > 0 {
		dataArray := []map[string]interface{}{}
		for _, val := range contexts {
			dataArray = append(dataArray, val.Get())
		}
		contextJson := *InitSelfDescribingJson(SCHEMA_CONTEXTS, dataArray)
		payload.AddJson(contextJson.Get(), t.Base64Encode, CONTEXT_ENCODED, CONTEXT)
	}

	// Add the event to the Emitter.
	t.Emitter.Add(payload)
}

// TrackPageView sends a page view event.
func (t Tracker) TrackPageView(e PageViewEvent) {
	e.Init()
	e.SetSubjectIfNil(t.Subject)
	t.track(e.Get(), e.Contexts)
}

// TrackStructEvent sends a structured event.
func (t Tracker) TrackStructEvent(e StructuredEvent) {
	e.Init()
	e.SetSubjectIfNil(t.Subject)
	t.track(e.Get(), e.Contexts)
}

// TrackSelfDescribingEvent sends a self-described event.
func (t Tracker) TrackSelfDescribingEvent(e SelfDescribingEvent) {
	e.Init()
	e.SetSubjectIfNil(t.Subject)
	t.track(e.Get(t.Base64Encode), e.Contexts)
}

// TrackScreenView sends a screen view event.
func (t Tracker) TrackScreenView(e ScreenViewEvent) {
	e.Init()
	t.TrackSelfDescribingEvent(e.Get())
}

// TrackTiming sends a timing event.
func (t Tracker) TrackTiming(e TimingEvent) {
	e.Init()
	t.TrackSelfDescribingEvent(e.Get())
}

// TrackEcommerceTransaction sends an ecommerce transaction event.
func (t Tracker) TrackEcommerceTransaction(e EcommerceTransactionEvent) {
	e.Init()
	e.SetSubjectIfNil(t.Subject)
	t.track(e.Get(), e.Contexts)
	for _, item := range e.Items {
		t.trackEcommerceTransationItem(item, e.OrderId, e.Currency, e.Timestamp, e.TrueTimestamp)
	}
}

// trackEcommerceTransationItem tracks the individual Ecommerce Items.
func (t Tracker) trackEcommerceTransationItem(e EcommerceTransactionItemEvent, orderId *string, currency *string, timestamp *int64, trueTimestamp *int64) {
	e.Init()
	e.SetSubjectIfNil(t.Subject)
	ep := e.Get()
	ep.Add(TI_ITEM_ID, orderId)
	ep.Add(TI_ITEM_CURRENCY, currency)
	ep.Add(TIMESTAMP, common.NewString(common.Int64ToString(timestamp)))
	ep.Add(TRUE_TIMESTAMP, common.NewString(common.Int64ToString(trueTimestamp)))
	t.track(ep, e.Contexts)
}

// --- Setters

// SetSubject updates the tracker with a new subject.
func (t *Tracker) SetSubject(subject *Subject) {
	t.Subject = subject
}

// SetEmitter updates the tracker with a new emitter.
func (t *Tracker) SetEmitter(emitter *Emitter) {
	t.Emitter = emitter
}

// SetNamespace updates the Tracker namespace value.
func (t *Tracker) SetNamespace(namespace string) {
	t.Namespace = namespace
}

// SetAppId updates the Tracker application id.
func (t *Tracker) SetAppId(appId string) {
	t.AppId = appId
}

// SetPlatform updates the platform from which the event is fired.
func (t *Tracker) SetPlatform(platform string) {
	t.Platform = platform
}

// SetBase64Encode updates whether to base64 encode contexts and unstructured events.
func (t *Tracker) SetBase64Encode(base64Encode bool) {
	t.Base64Encode = base64Encode
}
