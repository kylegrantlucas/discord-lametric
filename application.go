package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kylegrantlucas/lametric-discord/services"

	"github.com/bwmarrin/discordgo"
)

var allowedChannelIDs = []string{}
var serverID, lametricIP, lametricAPIKey string
var iconID int

func init() {
	var err error
	serverID = os.Getenv("DISCORD_SERVER_ID")
	lametricIP = os.Getenv("LAMETRIC_IP")
	lametricAPIKey = os.Getenv("LAMETRIC_API_KEY")
	iconID, err = strconv.Atoi(os.Getenv("LAMETRIC_ICON_ID"))
	if err != nil {
		iconID = 0
	}
}
func main() {
	discord, err := discordgo.New(os.Getenv("DISCORD_EMAIL"), os.Getenv("DISCORD_PASSWORD"))
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	allowedChannels := strings.Split(os.Getenv("DISCORD_CHANNELS"), ",")

	channels, err := discord.GuildChannels(serverID)
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	for _, ch := range channels {
		for _, allowed := range allowedChannels {
			if ch.Name == allowed {
				allowedChannelIDs = append(allowedChannelIDs, ch.ID)
			}
		}
	}

	log.Printf("allowing on channels: %v", allowedChannelIDs)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.GuildID == serverID && channelIsAllowed(m.ChannelID) {
		client := lametric.Client{Host: lametricIP, APIKey: lametricAPIKey}
		err := client.Notify(iconID, m.Content)
		if err != nil {
			log.Print(err)
		}
	}
}

func channelIsAllowed(id string) bool {
	for _, allowed := range allowedChannelIDs {
		if id == allowed {
			return true
		}
	}

	return false
}
