package main

import (
	"fmt"
	"github.com/cubeee/ent-notifier/lib"
	"slices"
	"time"
)

func main() {
	env := lib.LoadEnv()

	now := time.Now().Unix()
	lastCheckTime := now - int64(env.HistoryLookupSeconds)
	pastEvents := make([]*lib.PastEvent, 0)

	for {
		resp, err := checkEventsLoop(env, lastCheckTime, pastEvents)
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
			return now-pastEvent.Time > int64(env.PastEventMaxAge)
		})

		lastCheckTime = resp.LatestEventTime
		time.Sleep(time.Duration(env.SleepTime) * time.Second)
	}
}

func checkEventsLoop(
	env *lib.Env,
	lastCheckTime int64,
	pastEvents []*lib.PastEvent,
) (*lib.EventsResponse, error) {
	fmt.Println("Checking events since", lastCheckTime, "- stored past events:", len(pastEvents))

	eventsResponse, err := lib.GetEvents(env, lastCheckTime, pastEvents)
	if err != nil {
		panic(err)
	}

	for _, newEvent := range eventsResponse.NewEvents {
		fmt.Println(newEvent.DiscoveredTime, newEvent.X, newEvent.Y, newEvent.Area)
	}

	if len(eventsResponse.NewEvents) > 0 {
		fmt.Println("\tNew events:", len(eventsResponse.NewEvents))
		err := lib.NotifyEvents(env, eventsResponse.NewEvents, env.WebhookUrls)
		if err != nil {
			return nil, err
		}
	}

	return eventsResponse, nil
}
