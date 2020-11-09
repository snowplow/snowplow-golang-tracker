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
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
)

type StorageMemory struct {
	Db    *memdb.MemDB
	Index *uint32
}

func InitStorageMemory() *StorageMemory {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			DB_TABLE_NAME: {
				Name: DB_TABLE_NAME,
				Indexes: map[string]*memdb.IndexSchema{
					DB_COLUMN_ID: {
						Name:    DB_COLUMN_ID,
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: DB_COLUMN_ID},
					},
					DB_COLUMN_EVENT: {
						Name:    DB_COLUMN_EVENT,
						Indexer: &memdb.StringFieldIndex{Field: DB_COLUMN_EVENT},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	checkErr(err)

	return &StorageMemory{Db: db, Index: new(uint32)}
}

// AddEventRow adds a new event to the database
//
// NOTE: As entries are not auto-incremeneting the id is incremented manually which
//       limits inserts to 4,294,967,295 in single session
func (s StorageMemory) AddEventRow(payload Payload) bool {
	txn := s.Db.Txn(true)
	byteBuffer := SerializeMap(payload.Get())
	rer := &RawEventRowUint{event: byteBuffer, id: uint(atomic.AddUint32(s.Index, 1))}
	err := txn.Insert(DB_TABLE_NAME, rer)
	checkErr(err)
	txn.Commit()

	return true
}

// DeleteAllEventRows removes all rows within the memory store
func (s StorageMemory) DeleteAllEventRows() int64 {
	txn := s.Db.Txn(true)
	result, err := txn.DeleteAll(DB_TABLE_NAME, DB_COLUMN_ID)
	checkErr(err)
	txn.Commit()

	return int64(result)
}

// DeleteEventRows removes all rows with matching identifiers
func (s StorageMemory) DeleteEventRows(ids []int) int64 {
	txn := s.Db.Txn(true)
	deleteCount := 0

	for _, id := range ids {
		result, err := txn.DeleteAll(DB_TABLE_NAME, DB_COLUMN_ID, uint(id))
		checkErr(err)
		deleteCount += result
	}

	txn.Commit()

	return int64(deleteCount)
}

// GetAllEventRows returns all rows within the memory store
func (s StorageMemory) GetAllEventRows() []EventRow {
	eventItems := []EventRow{}
	txn := s.Db.Txn(false)
	defer txn.Abort()

	result, err := txn.Get(DB_TABLE_NAME, DB_COLUMN_ID)
	checkErr(err)
	for row := result.Next(); row != nil; row = result.Next() {
		item := row.(*RawEventRowUint)
		eventMap, _ := DeserializeMap(item.event)
		eventItems = append(eventItems, EventRow{int(item.id), Payload{eventMap}})
	}

	return eventItems
}

// GetEventRowsWithinRange returns all available events or a maximal slice
func (s StorageMemory) GetEventRowsWithinRange(eventRange int) []EventRow {
	eventItems := s.GetAllEventRows()
	if len(eventItems) <= eventRange {
		return eventItems
	} else {
		return eventItems[:eventRange]
	}
}
