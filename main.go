package main

import (
	"ent-notifier/lib"
	"fmt"
	"slices"
	"time"
)

var (
	SleepTime            = lib.GetEnvInt("SLEEP_TIME_SECONDS", 5)
	HistoryLookupSeconds = lib.GetEnvInt("HISTORY_LOOKUP_SECONDS", 15)
	WebHooksUrl          = lib.GetEnvList("DISCORD_WEBHOOK_URLS", ",")
	PastEventMaxAge      = lib.GetEnvInt("PAST_EVENT_MAX_AGE", 30*60)
)

func main() {
	now := time.Now().Unix()
	lastCheckTime := now - int64(HistoryLookupSeconds)
	pastEvents := make([]*lib.PastEvent, 0)

	for {
		resp, err := checkEventsLoop(lastCheckTime, pastEvents)
		if err != nil {
			panic(err)
		}

		for _, newEvent := range resp.NewEvents {
			pastEvents = append(pastEvents, &lib.PastEvent{
				EventType: newEvent.EventType,
				Time:      newEvent.DiscoveredTime,
				Area:      newEvent.Area,
			})
		}

		now = time.Now().Unix()
		pastEvents = slices.DeleteFunc(pastEvents, func(pastEvent *lib.PastEvent) bool {
			return now-pastEvent.Time > int64(PastEventMaxAge)
		})

		lastCheckTime = resp.LatestEventTime
		time.Sleep(time.Duration(SleepTime) * time.Second)
	}
}

func checkEventsLoop(lastCheckTime int64, pastEvents []*lib.PastEvent) (*lib.EventsResponse, error) {
	fmt.Println("Checking events since", lastCheckTime, "- stored past events:", len(pastEvents))

	eventsResponse, err := lib.GetEvents(lastCheckTime, pastEvents)
	if err != nil {
		panic(err)
	}

	for _, newEvent := range eventsResponse.NewEvents {
		fmt.Println(newEvent.DiscoveredTime, newEvent.X, newEvent.Y, newEvent.Area)
	}

	if len(eventsResponse.NewEvents) > 0 {
		fmt.Println("\tNew events:", len(eventsResponse.NewEvents))
		err := lib.NotifyEvents(eventsResponse.NewEvents, WebHooksUrl)
		if err != nil {
			return nil, err
		}
	}

	return eventsResponse, nil
}
