FROM golang:1.24.5-bullseye AS builder

WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o build/ent-notifier . \
 && chmod +x /app/build/ent-notifier \
 && chown 1000:1000 /app/build/ent-notifier

FROM debian:bullseye-20250721-slim
WORKDIR /app/

RUN apt-get update && apt-get -y --no-install-recommends install tini ca-certificates \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

COPY layers_osrs/mapsquares/-1/2/0_* /app/mapsquares/
COPY --from=builder /app/build/ent-notifier /app/ent-notifier
USER 1000

ENV MAP_FILE_PATH=/app/mapsquares
ENTRYPOINT ["tini", "-v", "--", "/app/ent-notifier"]
STOPSIGNAL SIGINT
