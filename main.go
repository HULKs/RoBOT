package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"RoBOT/config"
	"RoBOT/errors"
	"RoBOT/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new Discord session using the provided bot token
	session, err := discordgo.New("Bot " + config.RoBotConfig.Token)
	errors.Check(err, "Error creating Discord session")

	// Add handler for MessageCreate events
	session.AddHandler(handlers.MessageCreate)

	// Add Handler for Ready events
	session.AddHandler(handlers.Ready)

	// Add handler for VoiceStateUpdate
	session.AddHandler(handlers.VoiceStateUpdate)

	// Open a websocket connection to Discord and begin listening
	err = session.Open()
	errors.Check(err, "Error opening connection")

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session
	err = session.Close()
	errors.Check(err, "Failed to close session properly")
}
