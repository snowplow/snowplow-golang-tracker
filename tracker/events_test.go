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
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestPageViewInit(t *testing.T) {
  assert := assert.New(t)

  // Assert default init
  event := PageViewEvent{ 
    PageUrl: NewString("http://acme.com"),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("http://acme.com", *event.PageUrl)
  assert.Nil(event.PageTitle)
  assert.Nil(event.Referrer)
  assert.NotEqual("", event.Timestamp)
  assert.NotEqual("", event.EventId)
  assert.Nil(event.TrueTimestamp)

  // Assert full init
  event = PageViewEvent{ 
    PageUrl: NewString("http://acme.com"),
    PageTitle: NewString("Some Title"),
    Referrer: NewString("google.com"),
    Timestamp: NewInt64(1234567890123),
    EventId: NewString("custom-uuid-string"),
    TrueTimestamp: NewInt64(9876543210123),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("http://acme.com", *event.PageUrl)
  assert.Equal("Some Title", *event.PageTitle)
  assert.Equal("google.com", *event.Referrer)
  assert.Equal(int64(1234567890123), *event.Timestamp)
  assert.Equal("custom-uuid-string", *event.EventId)
  assert.Equal(int64(9876543210123), *event.TrueTimestamp)
  assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"pv\",\"eid\":\"custom-uuid-string\",\"page\":\"Some Title\",\"refr\":\"google.com\",\"ttm\":\"9876543210123\",\"url\":\"http://acme.com\"}", event.Get().String())

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
    Category: NewString("category-1"),
    Action: NewString("selected"),
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

  // Assert full init
  event = StructuredEvent{ 
    Category: NewString("category-1"),
    Action: NewString("selected"),
    Label: NewString("tshirts"),
    Property: NewString("sell"),
    Value: NewFloat64(123.456),
    Timestamp: NewInt64(1234567890123),
    EventId: NewString("custom-uuid-string"),
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
  assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"se\",\"eid\":\"custom-uuid-string\",\"se_ac\":\"selected\",\"se_ca\":\"category-1\",\"se_la\":\"tshirts\",\"se_pr\":\"sell\",\"se_va\":\"123.46\"}", event.Get().String())

  defer func() {
    if err := recover(); err != nil {
      assert.Equal("Category cannot be nil or empty.", err)
      defer func() {
        if err := recover(); err != nil {
          assert.Equal("Action cannot be nil or empty.", err)
        }
      }()
      event = StructuredEvent{ Category: NewString("category-1") }
      event.Init()
    }
  }()
  event = StructuredEvent{ Action: NewString("selected") }
  event.Init()
}

func TestSelfDescribingInit(t *testing.T) {
  assert := assert.New(t)

  // Assert default init
  event := SelfDescribingEvent{
    Event: InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e":"acme"}),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("{\"data\":{\"e\":\"acme\"},\"schema\":\"iglu:com.acme/event/jsonschema/1-0-0\"}", event.Event.String())
  assert.NotEqual("", *event.Timestamp)
  assert.NotEqual("", *event.EventId)
  assert.Nil(event.TrueTimestamp)

   // Assert full init
  event = SelfDescribingEvent{ 
    Event: InitSelfDescribingJson("iglu:com.acme/event/jsonschema/1-0-0", map[string]string{"e":"acme"}),
    Timestamp: NewInt64(1234567890123),
    EventId: NewString("custom-uuid-string"),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("{\"data\":{\"e\":\"acme\"},\"schema\":\"iglu:com.acme/event/jsonschema/1-0-0\"}", event.Event.String())
  assert.Equal(int64(1234567890123), *event.Timestamp)
  assert.Equal("custom-uuid-string", *event.EventId)
  assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"ue\",\"eid\":\"custom-uuid-string\",\"ue_pr\":\"{\\\"data\\\":{\\\"data\\\":{\\\"e\\\":\\\"acme\\\"},\\\"schema\\\":\\\"iglu:com.acme/event/jsonschema/1-0-0\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\\\"}\"}", event.Get(false).String())

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
    Name: NewString("some name"),
    Id: NewString("some id"),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("some name", *event.Name)
  assert.Equal("some id", *event.Id)
  assert.NotEqual("", *event.Timestamp)
  assert.NotEqual("", *event.EventId)
  assert.Nil(event.TrueTimestamp)

   // Assert full init
  event = ScreenViewEvent{ 
    Name: NewString("some name"),
    Id: NewString("some id"),
    Timestamp: NewInt64(1234567890123),
    EventId: NewString("custom-uuid-string"),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("some name", *event.Name)
  assert.Equal("some id", *event.Id)
  assert.Equal(int64(1234567890123), *event.Timestamp)
  assert.Equal("custom-uuid-string", *event.EventId)
  assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"ue\",\"eid\":\"custom-uuid-string\",\"ue_pr\":\"{\\\"data\\\":{\\\"data\\\":{\\\"id\\\":\\\"some id\\\",\\\"name\\\":\\\"some name\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/screen_view/jsonschema/1-0-0\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\\\"}\"}", event.Get().Get(false).String())

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
    Category: NewString("some category"), 
    Variable: NewString("some variable"),
    Timing: NewInt64(12345678),
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

   // Assert full init
  event = TimingEvent{ 
    Category: NewString("some category"),
    Variable: NewString("some variable"),
    Timing: NewInt64(12345678),
    Label: NewString("some label"),
    Timestamp: NewInt64(1234567890123),
    EventId: NewString("custom-uuid-string"),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("some category", *event.Category)
  assert.Equal("some variable", *event.Variable)
  assert.Equal(int64(12345678), *event.Timing)
  assert.Equal("some label", *event.Label)
  assert.Equal(int64(1234567890123), *event.Timestamp)
  assert.Equal("custom-uuid-string", *event.EventId)
  assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"ue\",\"eid\":\"custom-uuid-string\",\"ue_pr\":\"{\\\"data\\\":{\\\"data\\\":{\\\"category\\\":\\\"some category\\\",\\\"label\\\":\\\"some label\\\",\\\"timing\\\":12345678,\\\"variable\\\":\\\"some variable\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/timing/jsonschema/1-0-0\\\"},\\\"schema\\\":\\\"iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0\\\"}\"}", event.Get().Get(false).String())

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
          event = TimingEvent{ Category: NewString("some category"), Variable: NewString("some variable") }
          event.Init()
        }
      }()
      event = TimingEvent{ Category: NewString("some category") }
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
    OrderId: NewString("order-1"),
    TotalValue: NewFloat64(123456.789),
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

   // Assert full init
  event = EcommerceTransactionEvent{ 
    OrderId: NewString("order-1"), 
    TotalValue: NewFloat64(123456.789), 
    Affiliation: NewString("coffee"),
    TaxValue: NewFloat64(222.35),
    Shipping: NewFloat64(12.5),
    City: NewString("Dijon"),
    State: NewString("Bourgogne"),
    Country: NewString("France"),
    Currency: NewString("EUR"),
    Timestamp: NewInt64(1234567890123),
    EventId: NewString("custom-uuid-string"),
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
  assert.Equal("{\"dtm\":\"1234567890123\",\"e\":\"tr\",\"eid\":\"custom-uuid-string\",\"tr_af\":\"coffee\",\"tr_ci\":\"Dijon\",\"tr_co\":\"France\",\"tr_cu\":\"EUR\",\"tr_id\":\"order-1\",\"tr_sh\":\"12.50\",\"tr_st\":\"Bourgogne\",\"tr_tt\":\"123456.79\",\"tr_tx\":\"222.35\"}", event.Get().String())

  defer func() {
    if err := recover(); err != nil {
      assert.Equal("OrderID cannot be nil or empty.", err)
      defer func() {
        if err := recover(); err != nil {
          assert.Equal("TotalValue cannot be nil.", err)
        }
      }()
      event = EcommerceTransactionEvent{ OrderId: NewString("order-2") }
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
    Sku: NewString("123-dda"),
    Price: NewFloat64(123456.789),
    Quantity: NewInt64(5000),
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("123-dda", *event.Sku)
  assert.Equal(float64(123456.789), *event.Price)
  assert.Equal(int64(5000), *event.Quantity)
  assert.Nil(event.Name)
  assert.Nil(event.Category)
  assert.NotEqual("", *event.EventId)

   // Assert full init
  event = EcommerceTransactionItemEvent{ 
    Sku: NewString("123-dda"), 
    Price: NewFloat64(123456.789), 
    Quantity: NewInt64(5000),
    Name: NewString("power pc"),
    Category: NewString("servers"),
    EventId: NewString("custom-uuid-string"), 
  }
  event.Init()
  assert.NotNil(event)
  assert.Equal("123-dda", *event.Sku)
  assert.Equal(float64(123456.789), *event.Price)
  assert.Equal(int64(5000), *event.Quantity)
  assert.Equal("power pc", *event.Name)
  assert.Equal("servers", *event.Category)
  assert.Equal("custom-uuid-string", *event.EventId)
  assert.Equal("{\"e\":\"ti\",\"eid\":\"custom-uuid-string\",\"ti_ca\":\"servers\",\"ti_nm\":\"power pc\",\"ti_pr\":\"123456.79\",\"ti_qu\":\"5000\",\"ti_sk\":\"123-dda\"}", event.Get().String())

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
          event = EcommerceTransactionItemEvent{ Sku: NewString("abc123"), Price: NewFloat64(123.54) }
          event.Init()
        }
      }()
      event = EcommerceTransactionItemEvent{ Sku: NewString("abc123") }
      event.Init()
    }
  }()
  event = EcommerceTransactionItemEvent{}
  event.Init()
}
