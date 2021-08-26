package handlers

import (
	"strings"

	"RoBOT/commands"
	"RoBOT/errors"

	"github.com/bwmarrin/discordgo"
)

// MessageCreate is called every time a new message is created in a channel that the bot has access to
func MessageCreate(s *discordgo.Session, ev *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if ev.Author.ID == s.State.User.ID {
		return
	}
	// Parse the message content and split into arguments
	input := ParseInput(ev.Content)
	// If there are no arguments, it was no command
	if input == nil {
		return
	}

	command := commands.GetCommand(strings.ToLower(input[0]))
	// Check if command is nil before doing sth with it
	if command.Run == nil {
		// TODO Send message that user entered nonexistent command
		return
	}
	// Check for help command
	// TODO this should be done differently, maybe with a custom struct and helpers to make embeds
	if len(input) > 1 && input[1] == "help" {
		command.Help(s, ev, input)
		return
	}
	// Call all the methods for the associated command
	command.Run(s, ev, input[1:])
}

// Ready is called when receiving the "ready" event
func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	// Set the playing status
	err := s.UpdateGameStatus(0, "with gophers...")
	errors.Check(err, "Failed setting custom status")
}
