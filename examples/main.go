package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	sphelp "github.com/snowplow/snowplow-golang-tracker/v3/pkg/common"
	storagememory "github.com/snowplow/snowplow-golang-tracker/v3/pkg/storage/memory"
	sp "github.com/snowplow/snowplow-golang-tracker/v3/tracker"
)

func getUrlFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", errors.New("collector endpoint is required")
	}
	return os.Args[1], nil
}

func main() {
	collectorUri, err := getUrlFromArgs()

	if err != nil {
		log.Fatal(err)
	}

	emitter := sp.InitEmitter(
		sp.RequireCollectorUri(collectorUri),
		sp.RequireStorage(*storagememory.Init()),
		sp.OptionRequestType("POST"),
		sp.OptionProtocol("http"),
		sp.OptionSendLimit(4),
	)

	subject := sp.InitSubject()
	subject.SetLanguage("en")
	subject.SetScreenResolution(1280, 720)

	tracker := sp.InitTracker(
		sp.RequireEmitter(emitter),
		sp.OptionSubject(subject),
	)

	fmt.Println("Sending events to " + emitter.GetCollectorUrl())

	pageView := sp.PageViewEvent{
		PageUrl: sphelp.NewString("acme.com"),
	}
	tracker.TrackPageView(pageView)

	screenView := sp.ScreenViewEvent{
		Name: sphelp.NewString("name"),
		Id:   sphelp.NewString("Screen ID"),
	}
	tracker.TrackScreenView(screenView)

	structEvent := sp.StructuredEvent{
		Category: sphelp.NewString("shop"),
		Action:   sphelp.NewString("add-to-basket"),
		Property: sphelp.NewString("pcs"),
		Value:    sphelp.NewFloat64(2),
	}

	tracker.TrackStructEvent(structEvent)

	data := map[string]interface{}{
		"targetUrl": "https://www.snowplow.io",
	}

	sdj := sp.InitSelfDescribingJson("iglu:com.snowplowanalytics.snowplow/link_click/jsonschema/1-0-1", data)

	sde := sp.SelfDescribingEvent{Event: sdj}

	tracker.TrackSelfDescribingEvent(sde)

	tracker.Emitter.Stop()
	tracker.BlockingFlush(5, 10)

}
