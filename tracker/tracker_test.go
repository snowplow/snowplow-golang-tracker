//
// Copyright (c) 2016-2020 Snowplow Analytics Ltd. All rights reserved.
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
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestTrackerInit(t *testing.T) {
	assert := assert.New(t)
	tracker := InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme"),
			OptionDbName("test.db"),
		)),
		OptionSubject(InitSubject()),
		OptionNamespace("namespace"),
		OptionAppId("app-id"),
		OptionPlatform("mob"),
		OptionBase64Encode(false),
	)

	// Assert the option builders
	assert.NotNil(tracker)
	assert.NotNil(tracker.Emitter)
	assert.NotNil(tracker.Subject)
	assert.Equal("namespace", tracker.Namespace)
	assert.Equal("app-id", tracker.AppId)
	assert.Equal("mob", tracker.Platform)
	assert.Equal(false, tracker.Base64Encode)

	// Assert defaults
	tracker = InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme"),
			OptionDbName("test.db"),
		)),
	)
	assert.NotNil(tracker)
	assert.NotNil(tracker.Emitter)
	assert.Nil(tracker.Subject)
	assert.Equal("", tracker.Namespace)
	assert.Equal("", tracker.AppId)
	assert.Equal("srv", tracker.Platform)
	assert.Equal(true, tracker.Base64Encode)

	// Assert the set functions
	tracker.SetSubject(InitSubject())
	tracker.SetEmitter(InitEmitter(
		RequireCollectorUri("com.new"),
		OptionDbName("test.db"),
	))
	tracker.SetNamespace("some-namespace")
	tracker.SetAppId("some-app-id")
	tracker.SetPlatform("web")
	tracker.SetBase64Encode(false)
	assert.NotNil(tracker.Emitter)
	assert.NotNil(tracker.Subject)
	assert.Equal("some-namespace", tracker.Namespace)
	assert.Equal("some-app-id", tracker.AppId)
	assert.Equal("web", tracker.Platform)
	assert.Equal(false, tracker.Base64Encode)

	// Assert panic for no emitter set
	defer func() {
		if err := recover(); err != nil {
			assert.Equal("FATAL: Emitter cannot be nil.", err)
		}
	}()
	tracker = InitTracker()
}

func TestTrackFunctionsGET(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"http://com.acme.collector/i",
		httpmock.NewStringResponder(200, ""),
	)

	tracker := InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme.collector"),
			OptionRequestType("GET"),
			OptionStorage(*InitStorageMemory()),
			OptionCallback(func(g []CallbackResult, b []CallbackResult) {
				log.Println("Successes: " + IntToString(len(g)))
				log.Println("Failures: " + IntToString(len(b)))
			}),
			OptionHttpClient(http.DefaultClient),
		)),
		OptionSubject(InitSubject()),
		OptionNamespace("namespace"),
		OptionAppId("app-id"),
		OptionPlatform("mob"),
		OptionBase64Encode(false),
	)
	assert.NotNil(tracker)

	contextArray := []SelfDescribingJson{
		*InitSelfDescribingJson("iglu:com.acme/context/jsonschema/1-0-0", map[string]string{"e": "context"}),
	}

	// Track the bare minimum for all event types
	tracker.TrackPageView(PageViewEvent{
		PageUrl:  NewString("acme.com"),
		Contexts: contextArray,
	})
	tracker.TrackStructEvent(StructuredEvent{
		Category: NewString("some category"),
		Action:   NewString("some action"),
		Contexts: contextArray,
	})
	tracker.TrackSelfDescribingEvent(SelfDescribingEvent{
		Event:    InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
		Contexts: contextArray,
	})
	tracker.TrackScreenView(ScreenViewEvent{
		Id:       NewString("Screen ID"),
		Contexts: contextArray,
	})
	tracker.TrackTiming(TimingEvent{
		Category: NewString("Timing Category"),
		Variable: NewString("Some var"),
		Timing:   NewInt64(124578),
		Contexts: contextArray,
	})
	tracker.TrackEcommerceTransaction(EcommerceTransactionEvent{
		OrderId:    NewString("order-id"),
		TotalValue: NewFloat64(12345.68),
		Contexts:   contextArray,
		Items: []EcommerceTransactionItemEvent{
			{
				Sku:      NewString("a sku"),
				Price:    NewFloat64(12345.68),
				Quantity: NewInt64(1),
				Contexts: contextArray,
			},
		},
	})
	tracker.Emitter.Stop()
	tracker.BlockingFlush(5, 10)
}

