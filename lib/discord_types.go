package lib

type DiscordMessage struct {
	Files   *[]DiscordFile  `json:"-"`
	Content string          `json:"content,omitempty"`
	Embeds  *[]DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordFile struct {
	Name string
	Data *[]byte
}

type DiscordEmbed struct {
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Image       *DiscordEmbedImage   `json:"image,omitempty"`
	Fields      *[]DiscordEmbedField `json:"fields,omitempty"`
	Footer      *DiscordFooter       `json:"footer,omitempty"`
}

type DiscordEmbedImage struct {
	Url string `json:"url"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type DiscordFooter struct {
	Text    *string `json:"text,omitempty"`
	IconUrl *string `json:"icon_url,omitempty"`
}
