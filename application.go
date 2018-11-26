package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kylegrantlucas/discord-lametric/models"
	"github.com/kylegrantlucas/discord-lametric/services"

	"github.com/bwmarrin/discordgo"
)

var allowedChannelsMap = map[string]string{}
var allowedServersMap = map[string]string{}
var config models.Config
var lametricClient lametric.Client

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	var file string
	if os.Getenv("CONFIG_FILE") != "" {
		file = os.Getenv("CONFIG_FILE")
	} else {
		file = "./config.default.json"
	}

	configJSON, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	lametricClient = lametric.Client{Host: config.LaMetricIP, APIKey: config.LaMetricAPIKey}
}

func main() {
	discord, err := discordgo.New(config.DiscordEmail, config.DiscordPassword, config.DiscordToken)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	allowedServersMap, err = buildAllowedServerList(discord, config)
	if err != nil {
		log.Fatalf("Error building server list: %v", err)
	}

	allowedChannelsMap, err = buildAllowedChannelList(discord, allowedServersMap, config)
	if err != nil {
		log.Fatalf("Error building channel list: %v", err)
	}

	log.Printf("allowing on servers: %v", allowedServersMap)
	log.Printf("allowing on channels: %v", allowedChannelsMap)

	// Register messageCreate callback
	discord.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func buildLaMetricNotification(message string, channelConfig *models.ChannelConfig) lametric.Notification {
	notif := lametric.Notification{
		IconType: channelConfig.IconType,
		Priority: channelConfig.Priority,
		Model: lametric.Model{
			Frames: []lametric.Frame{
				{
					Icon: channelConfig.Icon,
					Text: message,
				},
			},
		},
	}

	if channelConfig.Sound != nil {
		notif.Model.Sound = &lametric.Sound{
			Category: channelConfig.Sound.Category,
			ID:       channelConfig.Sound.ID,
			Repeat:   channelConfig.Sound.Repeat,
		}
	}

	return notif
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if serverIsAllowed(m.GuildID) && channelIsAllowed(m.ChannelID) {
		channelConfig := fetchChannelConfig(m.ChannelID)
		if channelConfig != nil {
			err := lametricClient.Notify(
				buildLaMetricNotification(m.Message.Content, channelConfig),
			)
			if err != nil {
				log.Print(err)
			}
		} else {
			log.Printf("channel is not configured: %v, %v", allowedChannelsMap[m.ChannelID], m.ChannelID)
		}
	}
}

func fetchChannelConfig(channelID string) *models.ChannelConfig {
	for _, server := range config.Servers {
		for _, channel := range server.Channels {
			if channel.Name == allowedChannelsMap[channelID] {
				return &channel
			}
		}
	}

	return nil
}

func channelIsAllowed(id string) bool {
	if allowedChannelsMap[id] != "" {
		return true
	}

	return false
}

func serverIsAllowed(id string) bool {
	if allowedServersMap[id] != "" {
		return true
	}

	return false
}

func buildAllowedChannelList(discord *discordgo.Session, serversMap map[string]string, config models.Config) (map[string]string, error) {
	ids := map[string]string{}
	for serverID := range serversMap {
		channels, err := discord.GuildChannels(serverID)
		if err != nil {
			log.Fatalf("Error building channel list: %v", err)
			return map[string]string{}, err
		}

		for _, ch := range channels {
			for _, server := range config.Servers {
				for _, allowed := range server.Channels {
					if ch.Name == allowed.Name {
						ids[ch.ID] = ch.Name
					}
				}
			}
		}
	}

	return ids, nil
}

func buildAllowedServerList(discord *discordgo.Session, config models.Config) (map[string]string, error) {
	ids := map[string]string{}
	guilds, err := discord.UserGuilds(50, "", "")
	if err != nil {
		return map[string]string{}, err
	}

	for _, guild := range guilds {
		for _, allowed := range config.Servers {
			if guild.Name == allowed.Name {
				ids[guild.ID] = guild.Name
			}
		}
	}

	return ids, nil
}