func TestTrackFunctionsPOST(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://com.acme.collector/com.snowplowanalytics.snowplow/tp2",
		httpmock.NewStringResponder(200, ""),
	)

	tracker := InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme.collector"),
			OptionRequestType("POST"),
			OptionDbName("test.db"),
			OptionCallback(func(g []CallbackResult, b []CallbackResult) {
				log.Println("Successes: " + IntToString(len(g)))
				log.Println("Failures: " + IntToString(len(b)))
			}),
			OptionHttpClient(http.DefaultClient),
		)),
		OptionSubject(InitSubject()),
		OptionNamespace("namespace"),
		OptionAppId("app-id"),
		OptionPlatform("mob"),
		OptionBase64Encode(false),
	)
	assert.NotNil(tracker)

	contextArray := []SelfDescribingJson{
		*InitSelfDescribingJson("iglu:com.acme/context/jsonschema/1-0-0", map[string]string{"e": "context"}),
	}

	// Track the bare minimum for all event types
	tracker.TrackPageView(PageViewEvent{
		PageUrl:  NewString("acme.com"),
		Contexts: contextArray,
	})
	tracker.TrackStructEvent(StructuredEvent{
		Category: NewString("some category"),
		Action:   NewString("some action"),
		Contexts: contextArray,
	})
	tracker.TrackSelfDescribingEvent(SelfDescribingEvent{
		Event:    InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
		Contexts: contextArray,
	})
	tracker.TrackScreenView(ScreenViewEvent{
		Id:       NewString("Screen ID"),
		Contexts: contextArray,
	})
	tracker.TrackTiming(TimingEvent{
		Category: NewString("Timing Category"),
		Variable: NewString("Some var"),
		Timing:   NewInt64(124578),
		Contexts: contextArray,
	})
	tracker.TrackEcommerceTransaction(EcommerceTransactionEvent{
		OrderId:    NewString("order-id"),
		TotalValue: NewFloat64(12345.68),
		Contexts:   contextArray,
		Items: []EcommerceTransactionItemEvent{
			{
				Sku:      NewString("a sku"),
				Price:    NewFloat64(12345.68),
				Quantity: NewInt64(1),
				Contexts: contextArray,
			},
		},
	})
	tracker.Emitter.Stop()
	tracker.BlockingFlush(5, 10)
}

func TestTrackFunctionsFailingGET(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"http://com.acme.collector/i",
		httpmock.NewStringResponder(404, ""),
	)

	tracker := InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme.collector"),
			OptionRequestType("GET"),
			OptionDbName("test.db"),
			OptionCallback(func(g []CallbackResult, b []CallbackResult) {
				log.Println("Successes: " + IntToString(len(g)))
				log.Println("Failures: " + IntToString(len(b)))
			}),
			OptionHttpClient(http.DefaultClient),
		)),
		OptionSubject(InitSubject()),
		OptionNamespace("namespace"),
		OptionAppId("app-id"),
		OptionPlatform("mob"),
		OptionBase64Encode(false),
	)
	assert.NotNil(tracker)

	tracker.TrackPageView(PageViewEvent{PageUrl: NewString("acme.com")})
	tracker.BlockingFlush(5, 10)
}

func TestTrackFunctionsFailingPOST(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://com.acme.collector/com.snowplowanalytics.snowplow/tp2",
		httpmock.NewStringResponder(404, ""),
	)

	tracker := InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme.collector"),
			OptionRequestType("POST"),
			OptionStorage(*InitStorageMemory()),
			OptionCallback(func(g []CallbackResult, b []CallbackResult) {
				log.Println("Successes: " + IntToString(len(g)))
				log.Println("Failures: " + IntToString(len(b)))
			}),
			OptionHttpClient(http.DefaultClient),
		)),
		OptionSubject(InitSubject()),
		OptionNamespace("namespace"),
		OptionAppId("app-id"),
		OptionPlatform("mob"),
		OptionBase64Encode(false),
	)
	assert.NotNil(tracker)

	tracker.TrackPageView(PageViewEvent{PageUrl: NewString("acme.com")})
	tracker.BlockingFlush(5, 10)
}

func TestTrackFunctionsWithEventSubject(t *testing.T) {
	assert := assert.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://com.acme.collector/com.snowplowanalytics.snowplow/tp2",
		httpmock.NewStringResponder(200, ""),
	)

	tracker := InitTracker(
		RequireEmitter(InitEmitter(
			RequireCollectorUri("com.acme.collector"),
			OptionRequestType("POST"),
			OptionDbName("test.db"),
			OptionCallback(func(g []CallbackResult, b []CallbackResult) {
				log.Println("Successes: " + IntToString(len(g)))
				log.Println("Failures: " + IntToString(len(b)))
			}),
			OptionHttpClient(http.DefaultClient),
		)),
		OptionSubject(InitSubject()),
		OptionNamespace("namespace"),
		OptionAppId("app-id"),
		OptionPlatform("mob"),
		OptionBase64Encode(false),
	)
	assert.NotNil(tracker)

	contextArray := []SelfDescribingJson{
		*InitSelfDescribingJson("iglu:com.acme/context/jsonschema/1-0-0", map[string]string{"e": "context"}),
	}

	// Track the bare minimum for all event types with a event level subject
	tracker.TrackPageView(PageViewEvent{
		PageUrl:  NewString("acme.com"),
		Contexts: contextArray,
		Subject:  InitSubject(),
	})
	tracker.TrackStructEvent(StructuredEvent{
		Category: NewString("some category"),
		Action:   NewString("some action"),
		Contexts: contextArray,
		Subject:  InitSubject(),
	})
	tracker.TrackSelfDescribingEvent(SelfDescribingEvent{
		Event:    InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
		Contexts: contextArray,
		Subject:  InitSubject(),
	})
	tracker.TrackScreenView(ScreenViewEvent{
		Id:       NewString("Screen ID"),
		Contexts: contextArray,
		Subject:  InitSubject(),
	})
	tracker.TrackTiming(TimingEvent{
		Category: NewString("Timing Category"),
		Variable: NewString("Some var"),
		Timing:   NewInt64(124578),
		Contexts: contextArray,
		Subject:  InitSubject(),
	})
	tracker.TrackEcommerceTransaction(EcommerceTransactionEvent{
		OrderId:    NewString("order-id"),
		TotalValue: NewFloat64(12345.68),
		Contexts:   contextArray,
		Subject:    InitSubject(),
		Items: []EcommerceTransactionItemEvent{
			{
				Sku:      NewString("a sku"),
				Price:    NewFloat64(12345.68),
				Quantity: NewInt64(1),
				Contexts: contextArray,
			},
		},
	})
	tracker.Emitter.Stop()
	tracker.BlockingFlush(5, 10)
}
