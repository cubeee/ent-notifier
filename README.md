# Ent notifier

Sends Discord webhooks when ent events are scouted, along with a precise map image.

![Example webhook](/misc/example.png)

## Usage
Run in a Docker container with at least the required environment variables set.

## Environment variables
- `SLEEP_TIME_SECONDS`: Time slept between event api polls. Defaults to 5
- `HISTORY_LOOKUP_SECONDS`: Seconds to look for events in the past on first poll. Defaults to 15
- `EVENT_AREA_RADIUS`: Tiles in all directions of a discovered event. Other events/scouts in that area are considered duplicate. Defaults to 8
- `EVENTS_API_URL`: URL to get events from (**required**)
- `EVENTS_API_TIMEOUT`: Seconds to wait for a response from the events api url. Defaults to 5 seconds
- `EVENTS_API_USER_AGENT`: User agent to send when fetching from events api (**required**)
- `EVENTS_WORLD`: World to receive events from. Defaults to 444
- `EVENTS_ALLOWED`: Allowed event types to notify. Comma-delimiter string list
- `DISCORD_WEBHOOK_URLS`: Discord webhooks to post events to. Defaults to nothing. Use the optional format `<url>=<role_id>` to ping a role
- `MAP_FILE_PATH`: Directory where map chunk image files are located. Defaults to `layers_osrs/mapsquares/-1/2`
- `MAP_TILE_PIXELS`: Pixels per tile in the map chunk files. Defaults to 4
- `MAP_IMAGE_WIDTH`: Embedded map image width in pixels. Defaults to 400
- `MAP_IMAGE_HEIGHT`: Embedded map image height in pixels. Defaults to 250
- `LOCATION_COOLDOWN`: Seconds to wait between events at a specific location before not considering scouts duplicates
- `PAST_EVENT_MAX_AGE`: Number of seconds to store past events for, used for checking duplicate scouts. Defaults to 30 minutes.
