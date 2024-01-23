package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HULKs/RoBOT/commands"
	"github.com/HULKs/RoBOT/config"
	"github.com/HULKs/RoBOT/handlers"
	"github.com/HULKs/RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

func main() {
	// Print license notice
	fmt.Printf(
		"RoBOT Copyright (C) 2020-2022 HULKs e.V.\n" +
			"This program comes with ABSOLUTELY NO WARRANTY.\n" +
			"This is free software, and you are welcome to redistribute\n" +
			"it under certain conditions.\n" +
			"See https://www.gnu.org/licenses/gpl-3.0.html for details.\n\n",
	)

	// Create a new Discord session using the provided bot token
	s, err := dg.New("Bot " + config.RoBotConfig.Token)
	util.ErrCheck(err, "Error creating Discord session")

	// Explicitly enable state cache
	s.StateEnabled = true
	// Enable Intents
	s.Identify.Intents = dg.IntentsAll

	// Add Handler for Ready events
	s.AddHandler(handlers.Ready)

	// Add handler for VoiceStateUpdate
	s.AddHandler(handlers.VoiceStateUpdate)

	s.AddHandler(
		func(s *dg.Session, i *dg.InteractionCreate) {
			if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		},
	)

	// Open a websocket connection to Discord and begin listening
	err = s.Open()
	util.ErrCheck(err, "Error opening connection")

	log.Println("Adding commands...")
	registeredCommands := make([]*dg.ApplicationCommand, len(commands.Commands))
	for i, c := range commands.Commands {
		// TODO This won't work on first run? config isn't populated yet?
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, config.ServerConfig.GuildID, c)
		util.ErrCheck(err, fmt.Sprintf("Cannot create '%s' command", c.Name))
		registeredCommands[i] = cmd
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, os.Interrupt, os.Kill)
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	// Wait here until CTRL-C or other term signal is received
	<-sc

	log.Println("Removing commands...")
	allCommands, err := s.ApplicationCommands(s.State.User.ID, config.ServerConfig.GuildID)
	util.ErrCheck(err, "Could not fetch global commands")

	for _, cmd := range allCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, config.ServerConfig.GuildID, cmd.ID)
		util.ErrCheck(err, fmt.Sprintf("Cannot delete '%s' command", cmd.Name))
	}

	log.Println("Gracefully shutdowning")

	// Cleanly close down the Discord session
	err = s.Close()
	util.ErrCheck(err, "Failed to close session properly")
}
