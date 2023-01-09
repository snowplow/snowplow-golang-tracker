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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
)

func TestPageViewInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := PageViewEvent{
		PageUrl: common.NewString("http://acme.com"),
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("http://acme.com", *event.PageUrl)
	assert.Nil(event.PageTitle)
	assert.Nil(event.Referrer)
	assert.NotEqual("", event.Timestamp)
	assert.NotEqual("", event.EventId)
	assert.Nil(event.TrueTimestamp)
	assert.Nil(event.Subject)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = PageViewEvent{
		PageUrl:       common.NewString("http://acme.com"),
		PageTitle:     common.NewString("Some Title"),
		Referrer:      common.NewString("google.com"),
		Timestamp:     common.NewInt64(1234567890123),
		EventId:       common.NewString("custom-uuid-string"),
		TrueTimestamp: common.NewInt64(9876543210123),
		Subject:       subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("http://acme.com", *event.PageUrl)
	assert.Equal("Some Title", *event.PageTitle)
	assert.Equal("google.com", *event.Referrer)
	assert.Equal(int64(1234567890123), *event.Timestamp)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.Equal(int64(9876543210123), *event.TrueTimestamp)
	assert.NotNil(event.Subject)
	assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"pv\",\"eid\":\"custom-uuid-string\",\"page\":\"Some Title\",\"refr\":\"google.com\",\"ttm\":\"9876543210123\",\"uid\":\"1234\",\"url\":\"http://acme.com\"}", event.Get().String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("PageURL cannot be nil or empty.", err)
		}
	}()
	event = PageViewEvent{}
	event.Init()
}

func TestStructuredInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := StructuredEvent{
		Category: common.NewString("category-1"),
		Action:   common.NewString("selected"),
	}
	event.Init()

	assert.NotNil(event)
	assert.Equal("category-1", *event.Category)
	assert.Equal("selected", *event.Action)
	assert.Nil(event.Label)
	assert.Nil(event.Property)
	assert.Nil(event.Value)
	assert.NotEqual("", *event.Timestamp)
	assert.NotEqual("", *event.EventId)
	assert.Nil(event.TrueTimestamp)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = StructuredEvent{
		Category:  common.NewString("category-1"),
		Action:    common.NewString("selected"),
		Label:     common.NewString("tshirts"),
		Property:  common.NewString("sell"),
		Value:     common.NewFloat64(123.456),
		Timestamp: common.NewInt64(1234567890123),
		EventId:   common.NewString("custom-uuid-string"),
		Subject:   subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("category-1", *event.Category)
	assert.Equal("selected", *event.Action)
	assert.Equal("tshirts", *event.Label)
	assert.Equal("sell", *event.Property)
	assert.Equal(float64(123.456), *event.Value)
	assert.Equal(int64(1234567890123), *event.Timestamp)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.NotNil(event.Subject)
	assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"se\",\"eid\":\"custom-uuid-string\",\"se_ac\":\"selected\",\"se_ca\":\"category-1\",\"se_la\":\"tshirts\",\"se_pr\":\"sell\",\"se_va\":\"123.46\",\"uid\":\"1234\"}", event.Get().String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("Category cannot be nil or empty.", err)
			defer func() {
				if err := recover(); err != nil {
					assert.Equal("Action cannot be nil or empty.", err)
				}
			}()
			event = StructuredEvent{Category: common.NewString("category-1")}
			event.Init()
		}
	}()
	event = StructuredEvent{Action: common.NewString("selected")}
	event.Init()
}

func TestSelfDescribingInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := SelfDescribingEvent{
		Event: InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("{\"data\":{\"e\":\"acme\"},\"schema\":\"iglu:com.acme/event/jsonschema/1-0-0\"}", event.Event.String())
	assert.NotEqual("", *event.Timestamp)
	assert.NotEqual("", *event.EventId)
	assert.Nil(event.TrueTimestamp)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = SelfDescribingEvent{
		Event:     InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
		Timestamp: common.NewInt64(1234567890123),
		EventId:   common.NewString("custom-uuid-string"),
		Subject:   subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("{\"data\":{\"e\":\"acme\"},\"schema\":\"iglu:com.acme/event/jsonschema/1-0-0\"}", event.Event.String())
	assert.Equal(int64(1234567890123), *event.Timestamp)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"ue\",\"eid\":\"custom-uuid-string\",\"ue_pr\":\"{\\\"data\\\":{\\\"data\\\":{\\\"e\\\":\\\"acme\\\"},\\\"schema\\\":\\\"iglu:com.acme/event/jsonschema/1-0-0\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\\\"}\",\"uid\":\"1234\"}", event.Get(false).String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("Event cannot be nil.", err)
		}
	}()
	event = SelfDescribingEvent{}
	event.Init()
}

func TestScreenViewInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := ScreenViewEvent{
		Name: common.NewString("some name"),
		Id:   common.NewString("some id"),
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("some name", *event.Name)
	assert.Equal("some id", *event.Id)
	assert.NotEqual("", *event.Timestamp)
	assert.NotEqual("", *event.EventId)
	assert.Nil(event.TrueTimestamp)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = ScreenViewEvent{
		Name:      common.NewString("some name"),
		Id:        common.NewString("some id"),
		Timestamp: common.NewInt64(1234567890123),
		EventId:   common.NewString("custom-uuid-string"),
		Subject:   subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("some name", *event.Name)
	assert.Equal("some id", *event.Id)
	assert.Equal(int64(1234567890123), *event.Timestamp)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"ue\",\"eid\":\"custom-uuid-string\",\"ue_pr\":\"{\\\"data\\\":{\\\"data\\\":{\\\"id\\\":\\\"some id\\\",\\\"name\\\":\\\"some name\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/screen_view/jsonschema/1-0-0\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\\\"}\",\"uid\":\"1234\"}", event.Get().Get(false).String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("Name and ID cannot both be empty.", err)
		}
	}()
	event = ScreenViewEvent{}
	event.Init()
}

func TestTimingInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := TimingEvent{
		Category: common.NewString("some category"),
		Variable: common.NewString("some variable"),
		Timing:   common.NewInt64(12345678),
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("some category", *event.Category)
	assert.Equal("some variable", *event.Variable)
	assert.Equal(int64(12345678), *event.Timing)
	assert.Nil(event.Label)
	assert.NotEqual("", *event.Timestamp)
	assert.NotEqual("", *event.EventId)
	assert.Nil(event.TrueTimestamp)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = TimingEvent{
		Category:  common.NewString("some category"),
		Variable:  common.NewString("some variable"),
		Timing:    common.NewInt64(12345678),
		Label:     common.NewString("some label"),
		Timestamp: common.NewInt64(1234567890123),
		EventId:   common.NewString("custom-uuid-string"),
		Subject:   subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("some category", *event.Category)
	assert.Equal("some variable", *event.Variable)
	assert.Equal(int64(12345678), *event.Timing)
	assert.Equal("some label", *event.Label)
	assert.Equal(int64(1234567890123), *event.Timestamp)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"ue\",\"eid\":\"custom-uuid-string\",\"ue_pr\":\"{\\\"data\\\":{\\\"data\\\":{\\\"category\\\":\\\"some category\\\",\\\"label\\\":\\\"some label\\\",\\\"timing\\\":12345678,\\\"variable\\\":\\\"some variable\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/timing/jsonschema/1-0-0\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\\\"}\",\"uid\":\"1234\"}", event.Get().Get(false).String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("Category cannot be nil or empty.", err)
			defer func() {
				if err := recover(); err != nil {
					assert.Equal("Variable cannot be nil or empty.", err)
					defer func() {
						if err := recover(); err != nil {
							assert.Equal("Timing cannot be nil.", err)
						}
					}()
					event = TimingEvent{Category: common.NewString("some category"), Variable: common.NewString("some variable")}
					event.Init()
				}
			}()
			event = TimingEvent{Category: common.NewString("some category")}
			event.Init()
		}
	}()
	event = TimingEvent{}
	event.Init()
}

func TestEcommerceTransactionInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := EcommerceTransactionEvent{
		OrderId:    common.NewString("order-1"),
		TotalValue: common.NewFloat64(123456.789),
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("order-1", *event.OrderId)
	assert.Equal(float64(123456.789), *event.TotalValue)
	assert.Nil(event.Affiliation)
	assert.Nil(event.TaxValue)
	assert.Nil(event.Shipping)
	assert.Nil(event.City)
	assert.Nil(event.State)
	assert.Nil(event.Country)
	assert.Nil(event.Currency)
	assert.Nil(event.Items)
	assert.NotEqual("", event.Timestamp)
	assert.NotEqual("", event.EventId)
	assert.Nil(event.TrueTimestamp)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = EcommerceTransactionEvent{
		OrderId:     common.NewString("order-1"),
		TotalValue:  common.NewFloat64(123456.789),
		Affiliation: common.NewString("coffee"),
		TaxValue:    common.NewFloat64(222.35),
		Shipping:    common.NewFloat64(12.5),
		City:        common.NewString("Dijon"),
		State:       common.NewString("Bourgogne"),
		Country:     common.NewString("France"),
		Currency:    common.NewString("EUR"),
		Timestamp:   common.NewInt64(1234567890123),
		EventId:     common.NewString("custom-uuid-string"),
		Subject:     subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("order-1", *event.OrderId)
	assert.Equal(float64(123456.789), *event.TotalValue)
	assert.Equal("coffee", *event.Affiliation)
	assert.Equal(float64(222.35), *event.TaxValue)
	assert.Equal(float64(12.5), *event.Shipping)
	assert.Equal("Dijon", *event.City)
	assert.Equal("Bourgogne", *event.State)
	assert.Equal("France", *event.Country)
	assert.Equal("EUR", *event.Currency)
	assert.Nil(event.Items)
	assert.Equal(int64(1234567890123), *event.Timestamp)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"tr\",\"eid\":\"custom-uuid-string\",\"tr_af\":\"coffee\",\"tr_ci\":\"Dijon\",\"tr_co\":\"France\",\"tr_cu\":\"EUR\",\"tr_id\":\"order-1\",\"tr_sh\":\"12.50\",\"tr_st\":\"Bourgogne\",\"tr_tt\":\"123456.79\",\"tr_tx\":\"222.35\",\"uid\":\"1234\"}", event.Get().String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("OrderID cannot be nil or empty.", err)
			defer func() {
				if err := recover(); err != nil {
					assert.Equal("TotalValue cannot be nil.", err)
				}
			}()
			event = EcommerceTransactionEvent{OrderId: common.NewString("order-2")}
			event.Init()
		}
	}()
	event = EcommerceTransactionEvent{}
	event.Init()
}

