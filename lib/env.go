package lib

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	EventAreaRadius      int
	MapImageWidth        int
	MapImageHeight       int
	EmbedFooter          string
	ApiUrl               string
	ApiTimeout           int
	ApiUserAgent         string
	EventsWorld          int
	EventsAllowed        []string
	LocationCooldown     int
	MapFilePath          string
	MapTilePixels        int
	SleepTime            int
	HistoryLookupSeconds int
	WebhookUrls          []string
	PastEventMaxAge      int
	MappedLocationsFile  string
}

func LoadEnv() *Env {
	return &Env{
		EventAreaRadius:      GetEnvInt("EVENT_AREA_RADIUS", 8),
		MapImageWidth:        GetEnvInt("MAP_IMAGE_WIDTH", 400),
		MapImageHeight:       GetEnvInt("MAP_IMAGE_HEIGHT", 250),
		EmbedFooter:          GetEnv("EMBED_FOOTER", ""),
		ApiUrl:               GetRequiredEnv("EVENTS_API_URL"),
		ApiTimeout:           GetEnvInt("EVENTS_API_TIMEOUT", 5),
		ApiUserAgent:         os.Getenv("EVENTS_API_USER_AGENT"),
		EventsWorld:          GetEnvInt("EVENTS_WORLD", 444),
		EventsAllowed:        GetEnvList("EVENTS_ALLOWED", ","),
		LocationCooldown:     GetEnvInt("LOCATION_COOLDOWN", 300),
		MapFilePath:          GetEnv("MAP_FILE_PATH", "layers_osrs/mapsquares/-1/2"),
		MapTilePixels:        GetEnvInt("MAP_TILE_PIXELS", 4),
		SleepTime:            GetEnvInt("SLEEP_TIME_SECONDS", 5),
		HistoryLookupSeconds: GetEnvInt("HISTORY_LOOKUP_SECONDS", 0),
		WebhookUrls:          GetEnvList("DISCORD_WEBHOOK_URLS", ","),
		PastEventMaxAge:      GetEnvInt("PAST_EVENT_MAX_AGE", 30*60),
		MappedLocationsFile:  os.Getenv("MAPPED_LOCATIONS_FILE"),
	}
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic(fmt.Errorf("required environment variable missing: %s", key))
	}
	return value
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	value := GetEnv(key, strconv.Itoa(fallback))
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Println("Failed to read env var into an int")
		panic(err)
	}
	return intValue
}

func GetEnvList(key string, delimiter string) []string {
	value := GetEnv(key, "")
	return strings.Split(value, delimiter)
}
