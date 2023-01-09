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
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
)

// --- PageView Event

type PageViewEvent struct {
	PageUrl       *string              // Required
	PageTitle     *string              // Optional
	Referrer      *string              // Optional
	Timestamp     *int64               // Optional
	EventId       *string              // Optional
	TrueTimestamp *int64               // Optional
	Contexts      []SelfDescribingJson // Optional
	Subject       *Subject             // Optional
}

// Init checks and validates the struct.
func (e *PageViewEvent) Init() {
	if e.PageUrl == nil || *e.PageUrl == "" {
		panic("PageURL cannot be nil or empty.")
	}
	if e.Timestamp == nil {
		e.Timestamp = common.NewInt64(common.GetTimestamp())
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// SetSubjectIfNil will set the event level subject if a Subject isn't already set
func (e *PageViewEvent) SetSubjectIfNil(subject *Subject) {
	if e.Subject == nil {
		e.Subject = subject
	}
}

// Get returns the event payload.
func (e PageViewEvent) Get() payload.Payload {
	ep := *payload.Init()
	ep.Add(EVENT, common.NewString(EVENT_PAGE_VIEW))
	ep.Add(PAGE_URL, e.PageUrl)
	ep.Add(PAGE_TITLE, e.PageTitle)
	ep.Add(PAGE_REFR, e.Referrer)
	ep.Add(TIMESTAMP, common.NewString(common.Int64ToString(e.Timestamp)))
	ep.Add(EID, e.EventId)
	ep.Add(TRUE_TIMESTAMP, common.NewString(common.Int64ToString(e.TrueTimestamp)))
	if e.Subject != nil {
		ep.AddDict(e.Subject.Get())
	}
	return ep
}

// --- Structured Event

type StructuredEvent struct {
	Category      *string              // Required
	Action        *string              // Required
	Label         *string              // Optional
	Property      *string              // Optional
	Value         *float64             // Optional
	Timestamp     *int64               // Optional
	EventId       *string              // Optional
	TrueTimestamp *int64               // Optional
	Contexts      []SelfDescribingJson // Optional
	Subject       *Subject             // Optional
}

// Init checks and validates the struct.
func (e *StructuredEvent) Init() {
	if e.Category == nil || *e.Category == "" {
		panic("Category cannot be nil or empty.")
	}
	if e.Action == nil || *e.Action == "" {
		panic("Action cannot be nil or empty.")
	}
	if e.Timestamp == nil {
		e.Timestamp = common.NewInt64(common.GetTimestamp())
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// SetSubjectIfNil will set the event level subject if a Subject isn't already set
func (e *StructuredEvent) SetSubjectIfNil(subject *Subject) {
	if e.Subject == nil {
		e.Subject = subject
	}
}

// Get returns the event payload.
func (e StructuredEvent) Get() payload.Payload {
	ep := *payload.Init()
	ep.Add(EVENT, common.NewString(EVENT_STRUCTURED))
	ep.Add(SE_CATEGORY, e.Category)
	ep.Add(SE_ACTION, e.Action)
	ep.Add(SE_LABEL, e.Label)
	ep.Add(SE_PROPERTY, e.Property)
	ep.Add(SE_VALUE, common.NewString(common.Float64ToString(e.Value, 2)))
	ep.Add(TIMESTAMP, common.NewString(common.Int64ToString(e.Timestamp)))
	ep.Add(EID, e.EventId)
	ep.Add(TRUE_TIMESTAMP, common.NewString(common.Int64ToString(e.TrueTimestamp)))
	if e.Subject != nil {
		ep.AddDict(e.Subject.Get())
	}
	return ep
}

// --- SelfDescribing Event

type SelfDescribingEvent struct {
	Event         *SelfDescribingJson  // Required
	Timestamp     *int64               // Optional
	EventId       *string              // Optional
	TrueTimestamp *int64               // Optional
	Contexts      []SelfDescribingJson // Optional
	Subject       *Subject             // Optional
}

// Init checks and validates the struct.
func (e *SelfDescribingEvent) Init() {
	if e.Event == nil {
		panic("Event cannot be nil.")
	}
	if e.Timestamp == nil {
		e.Timestamp = common.NewInt64(common.GetTimestamp())
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// SetSubjectIfNil will set the event level subject if a Subject isn't already set
func (e *SelfDescribingEvent) SetSubjectIfNil(subject *Subject) {
	if e.Subject == nil {
		e.Subject = subject
	}
}

// Get returns the event payload.
func (e SelfDescribingEvent) Get(base64Encode bool) payload.Payload {
	sdj := *InitSelfDescribingJson(SCHEMA_UNSTRUCT_EVENT, e.Event.Get())
	ep := *payload.Init()
	ep.Add(EVENT, common.NewString(EVENT_UNSTRUCTURED))
	ep.Add(TIMESTAMP, common.NewString(common.Int64ToString(e.Timestamp)))
	ep.Add(EID, e.EventId)
	ep.Add(TRUE_TIMESTAMP, common.NewString(common.Int64ToString(e.TrueTimestamp)))
	ep.AddJson(sdj.Get(), base64Encode, UNSTRUCTURED_ENCODED, UNSTRUCTURED)
	if e.Subject != nil {
		ep.AddDict(e.Subject.Get())
	}
	return ep
}

// --- ScreenView Event

type ScreenViewEvent struct {
	Name          *string              // Optional
	Id            *string              // Optional
	Timestamp     *int64               // Optional
	EventId       *string              // Optional
	TrueTimestamp *int64               // Optional
	Contexts      []SelfDescribingJson // Optional
	Subject       *Subject             // Optional
}

// Init checks and validates the struct.
func (e *ScreenViewEvent) Init() {
	if (e.Name == nil || *e.Name == "") && (e.Id == nil || *e.Id == "") {
		panic("Name and ID cannot both be empty.")
	}
	if e.Timestamp == nil {
		e.Timestamp = common.NewInt64(common.GetTimestamp())
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// Get returns the event payload.
func (e ScreenViewEvent) Get() SelfDescribingEvent {
	ep := *payload.Init()
	ep.Add(SV_NAME, e.Name)
	ep.Add(SV_ID, e.Id)
	sdj := InitSelfDescribingJson(SCHEMA_SCREEN_VIEW, ep.Get())
	return SelfDescribingEvent{Event: sdj, Timestamp: e.Timestamp, EventId: e.EventId, TrueTimestamp: e.TrueTimestamp, Contexts: e.Contexts, Subject: e.Subject}
}

// --- Timing Event

type TimingEvent struct {
	Category      *string              // Required
	Variable      *string              // Required
	Timing        *int64               // Required
	Label         *string              // Optional
	Timestamp     *int64               // Optional
	EventId       *string              // Optional
	TrueTimestamp *int64               // Optional
	Contexts      []SelfDescribingJson // Optional
	Subject       *Subject             // Optional
}

// Init checks and validates the struct.
func (e *TimingEvent) Init() {
	if e.Category == nil || *e.Category == "" {
		panic("Category cannot be nil or empty.")
	}
	if e.Variable == nil || *e.Variable == "" {
		panic("Variable cannot be nil or empty.")
	}
	if e.Timing == nil {
		panic("Timing cannot be nil.")
	}
	if e.Timestamp == nil {
		e.Timestamp = common.NewInt64(common.GetTimestamp())
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// Get returns the event payload.
func (e TimingEvent) Get() SelfDescribingEvent {
	data := map[string]interface{}{
		UT_CATEGORY: *e.Category,
		UT_VARIABLE: *e.Variable,
		UT_TIMING:   *e.Timing,
	}
	if e.Label != nil && *e.Label != "" {
		data[UT_LABEL] = *e.Label
	}

	sdj := InitSelfDescribingJson(SCHEMA_USER_TIMINGS, data)
	return SelfDescribingEvent{Event: sdj, Timestamp: e.Timestamp, EventId: e.EventId, TrueTimestamp: e.TrueTimestamp, Contexts: e.Contexts, Subject: e.Subject}
}

// --- EcommerceTransaction Event

type EcommerceTransactionEvent struct {
	OrderId       *string                         // Required
	TotalValue    *float64                        // Required
	Affiliation   *string                         // Optional
	TaxValue      *float64                        // Optional
	Shipping      *float64                        // Optional
	City          *string                         // Optional
	State         *string                         // Optional
	Country       *string                         // Optional
	Currency      *string                         // Optional
	Items         []EcommerceTransactionItemEvent // Optional
	Timestamp     *int64                          // Optional
	EventId       *string                         // Optional
	TrueTimestamp *int64                          // Optional
	Contexts      []SelfDescribingJson            // Optional
	Subject       *Subject                        // Optional
}

// Init checks and validates the struct.
func (e *EcommerceTransactionEvent) Init() {
	if e.OrderId == nil || *e.OrderId == "" {
		panic("OrderID cannot be nil or empty.")
	}
	if e.TotalValue == nil {
		panic("TotalValue cannot be nil.")
	}
	if e.Timestamp == nil {
		e.Timestamp = common.NewInt64(common.GetTimestamp())
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// SetSubjectIfNil will set the event level subject if a Subject isn't already set
func (e *EcommerceTransactionEvent) SetSubjectIfNil(subject *Subject) {
	if e.Subject == nil {
		e.Subject = subject
	}
}

// Get returns the event payload.
func (e EcommerceTransactionEvent) Get() payload.Payload {
	ep := *payload.Init()
	ep.Add(EVENT, common.NewString(EVENT_ECOMM))
	ep.Add(TR_ID, e.OrderId)
	ep.Add(TR_TOTAL, common.NewString(common.Float64ToString(e.TotalValue, 2)))
	ep.Add(TR_AFFILIATION, e.Affiliation)
	ep.Add(TR_TAX, common.NewString(common.Float64ToString(e.TaxValue, 2)))
	ep.Add(TR_SHIPPING, common.NewString(common.Float64ToString(e.Shipping, 2)))
	ep.Add(TR_CITY, e.City)
	ep.Add(TR_STATE, e.State)
	ep.Add(TR_COUNTRY, e.Country)
	ep.Add(TR_CURRENCY, e.Currency)
	ep.Add(TIMESTAMP, common.NewString(common.Int64ToString(e.Timestamp)))
	ep.Add(EID, e.EventId)
	ep.Add(TRUE_TIMESTAMP, common.NewString(common.Int64ToString(e.TrueTimestamp)))
	if e.Subject != nil {
		ep.AddDict(e.Subject.Get())
	}
	return ep
}

// --- EcommerceTransactionItem Event

type EcommerceTransactionItemEvent struct {
	Sku      *string              // Required
	Price    *float64             // Required
	Quantity *int64               // Required
	Name     *string              // Optional
	Category *string              // Optional
	EventId  *string              // Optional
	Contexts []SelfDescribingJson // Optional
	Subject  *Subject             // Optional
}

// Init checks and validates the struct.
func (e *EcommerceTransactionItemEvent) Init() {
	if e.Sku == nil || *e.Sku == "" {
		panic("Sku cannot be nil or empty.")
	}
	if e.Price == nil {
		panic("Price cannot be nil.")
	}
	if e.Quantity == nil {
		panic("Quantity cannot be nil.")
	}
	if e.EventId == nil {
		e.EventId = common.NewString(common.GetUUID())
	}
}

// SetSubjectIfNil will set the event level subject if a Subject isn't already set
func (e *EcommerceTransactionItemEvent) SetSubjectIfNil(subject *Subject) {
	if e.Subject == nil {
		e.Subject = subject
	}
}

// Get returns the event payload.
func (e EcommerceTransactionItemEvent) Get() payload.Payload {
	ep := *payload.Init()
	ep.Add(EVENT, common.NewString(EVENT_ECOMM_ITEM))
	ep.Add(TI_ITEM_SKU, e.Sku)
	ep.Add(TI_ITEM_PRICE, common.NewString(common.Float64ToString(e.Price, 2)))
	ep.Add(TI_ITEM_QUANTITY, common.NewString(common.Int64ToString(e.Quantity)))
	ep.Add(TI_ITEM_NAME, e.Name)
	ep.Add(TI_ITEM_CATEGORY, e.Category)
	ep.Add(EID, e.EventId)
	if e.Subject != nil {
		ep.AddDict(e.Subject.Get())
	}
	return ep
}
