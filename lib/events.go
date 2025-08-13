package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
)

type PastEvent struct {
	EventType string
	Time      int64
	Area      *Area
}

type Event struct {
	EventType      string
	DiscoveredTime int64
	X              int
	Y              int
	Area           *Area
}

type EventsResponse struct {
	NewEvents       []*Event
	LatestEventTime int64
}

type EventsApiEvent struct {
	EventType      string `json:"event_type"`
	World          int    `json:"world"`
	DiscoveredTime int64  `json:"discovered_time"`
	X              int    `json:"x_coord"`
	Y              int    `json:"y_coord"`
}

type EventsApiResponse struct {
	Items []EventsApiEvent `json:"items"`
}

var (
	ApiUrl           = os.Getenv("EVENTS_API_URL")
	ApiTimeout       = GetEnvInt("EVENTS_API_TIMEOUT", 5)
	EventsWorld      = GetEnvInt("EVENTS_WORLD", 444)
	EventsAllowed    = GetEnvList("EVENTS_ALLOWED", ",")
	UserAgent        = os.Getenv("USER_AGENT")
	LocationCooldown = GetEnvInt("LOCATION_COOLDOWN", 300)
)

func GetEvents(lastCheckTime int64, pastEvents []*PastEvent) (*EventsResponse, error) {
	response, err := getEvents(ApiUrl)
	if err != nil {
		return nil, err
	}

	var events []*Event
	for _, eventItem := range response.Items {
		if eventItem.DiscoveredTime < lastCheckTime {
			continue
		}

		if eventItem.World != EventsWorld {
			continue
		}

		if !slices.Contains(EventsAllowed, eventItem.EventType) {
			continue
		}

		event := &Event{
			EventType:      eventItem.EventType,
			DiscoveredTime: eventItem.DiscoveredTime,
			X:              eventItem.X,
			Y:              eventItem.Y,
			Area:           CreateEventArea(eventItem.X, eventItem.Y),
		}

		if overlapsEvent(event, events) {
			continue
		}

		if overlapsPastEvent(event, pastEvents) {
			continue
		}

		events = append(events, event)
	}

	var latestEventTime int64
	for _, eventItem := range response.Items {
		if eventItem.DiscoveredTime > latestEventTime {
			latestEventTime = eventItem.DiscoveredTime
		}
	}

	return &EventsResponse{
		NewEvents:       events,
		LatestEventTime: latestEventTime,
	}, nil
}

func overlapsEvent(event *Event, otherEvents []*Event) bool {
	for i, _ := range otherEvents {
		other := *otherEvents[i]
		if event.EventType == other.EventType && (event.DiscoveredTime-other.DiscoveredTime < int64(LocationCooldown)) && event.Area.IntersectsArea(other.Area) {
			return true
		}
	}
	return false
}

func overlapsPastEvent(event *Event, otherEvents []*PastEvent) bool {
	for i, _ := range otherEvents {
		other := *otherEvents[i]
		if event.EventType == other.EventType && (event.DiscoveredTime-other.Time < int64(LocationCooldown)) && event.Area.IntersectsArea(other.Area) {
			return true
		}
	}
	return false
}

func getEvents(url string) (*EventsApiResponse, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(ApiTimeout),
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				panic(err)
			}
		}(res.Body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	events := EventsApiResponse{}
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}

	return &events, nil
}
