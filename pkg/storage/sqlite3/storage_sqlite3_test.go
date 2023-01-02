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

package sqlite3

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/storageiface"
)

// TestStorageSQLite3Init asserts behaviour of SQLite storage functions.
func TestStorageSQLite3Init(t *testing.T) {
	assert := assert.New(t)
	storage := *Init("test.db")
	assert.NotNil(storage)
	assert.Equal("test.db", storage.DbName)
}

// TestSQLite3AddGetDeletePayload asserts ability to add, delete and get payloads.
func TestSQLite3AddGetDeletePayload(t *testing.T) {
	assert := assert.New(t)
	storage := *Init("test.db")
	assertDatabaseAddGetDeletePayload(assert, storage)
}

func TestSQLite3PanicRecovery(t *testing.T) {
	assert := assert.New(t)

	result := execDeleteQuery(nil, "")
	assert.Equal(int64(0), result)

	eventRows := execGetQuery(nil, "")
	assert.Equal(0, len(eventRows))

	addResult := execAddStatement(nil, nil)
	assert.False(addResult)
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
	assert.Equal(int64(1), storage.DeleteEventRows([]int{1}))
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
