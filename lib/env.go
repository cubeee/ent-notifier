package lib

import "os"

var (
	EventAreaRadius      = GetEnvInt("EVENT_AREA_RADIUS", 8)
	MapImageWidth        = GetEnvInt("MAP_IMAGE_WIDTH", 400)
	MapImageHeight       = GetEnvInt("MAP_IMAGE_HEIGHT", 250)
	EmbedFooter          = GetEnv("EMBED_FOOTER", "")
	ApiUrl               = os.Getenv("EVENTS_API_URL")
	ApiTimeout           = GetEnvInt("EVENTS_API_TIMEOUT", 5)
	EventsWorld          = GetEnvInt("EVENTS_WORLD", 444)
	EventsAllowed        = GetEnvList("EVENTS_ALLOWED", ",")
	UserAgent            = os.Getenv("USER_AGENT")
	LocationCooldown     = GetEnvInt("LOCATION_COOLDOWN", 300)
	MapFilePath          = GetEnv("MAP_FILE_PATH", "layers_osrs/mapsquares/-1/2")
	MapTilePixels        = GetEnvInt("MAP_TILE_PIXELS", 4)
	SleepTime            = GetEnvInt("SLEEP_TIME_SECONDS", 5)
	HistoryLookupSeconds = GetEnvInt("HISTORY_LOOKUP_SECONDS", 0)
	WebhookUrls          = GetEnvList("DISCORD_WEBHOOK_URLS", ",")
	PastEventMaxAge      = GetEnvInt("PAST_EVENT_MAX_AGE", 30*60)
)
