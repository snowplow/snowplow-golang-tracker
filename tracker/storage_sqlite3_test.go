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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStorageSQLite3Init asserts behaviour of SQLite storage functions.
func TestStorageSQLite3Init(t *testing.T) {
	assert := assert.New(t)
	storage := *InitStorageSQLite3("test.db")
	assert.NotNil(storage)
	assert.Equal("test.db", storage.DbName)

	defer func() {
		if err := recover(); err != nil {
			assert.NotNil(err)
		}
	}()
	storage = *InitStorageSQLite3("~/")
}

// TestSQLite3AddGetDeletePayload asserts ability to add, delete and get payloads.
func TestSQLite3AddGetDeletePayload(t *testing.T) {
	assert := assert.New(t)
	storage := *InitStorageSQLite3("test.db")
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
