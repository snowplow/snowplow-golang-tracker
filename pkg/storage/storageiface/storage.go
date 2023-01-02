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

package storageiface

import (
	"github.com/snowplow/snowplow-golang-tracker/v3/pkg/payload"
)

const (
	DB_TABLE_NAME   = "events"
	DB_COLUMN_ID    = "id"
	DB_COLUMN_EVENT = "event"
)

type EventRow struct {
	Id    int
	Event payload.Payload
}

type Storage interface {
	AddEventRow(payload payload.Payload) bool
	DeleteAllEventRows() int64
	DeleteEventRows(ids []int) int64
	GetAllEventRows() []EventRow
	GetEventRowsWithinRange(eventRange int) []EventRow
}
