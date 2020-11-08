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

	_ "github.com/mattn/go-sqlite3"
)

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
