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
	"database/sql"
	"log"
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_DRIVER       = "sqlite3"
	DB_TABLE_NAME   = "events"
	DB_COLUMN_ID    = "id"
	DB_COLUMN_EVENT = "event"
)

type Storage interface {
	AddEventRow(payload Payload) bool
	DeleteAllEventRows() int64
	DeleteEventRows(ids []int) int64
	GetAllEventRows() []EventRow
	GetEventRowsWithinRange(eventRange int) []EventRow
}

type RawEventRow struct {
	id    int
	event []byte
}

type RawEventRowUint struct {
	id    uint
	event []byte
}

type EventRow struct {
	id    int
	event Payload
}

// --- Memory Storage Implementation

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

// --- SQLite3 Storage Implementation

type StorageSQLite3 struct {
	DbName string
}

func InitStorageSQLite3(dbName string) *StorageSQLite3 {
	db := getDbConn(dbName)
	defer db.Close()

	db.SetMaxOpenConns(1)

	// Enable Write-Ahead-Logging for concurrent read and write
	_, err1 := db.Exec("PRAGMA journal_mode=WAL;")
	checkErr(err1)

	// Create the Events Table
	query :=
		"CREATE TABLE IF NOT EXISTS " + DB_TABLE_NAME + "(" +
			DB_COLUMN_ID + " INTEGER PRIMARY KEY, " +
			DB_COLUMN_EVENT + " BLOB" +
			");"
	_, err2 := db.Exec(query)
	checkErr(err2)

	return &StorageSQLite3{DbName: dbName}
}

func getDbConn(dbName string) *sql.DB {
	db, err := sql.Open(DB_DRIVER, dbName)
	checkErr(err)
	return db
}

// --- ADD

// Add stores an event payload in the database.
func (s StorageSQLite3) AddEventRow(payload Payload) bool {
	db := getDbConn(s.DbName)
	defer db.Close()

	// Prepare Add Statement
	query :=
		"INSERT INTO " + DB_TABLE_NAME + "(" +
			DB_COLUMN_EVENT +
			") values(?);"
	addStmt, err1 := db.Prepare(query)
	checkErr(err1)

	byteBuffer := SerializeMap(payload.Get())
	return execAddStatement(addStmt, byteBuffer)
}

// execAddStatement executes the add statement passed to it.
func execAddStatement(stmt *sql.Stmt, byteBuffer []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	res, err := stmt.Exec(byteBuffer)
	checkErr(err)
	affected, err2 := res.RowsAffected()
	checkErr(err2)

	return affected == 1
}

// --- DELETE

// DeleteAllEventRows removes all events from the database.
func (s StorageSQLite3) DeleteAllEventRows() int64 {
	db := getDbConn(s.DbName)
	defer db.Close()

	query := "DELETE FROM " + DB_TABLE_NAME + ";"
	return execDeleteQuery(db, query)
}

// DeleteEventRows removes a range of ids from the database.
func (s StorageSQLite3) DeleteEventRows(ids []int) int64 {
	db := getDbConn(s.DbName)
	defer db.Close()

	if len(ids) > 0 {
		query :=
			"DELETE FROM " + DB_TABLE_NAME + " " +
				"WHERE " + DB_COLUMN_ID + " in(" + IntArrayToString(ids, ",") + ");"
		return execDeleteQuery(db, query)
	} else {
		return 0
	}
}

// execDeleteQuery is used to run queries which removed event rows from the database.
func execDeleteQuery(db *sql.DB, query string) int64 {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	stmt, err := db.Prepare(query)
	checkErr(err)
	defer stmt.Close()
	res, err2 := stmt.Exec()
	checkErr(err2)
	affected, err3 := res.RowsAffected()
	checkErr(err3)

	return affected
}

// --- GET

// GetAllEventRows returns all events in the database.
func (s StorageSQLite3) GetAllEventRows() []EventRow {
	db := getDbConn(s.DbName)
	defer db.Close()

	query := "SELECT " + DB_COLUMN_ID + ", " + DB_COLUMN_EVENT + " FROM " + DB_TABLE_NAME + ";"
	return execGetQuery(db, query)
}

// GetEventRowsWithinRange returns a specified range of events from the database.
func (s StorageSQLite3) GetEventRowsWithinRange(eventRange int) []EventRow {
	db := getDbConn(s.DbName)
	defer db.Close()

	query :=
		"SELECT " + DB_COLUMN_ID + ", " + DB_COLUMN_EVENT + " FROM " + DB_TABLE_NAME + " " +
			"ORDER BY " + DB_COLUMN_ID + " DESC LIMIT " + IntToString(eventRange) + ";"
	return execGetQuery(db, query)
}

// execGetQuery is used to run queries to fetch event rows from the database.
func execGetQuery(db *sql.DB, query string) []EventRow {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	eventItems := []EventRow{}
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		item := RawEventRow{}
		rows.Scan(&item.id, &item.event)
		eventMap, _ := DeserializeMap(item.event)
		eventItems = append(eventItems, EventRow{item.id, Payload{eventMap}})
	}

	return eventItems
}

// --- Helpers

// checkErr throws a panic for all non-nil errors passed to it.
func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
