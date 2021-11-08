package handlers

import (
	"log"
	"strings"

	"RoBOT/commands"
	"RoBOT/config"
	"RoBOT/util"

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

// VoiceStateUpdate is called every time a user joins a voice channel,
// but we only care about the magic meeting creation channel
func VoiceStateUpdate(s *discordgo.Session, ev *discordgo.VoiceStateUpdate) {
	// Ignore events caused by bot itself
	if ev.UserID == s.State.User.ID {
		return
	}
	// Check if Magic channel
	if ev.ChannelID != config.ServerConfig.VoiceChannelCreateID {
		return
	}
	// Get guild
	g, err := s.Guild(ev.GuildID)
	util.ErrCheck(err, "[VoiceStateUpdate] Failed getting guild for ID "+ev.GuildID)
	// Get user
	user, err := s.User(ev.UserID)
	util.ErrCheck(err, "[VoiceStateUpdate] Failed getting user for UserID "+ev.UserID)
	// Get Archive category
	catArchive, err := s.Channel(config.ServerConfig.ArchiveCategoryID)
	util.ErrCheck(err, "[VoiceStateUpdate] Failed getting Archive category channel")

	// Create new category and new text and voice channel
	log.Printf("[VoiceStateUpdate] Creating new meeting for user %s", user.Username)
	catNew, err := s.GuildChannelCreateComplex(
		g.ID, discordgo.GuildChannelCreateData{
			Name:     "New Meeting",
			Type:     discordgo.ChannelTypeGuildCategory,
			Position: catArchive.Position - 1,
			PermissionOverwrites: append(
				// Hide for @everyone, Default permissions for @Participant
				util.PermOverwriteHideForAShowForB(
					config.ServerConfig.EveryoneRoleID,
					config.ServerConfig.ParticipantRoleID,
				),
				// Add management permissions for meeting organizer
				&discordgo.PermissionOverwrite{
					ID:    ev.UserID,
					Type:  discordgo.PermissionOverwriteTypeRole,
					Deny:  0,
					Allow: discordgo.PermissionManageChannels,
				},
			),
		},
	)
	// Create Text channel and Voice channel
	util.CreateChannel(
		s, g, "text", "", catNew.ID, discordgo.ChannelTypeGuildText, nil,
	)
	chVoice := util.CreateChannel(
		s, g, "voice", "", catNew.ID, discordgo.ChannelTypeGuildVoice, nil,
	)

	// Move user to new channel
	err = s.GuildMemberMove(g.ID, user.ID, &chVoice.ID)
	util.ErrCheck(err, "[VoiceStatUpdate] Failed moving user "+user.Username+" to new voice channel")
	log.Printf("[VoiceStateUpdate] Moved user %s to channel %s after creating new meeting", user.Username, chVoice.Name)
}

// Ready is called when receiving the "ready" event
func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	// Set the playing status
	err := s.UpdateGameStatus(0, "with gophers...")
	util.ErrCheck(err, "Failed setting custom status")
}
