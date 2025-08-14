package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
)

const eventsLocation = "[Log Hunters](https://discord.com/channels/922245627092541450/1166753438177173534)"

func NotifyEvents(events []*Event, webhookUrls []string) error {
	var embeds []DiscordEmbed
	var files []DiscordFile

	sort.Slice(events, func(a, b int) bool {
		return events[a].DiscoveredTime < events[b].DiscoveredTime
	})

	for num, event := range events {
		mapLink := fmt.Sprintf("[Map](https://mejrs.github.io/osrs?m=-1&z=4&p=0&x=%d&y=%d)", event.X, event.Y)

		mapImage, err := CreateThumbnail(event.X, event.Y, MapImageWidth, MapImageHeight)

		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, mapImage)
		if err != nil {
			return err
		}
		imageName := fmt.Sprintf("map%d.png", num)

		embeds = append(embeds, DiscordEmbed{
			Title: event.EventType,
			Fields: &[]DiscordEmbedField{
				{
					Name:   "Discovered",
					Value:  fmt.Sprintf("<t:%d:R>", event.DiscoveredTime),
					Inline: true,
				},
				{
					Name:   "Links",
					Value:  fmt.Sprintf("%s\n%s", eventsLocation, mapLink),
					Inline: true,
				},
			},
			Image: &DiscordEmbedImage{
				Url: fmt.Sprintf("attachment://%s", imageName),
			},
		})

		imageBuffer := buffer.Bytes()
		files = append(files, DiscordFile{
			Name: imageName,
			Data: &imageBuffer,
		})
	}

	message := &DiscordMessage{
		Embeds: &embeds,
		Files:  &files,
	}

	for _, url := range webhookUrls {
		if len(url) == 0 {
			continue
		}

		webhookUrl := url
		var roleId *string

		if strings.Contains(webhookUrl, "=") {
			parts := strings.Split(webhookUrl, "=")
			webhookUrl = parts[0]
			roleId = &parts[1]
		}

		if roleId != nil {
			message.Content = fmt.Sprintf("<@&%s>", *roleId)
		}

		err := postMessage(webhookUrl, message)
		if err != nil {
			fmt.Println("Failed to post to webhook url", url, err)
		}
	}
	return nil
}

func postMessage(webhookUrl string, message *DiscordMessage) error {
	payload, contentType, err := encodeMessage(message)
	if err != nil {
		return err
	}

	_, err = http.Post(webhookUrl, contentType, payload)
	if err != nil {
		return err
	}

	return nil
}

func encodeMessage(message *DiscordMessage) (*bytes.Buffer, string, error) {
	payload := new(bytes.Buffer)

	if message.Files != nil {
		writer := multipart.NewWriter(payload)

		partWriter, err := writer.CreateFormField("payload_json")
		if err != nil {
			return nil, "", err
		}

		err = json.NewEncoder(partWriter).Encode(message)
		if err != nil {
			return nil, "", err
		}

		for index, file := range *message.Files {
			partWriter, err := writer.CreateFormFile(fmt.Sprintf("files[%v]", index), file.Name)
			if err != nil {
				return nil, "", err
			}

			_, err = partWriter.Write(*file.Data)
			if err != nil {
				return nil, "", err
			}
		}

		err = writer.Close()
		if err != nil {
			return nil, "", err
		}

		return payload, "multipart/form-data; boundary=" + writer.Boundary(), nil
	} else {
		err := json.NewEncoder(payload).Encode(message)
		return payload, "application/json", err
	}
}
