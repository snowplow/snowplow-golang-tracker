module github.com/snowplow/snowplow-golang-tracker/tracker

go 1.12

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hashicorp/go-immutable-radix v0.0.0-20180129170900-7f3cd4390caa
	github.com/hashicorp/go-memdb v0.0.0-20180223233045-1289e7fffe71
	github.com/hashicorp/golang-lru v0.0.0-20180201235237-0fb14efe8c47
	github.com/jarcoal/httpmock v0.0.0-20180424175123-9c70cfe4a1da
	github.com/mattn/go-sqlite3 v1.9.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/snowplow/snowplow-golang-tracker v2.0.0+incompatible // indirect
	github.com/stretchr/testify v1.2.2
	github.com/twinj/uuid v1.0.0
)

replace github.com/snowplow/snowplow-golang-tracker => ./