func TestEcommerceTransactionItemInit(t *testing.T) {
	assert := assert.New(t)

	// Assert default init
	event := EcommerceTransactionItemEvent{
		Sku:      common.NewString("123-dda"),
		Price:    common.NewFloat64(123456.789),
		Quantity: common.NewInt64(5000),
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("123-dda", *event.Sku)
	assert.Equal(float64(123456.789), *event.Price)
	assert.Equal(int64(5000), *event.Quantity)
	assert.Nil(event.Name)
	assert.Nil(event.Category)
	assert.NotEqual("", *event.EventId)

	subject := InitSubject()
	subject.SetUserId("1234")

	// Assert full init
	event = EcommerceTransactionItemEvent{
		Sku:      common.NewString("123-dda"),
		Price:    common.NewFloat64(123456.789),
		Quantity: common.NewInt64(5000),
		Name:     common.NewString("power pc"),
		Category: common.NewString("servers"),
		EventId:  common.NewString("custom-uuid-string"),
		Subject:  subject,
	}
	event.Init()
	assert.NotNil(event)
	assert.Equal("123-dda", *event.Sku)
	assert.Equal(float64(123456.789), *event.Price)
	assert.Equal(int64(5000), *event.Quantity)
	assert.Equal("power pc", *event.Name)
	assert.Equal("servers", *event.Category)
	assert.Equal("custom-uuid-string", *event.EventId)
	assert.Equal("{\"e\":\"ti\",\"eid\":\"custom-uuid-string\",\"ti_ca\":\"servers\",\"ti_nm\":\"power pc\",\"ti_pr\":\"123456.79\",\"ti_qu\":\"5000\",\"ti_sk\":\"123-dda\",\"uid\":\"1234\"}", event.Get().String())

	defer func() {
		if err := recover(); err != nil {
			assert.Equal("Sku cannot be nil or empty.", err)
			defer func() {
				if err := recover(); err != nil {
					assert.Equal("Price cannot be nil.", err)
					defer func() {
						if err := recover(); err != nil {
							assert.Equal("Quantity cannot be nil.", err)
						}
					}()
					event = EcommerceTransactionItemEvent{Sku: common.NewString("abc123"), Price: common.NewFloat64(123.54)}
					event.Init()
				}
			}()
			event = EcommerceTransactionItemEvent{Sku: common.NewString("abc123")}
			event.Init()
		}
	}()
	event = EcommerceTransactionItemEvent{}
	event.Init()
}

func TestPageViewSubjectOverride(t *testing.T) {
	assert := assert.New(t)

	subject1 := InitSubject()
	subject1.SetUserId("1234")
	subject2 := InitSubject()
	subject2.SetUserId("9876")

	// Assert subject isn't overwritten
	event := PageViewEvent{
		PageUrl: common.NewString("http://acme.com"),
		Subject: subject1,
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "1234")

	// Assert subject is overwritten
	event = PageViewEvent{
		PageUrl: common.NewString("http://acme.com"),
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "9876")
}

func TestStructuredSubjectOverride(t *testing.T) {
	assert := assert.New(t)

	subject1 := InitSubject()
	subject1.SetUserId("1234")
	subject2 := InitSubject()
	subject2.SetUserId("9876")

	// Assert subject isn't overwritten
	event := StructuredEvent{
		Category: common.NewString("category-1"),
		Action:   common.NewString("selected"),
		Subject:  subject1,
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "1234")

	// Assert subject is overwritten
	event = StructuredEvent{
		Category: common.NewString("category-1"),
		Action:   common.NewString("selected"),
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "9876")
}

func TestSelfDescribingSubjectOverride(t *testing.T) {
	assert := assert.New(t)

	subject1 := InitSubject()
	subject1.SetUserId("1234")
	subject2 := InitSubject()
	subject2.SetUserId("9876")

	// Assert subject isn't overwritten
	event := SelfDescribingEvent{
		Event:   InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
		Subject: subject1,
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "1234")

	// Assert subject is overwritten
	event = SelfDescribingEvent{
		Event: InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e": "acme"}),
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "9876")
}

func TestEcommerceTransactionSubjectOverride(t *testing.T) {
	assert := assert.New(t)

	subject1 := InitSubject()
	subject1.SetUserId("1234")
	subject2 := InitSubject()
	subject2.SetUserId("9876")

	// Assert subject isn't overwritten
	event := EcommerceTransactionEvent{
		OrderId:    common.NewString("order-1"),
		TotalValue: common.NewFloat64(123456.789),
		Subject:    subject1,
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "1234")

	// Assert subject is overwritten
	event = EcommerceTransactionEvent{
		OrderId:    common.NewString("order-1"),
		TotalValue: common.NewFloat64(123456.789),
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "9876")
}

func TestEcommerceTransactionItemSubjectOverride(t *testing.T) {
	assert := assert.New(t)

	subject1 := InitSubject()
	subject1.SetUserId("1234")
	subject2 := InitSubject()
	subject2.SetUserId("9876")

	// Assert subject isn't overwritten
	event := EcommerceTransactionItemEvent{
		Sku:      common.NewString("123-dda"),
		Price:    common.NewFloat64(123456.789),
		Quantity: common.NewInt64(5000),
		Subject:  subject1,
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "1234")

	// Assert subject is overwritten
	event = EcommerceTransactionItemEvent{
		Sku:      common.NewString("123-dda"),
		Price:    common.NewFloat64(123456.789),
		Quantity: common.NewInt64(5000),
	}
	event.Init()
	event.SetSubjectIfNil(subject2)
	assert.Equal(event.Subject.Get()[UID], "9876")
}
