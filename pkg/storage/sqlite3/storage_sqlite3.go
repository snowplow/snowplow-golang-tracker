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
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/storageiface"
)

type StorageSQLite3 struct {
	DbName string
}

type RawEventRow struct {
	id    int
	event []byte
}

func Init(dbName string) *StorageSQLite3 {
	db := getDbConn(dbName)
	defer db.Close()

	db.SetMaxOpenConns(1)

	// Enable Write-Ahead-Logging for concurrent read and write
	_, err1 := db.Exec("PRAGMA journal_mode=WAL;")
	common.CheckErr(err1)

	// Create the Events Table
	query :=
		"CREATE TABLE IF NOT EXISTS " + storageiface.DB_TABLE_NAME + "(" +
			storageiface.DB_COLUMN_ID + " INTEGER PRIMARY KEY, " +
			storageiface.DB_COLUMN_EVENT + " BLOB" +
			");"
	_, err2 := db.Exec(query)
	common.CheckErr(err2)

	return &StorageSQLite3{DbName: dbName}
}

func getDbConn(dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	common.CheckErr(err)
	return db
}

// --- ADD

// Add stores an event payload in the database.
func (s StorageSQLite3) AddEventRow(payload payload.Payload) bool {
	db := getDbConn(s.DbName)
	defer db.Close()

	// Prepare Add Statement
	query :=
		"INSERT INTO " + storageiface.DB_TABLE_NAME + "(" +
			storageiface.DB_COLUMN_EVENT +
			") values(?);"
	addStmt, err1 := db.Prepare(query)
	common.CheckErr(err1)

	byteBuffer := common.SerializeMap(payload.Get())
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
	common.CheckErr(err)
	affected, err2 := res.RowsAffected()
	common.CheckErr(err2)

	return affected == 1
}

// --- DELETE

// DeleteAllEventRows removes all events from the database.
func (s StorageSQLite3) DeleteAllEventRows() int64 {
	db := getDbConn(s.DbName)
	defer db.Close()

	query := "DELETE FROM " + storageiface.DB_TABLE_NAME + ";"
	return execDeleteQuery(db, query)
}

// DeleteEventRows removes a range of ids from the database.
func (s StorageSQLite3) DeleteEventRows(ids []int) int64 {
	db := getDbConn(s.DbName)
	defer db.Close()

	if len(ids) > 0 {
		query :=
			"DELETE FROM " + storageiface.DB_TABLE_NAME + " " +
				"WHERE " + storageiface.DB_COLUMN_ID + " in(" + common.IntArrayToString(ids, ",") + ");"
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
	common.CheckErr(err)
	defer stmt.Close()
	res, err2 := stmt.Exec()
	common.CheckErr(err2)
	affected, err3 := res.RowsAffected()
	common.CheckErr(err3)

	return affected
}

// --- GET

// GetAllEventRows returns all events in the database.
func (s StorageSQLite3) GetAllEventRows() []storageiface.EventRow {
	db := getDbConn(s.DbName)
	defer db.Close()

	query := "SELECT " + storageiface.DB_COLUMN_ID + ", " + storageiface.DB_COLUMN_EVENT + " FROM " + storageiface.DB_TABLE_NAME + ";"
	return execGetQuery(db, query)
}

// GetEventRowsWithinRange returns a specified range of events from the database.
func (s StorageSQLite3) GetEventRowsWithinRange(eventRange int) []storageiface.EventRow {
	db := getDbConn(s.DbName)
	defer db.Close()

	query :=
		"SELECT " + storageiface.DB_COLUMN_ID + ", " + storageiface.DB_COLUMN_EVENT + " FROM " + storageiface.DB_TABLE_NAME + " " +
			"ORDER BY " + storageiface.DB_COLUMN_ID + " DESC LIMIT " + common.IntToString(eventRange) + ";"
	return execGetQuery(db, query)
}

// execGetQuery is used to run queries to fetch event rows from the database.
func execGetQuery(db *sql.DB, query string) []storageiface.EventRow {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	eventItems := []storageiface.EventRow{}
	rows, err := db.Query(query)
	common.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		item := RawEventRow{}
		rows.Scan(&item.id, &item.event)
		eventMap, _ := common.DeserializeMap(item.event)
		eventItems = append(eventItems, storageiface.EventRow{item.id, payload.Payload{eventMap}})
	}

	return eventItems
}
