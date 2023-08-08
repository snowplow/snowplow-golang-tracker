# Example App

This example shows how to set up the golang tracker to send events to a Snowplow pipeline.

#### Installation
- Install the golang tracker.

`$host go get "github.com/snowplow/snowplow-golang-tracker/v3/tracker"` 

#### Usage
Navigate to the example folder.

`cd examples`

To send events to your pipeline, run `go run main.go {{your_collector_endpoint}}`. You should see 4 events in your pipleine.

