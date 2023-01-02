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
	"sync/atomic"

	"github.com/hashicorp/go-memdb"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/storageiface"
)

type StorageMemory struct {
	Db    *memdb.MemDB
	Index *uint32
}

type RawEventRowUint struct {
	id    uint
	event []byte
}

func Init() *StorageMemory {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			storageiface.DB_TABLE_NAME: {
				Name: storageiface.DB_TABLE_NAME,
				Indexes: map[string]*memdb.IndexSchema{
					storageiface.DB_COLUMN_ID: {
						Name:    storageiface.DB_COLUMN_ID,
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: storageiface.DB_COLUMN_ID},
					},
					storageiface.DB_COLUMN_EVENT: {
						Name:    storageiface.DB_COLUMN_EVENT,
						Indexer: &memdb.StringFieldIndex{Field: storageiface.DB_COLUMN_EVENT},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	common.CheckErr(err)

	return &StorageMemory{Db: db, Index: new(uint32)}
}

// AddEventRow adds a new event to the database
//
// NOTE: As entries are not auto-incremeneting the id is incremented manually which
//
//	limits inserts to 4,294,967,295 in single session
func (s StorageMemory) AddEventRow(payload payload.Payload) bool {
	txn := s.Db.Txn(true)
	byteBuffer := common.SerializeMap(payload.Get())
	rer := &RawEventRowUint{event: byteBuffer, id: uint(atomic.AddUint32(s.Index, 1))}
	err := txn.Insert(storageiface.DB_TABLE_NAME, rer)
	common.CheckErr(err)
	txn.Commit()

	return true
}

// DeleteAllEventRows removes all rows within the memory store
func (s StorageMemory) DeleteAllEventRows() int64 {
	txn := s.Db.Txn(true)
	result, err := txn.DeleteAll(storageiface.DB_TABLE_NAME, storageiface.DB_COLUMN_ID)
	common.CheckErr(err)
	txn.Commit()

	return int64(result)
}

// DeleteEventRows removes all rows with matching identifiers
func (s StorageMemory) DeleteEventRows(ids []int) int64 {
	txn := s.Db.Txn(true)
	deleteCount := 0

	for _, id := range ids {
		result, err := txn.DeleteAll(storageiface.DB_TABLE_NAME, storageiface.DB_COLUMN_ID, uint(id))
		common.CheckErr(err)
		deleteCount += result
	}

	txn.Commit()

	return int64(deleteCount)
}

// GetAllEventRows returns all rows within the memory store
func (s StorageMemory) GetAllEventRows() []storageiface.EventRow {
	eventItems := []storageiface.EventRow{}
	txn := s.Db.Txn(false)
	defer txn.Abort()

	result, err := txn.Get(storageiface.DB_TABLE_NAME, storageiface.DB_COLUMN_ID)
	common.CheckErr(err)
	for row := result.Next(); row != nil; row = result.Next() {
		item := row.(*RawEventRowUint)
		eventMap, _ := common.DeserializeMap(item.event)
		eventItems = append(eventItems, storageiface.EventRow{int(item.id), payload.Payload{eventMap}})
	}

	return eventItems
}

// GetEventRowsWithinRange returns all available events or a maximal slice
func (s StorageMemory) GetEventRowsWithinRange(eventRange int) []storageiface.EventRow {
	eventItems := s.GetAllEventRows()
	if len(eventItems) <= eventRange {
		return eventItems
	} else {
		return eventItems[:eventRange]
	}
}
