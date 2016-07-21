//
// Copyright (c) 2016 Snowplow Analytics Ltd. All rights reserved.
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

const (
  TRACKER_VERSION = "golang-1.1.0"

  // POST Requests
  POST_PROTOCOL_VENDOR = "com.snowplowanalytics.snowplow"
  POST_PROTOCOL_VERSION = "tp2"
  POST_CONTENT_TYPE = "application/json; charset=utf-8"

  // GET Requests
  GET_PROTOCOL_PATH = "i"

  // Schema Versions
  SCHEMA_PAYLOAD_DATA = "iglu:com.snowplowanalytics.snowplow/payload_data/jsonschema/1-0-4"
  SCHEMA_CONTEXTS = "iglu:com.snowplowanalytics.snowplow/contexts/jsonschema/1-0-1"
  SCHEMA_UNSTRUCT_EVENT = "iglu:com.snowplowanalytics.snowplow/unstruct_event/jsonschema/1-0-0"
  SCHEMA_SCREEN_VIEW = "iglu:com.snowplowanalytics.snowplow/screen_view/jsonschema/1-0-0"
  SCHEMA_USER_TIMINGS = "iglu:com.snowplowanalytics.snowplow/timing/jsonschema/1-0-0"

  // Event Types
  EVENT_PAGE_VIEW = "pv"
  EVENT_STRUCTURED = "se"
  EVENT_UNSTRUCTURED = "ue"
  EVENT_ECOMM = "tr"
  EVENT_ECOMM_ITEM = "ti"

  // General
  SCHEMA = "schema"
  DATA = "data"
  EVENT = "e"
  EID = "eid"
  TIMESTAMP = "dtm"
  SENT_TIMESTAMP = "stm"
  TRUE_TIMESTAMP = "ttm"
  T_VERSION = "tv"
  APP_ID = "aid"
  NAMESPACE = "tna"
  PLATFORM = "p"

  CONTEXT = "co"
  CONTEXT_ENCODED = "cx"
  UNSTRUCTURED = "ue_pr"
  UNSTRUCTURED_ENCODED = "ue_px"

  // Subject class
  UID = "uid"
  RESOLUTION = "res"
  VIEWPORT = "vp"
  COLOR_DEPTH = "cd"
  TIMEZONE = "tz"
  LANGUAGE = "lang"
  IP_ADDRESS = "ip"
  USERAGENT = "ua"
  DOMAIN_UID = "duid"
  NETWORK_UID = "tnuid"

  // Page View
  PAGE_URL = "url"
  PAGE_TITLE = "page"
  PAGE_REFR = "refr"

  // Structured Event
  SE_CATEGORY = "se_ca"
  SE_ACTION = "se_ac"
  SE_LABEL = "se_la"
  SE_PROPERTY = "se_pr"
  SE_VALUE = "se_va"

  // Ecomm Transaction
  TR_ID = "tr_id"
  TR_TOTAL = "tr_tt"
  TR_AFFILIATION = "tr_af"
  TR_TAX = "tr_tx"
  TR_SHIPPING = "tr_sh"
  TR_CITY = "tr_ci"
  TR_STATE = "tr_st"
  TR_COUNTRY = "tr_co"
  TR_CURRENCY = "tr_cu"

  // Transaction Item
  TI_ITEM_ID = "ti_id"
  TI_ITEM_SKU = "ti_sk"
  TI_ITEM_NAME = "ti_nm"
  TI_ITEM_CATEGORY = "ti_ca"
  TI_ITEM_PRICE = "ti_pr"
  TI_ITEM_QUANTITY = "ti_qu"
  TI_ITEM_CURRENCY = "ti_cu"

  // Screen View
  SV_ID = "id"
  SV_NAME = "name"

  // User Timing
  UT_CATEGORY = "category"
  UT_VARIABLE = "variable"
  UT_TIMING = "timing"
  UT_LABEL = "label"
)
