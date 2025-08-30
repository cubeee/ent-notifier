package main

import (
	"encoding/json"
	"github.com/cubeee/ent-notifier/lib"
	"io"
	"log"
	"os"
	"slices"
	"time"
)

func main() {
	env := lib.LoadEnv()

	now := time.Now().Unix()
	lastCheckTime := now - int64(env.HistoryLookupSeconds)
	pastEvents := make([]*lib.PastEvent, 0)
	mappedLocations := loadMappedLocations(env.MappedLocationsFile)
	log.Println("Loaded", len(mappedLocations), "mapped locations")

	for {
		resp, err := checkEventsLoop(env, lastCheckTime, mappedLocations, pastEvents)
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
	mappedLocations []*lib.MappedLocation,
	pastEvents []*lib.PastEvent,
) (*lib.EventsResponse, error) {
	log.Println("Checking events since", lastCheckTime, "- stored past events:", len(pastEvents))

	eventsResponse, err := lib.GetEvents(env, lastCheckTime, mappedLocations, pastEvents)
	if err != nil {
		panic(err)
	}

	for _, newEvent := range eventsResponse.NewEvents {
		log.Println(newEvent.DiscoveredTime, newEvent.X, newEvent.Y, newEvent.Area)
		if newEvent.MappedLocation == nil {
			log.Printf("\tunmapped location: %d,%d", newEvent.X, newEvent.Y)
		}
	}

	if len(eventsResponse.NewEvents) > 0 {
		log.Println("\tNew events:", len(eventsResponse.NewEvents))
		err := lib.NotifyEvents(env, eventsResponse.NewEvents, env.WebhookUrls)
		if err != nil {
			return nil, err
		}
	}

	return eventsResponse, nil
}

func loadMappedLocations(path string) []*lib.MappedLocation {
	var locations []*lib.MappedLocation
	if len(path) == 0 {
		log.Println("Mapped locations file path not set")
		return locations
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println("failed to read mapped locations file:", err)
		return locations
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Panicln("failed to close mapped locations file:", err)
		}
	}(jsonFile)

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println("failed to read mapped locations file contents:", err)
		return locations
	}

	err = json.Unmarshal(byteValue, &locations)
	if err != nil {
		log.Println("failed to unmarshal mapped locations:", err)
		return locations
	}

	return locations
}
