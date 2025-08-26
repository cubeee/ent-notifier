package lib

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
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
	World          int
	DiscoveredTime int64
	X              int
	Y              int
	Area           *Area
	MappedLocation *MappedLocation
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

type MappedLocation struct {
	Name      string   `json:"name"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	Radius    int      `json:"radius"`
	Teleports []string `json:"teleports"`
}

func GetEvents(
	env *Env,
	lastCheckTime int64,
	mappedLocations []*MappedLocation,
	pastEvents []*PastEvent,
) (*EventsResponse, error) {
	response, err := getEvents(env.ApiUrl, env.ApiTimeout, env.ApiUserAgent)
	if err != nil {
		return nil, err
	}

	var events []*Event
	for _, eventItem := range response.Items {
		if eventItem.DiscoveredTime < lastCheckTime {
			continue
		}

		if eventItem.World != env.EventsWorld {
			continue
		}

		if !slices.Contains(env.EventsAllowed, eventItem.EventType) {
			continue
		}

		event := &Event{
			EventType:      eventItem.EventType,
			World:          eventItem.World,
			DiscoveredTime: eventItem.DiscoveredTime,
			X:              eventItem.X,
			Y:              eventItem.Y,
			Area:           CreateEventArea(eventItem.X, eventItem.Y, env.EventAreaRadius),
			MappedLocation: GetMappedLocation(eventItem.X, eventItem.Y, mappedLocations),
		}

		if overlapsEvent(event, events, env.LocationCooldown) {
			continue
		}

		if overlapsPastEvent(event, pastEvents, env.LocationCooldown) {
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

func overlapsEvent(event *Event, otherEvents []*Event, locationCooldown int) bool {
	for i, _ := range otherEvents {
		other := *otherEvents[i]
		if event.EventType == other.EventType && (event.DiscoveredTime-other.DiscoveredTime < int64(locationCooldown)) && event.Area.IntersectsArea(other.Area) {
			return true
		}
	}
	return false
}

func overlapsPastEvent(event *Event, otherEvents []*PastEvent, locationCooldown int) bool {
	for i, _ := range otherEvents {
		other := *otherEvents[i]
		if event.EventType == other.EventType && (event.DiscoveredTime-other.Time < int64(locationCooldown)) && event.Area.IntersectsArea(other.Area) {
			return true
		}
	}
	return false
}

func GetMappedLocation(x int, y int, mappedLocations []*MappedLocation) *MappedLocation {
	var closestLocation *MappedLocation
	closestDistance := 1_000_000.0

	for _, location := range mappedLocations {
		xDiff := math.Abs(float64(location.X - x))
		yDiff := math.Abs(float64(location.Y - y))
		avg := (xDiff + yDiff) / 2
		if avg < closestDistance && avg < float64(location.Radius) {
			closestLocation = location
			closestDistance = math.Min(xDiff, yDiff)
		}
	}
	return closestLocation
}

func getEvents(url string, timeout int, userAgent string) (*EventsApiResponse, error) {
	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if len(userAgent) > 0 {
		req.Header.Add("User-Agent", userAgent)
	}
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
