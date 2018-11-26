package models

type Config struct {
	Servers         []ServerConfig `json:"servers,omitempty"`
	DiscordEmail    string         `json:"discord_email,omitempty"`
	DiscordPassword string         `json:"discord_password,omitempty"`
	DiscordToken    string         `json:"discord_token,omitempty"`
	LaMetricIP      string         `json:"lametric_ip,omitempty"`
	LaMetricAPIKey  string         `json:"lametric_api_key,omitempty"`
}

type ServerConfig struct {
	Name     string          `json:"server_name,omitempty"`
	Channels []ChannelConfig `json:"channels,omitempty"`
}

type ChannelConfig struct {
	Name     string       `json:"channel_name,omitempty"`
	Icon     *string      `json:"icon,omitempty"`
	Priority *string      `json:"priority,omitempty"`
	IconType *string      `json:"icon_type,omitempty"`
	Lifetime *string      `json:"lifetime,omitempty"`
	Sound    *SoundConfig `json:"sound,omitempty"`
}

type SoundConfig struct {
	Category string `json:"category,omitempty"`
	ID       string `json:"id,omitempty"`
	Repeat   int    `json:"repeat,omitempty"`
}
