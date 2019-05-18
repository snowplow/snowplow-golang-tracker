package main

import (
	"fmt"
	sp "github.com/snowplow/snowplow-golang-tracker/tracker"
	"log"
	"net/http"
	"strconv"
	"time"
)

var client *http.Client

func startWebserver() {
	srv := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/com.snowplowanalytics.snowplow/tp2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
}

func startLoadWorker(number int) {
	count := 0

	emitter := sp.InitEmitter(
		sp.RequireCollectorUri("127.0.0.1:8080"),
		sp.OptionRequestType("POST"),
		sp.OptionDbName("event" + strconv.Itoa(number) + ".db"),
		sp.OptionHttpClient(client),
	)

	tracker := sp.InitTracker(
		sp.RequireEmitter(emitter),
	)

	data := map[string]interface{}{
		"level": 5,
		"saveId": "ju302",
		"hardMode": true,
	}
	sdj := sp.InitSelfDescribingJson("iglu:com.example_company/save-game/jsonschema/1-0-2", data)
	for {
		time.Sleep(time.Millisecond * time.Duration(10))

		tracker.TrackSelfDescribingEvent(sp.SelfDescribingEvent{
			Event: sdj,
		})
		log.Printf("Tracker #%v finished POST request #%v", number, count)
		count += 1
	}

}

func main() {
	startWebserver()

	// Setup HttpClient
	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 1024
	defaultTransport.MaxIdleConnsPerHost = 1024

	client = &http.Client{Transport: &defaultTransport}

	for n := 0; n < 100; n++ {
		go startLoadWorker(n)
	}

	time.Sleep(time.Second * 2400)
}
