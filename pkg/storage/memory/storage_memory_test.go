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

package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/storageiface"
)

// TestStorageMemoryInit asserts behaviour of memdb storage functions.
func TestStorageMemoryInit(t *testing.T) {
	assert := assert.New(t)
	storage := Init()
	assert.NotNil(storage)
	assert.NotNil(storage.Db)
}

// TestMemoryAddGetDeletePayload asserts ability to add, delete and get payloads.
func TestMemoryAddGetDeletePayload(t *testing.T) {
	assert := assert.New(t)
	storage := *Init()
	assertDatabaseAddGetDeletePayload(assert, storage)
}

// TestMemoryAddGetDeletePayload_WithIndexOverflow asserts behaviour when index overflows beyond UInt32 max.
func TestMemoryAddGetDeletePayload_WithIndexOverflow(t *testing.T) {
	assert := assert.New(t)
	storage := *Init()

	// Set to max value available for Index
	var index uint32 = 4294967295
	storage.Index = &index

	assertDatabaseAddGetDeletePayload(assert, storage)
}

// TestMemoryAddPayload_WithUniqueIndexCollision asserts behaviour when we hit a collision in the database which can only happen
// if we rollover the Index and have not yet deleted the events (creating faster than sending).
func TestMemoryAddPayload_WithUniqueIndexCollision(t *testing.T) {
	assert := assert.New(t)
	storage := *Init()
	assert.Equal(uint32(0), *storage.Index)

	payload := *payload.Init()
	payload.Add("e", common.NewString("pv"))

	// Add 20 payloads
	for i := 0; i < 20; i++ {
		result := storage.AddEventRow(payload)
		assert.True(result)
	}

	assert.Equal(uint32(20), *storage.Index)
	eventRows := storage.GetAllEventRows()
	assert.Equal(20, len(eventRows))

	// Reset index to 0
	var index uint32 = 0
	storage.Index = &index
	assert.Equal(uint32(0), *storage.Index)

	// TODO: Ideally this should fail but doesn't due to this - https://github.com/hashicorp/go-memdb/issues/7

	// Add 20 payloads
	for i := 0; i < 20; i++ {
		result := storage.AddEventRow(payload)
		assert.True(result)
	}

	assert.Equal(uint32(20), *storage.Index)
	eventRows2 := storage.GetAllEventRows()
	assert.Equal(20, len(eventRows2))
}

// --- Common

func assertDatabaseAddGetDeletePayload(assert *assert.Assertions, storage storageiface.Storage) {
	storage.DeleteAllEventRows()
	payload := *payload.Init()
	payload.Add("e", common.NewString("pv"))

	// Add a Payload
	assert.True(storage.AddEventRow(payload))
	eventRows := storage.GetAllEventRows()
	assert.Equal(1, len(eventRows))
	assert.Equal("pv", eventRows[0].Event.Get()["e"])

	// Delete the added row
	assert.Equal(int64(1), storage.DeleteEventRows([]int{eventRows[0].Id}))
	eventRows = storage.GetAllEventRows()
	assert.Equal(0, len(eventRows))

	// Add 20 payloads
	for i := 0; i < 20; i++ {
		result := storage.AddEventRow(payload)
		assert.True(result)
	}

	eventRows = storage.GetEventRowsWithinRange(10)
	assert.Equal(10, len(eventRows))
	eventRows = storage.GetEventRowsWithinRange(30)
	assert.Equal(20, len(eventRows))
	eventRows = storage.GetAllEventRows()
	assert.Equal(20, len(eventRows))
	assert.Equal(int64(20), storage.DeleteAllEventRows())
	eventRows = storage.GetAllEventRows()
	assert.Equal(0, len(eventRows))
	assert.Equal(int64(0), storage.DeleteEventRows([]int{}))
}
