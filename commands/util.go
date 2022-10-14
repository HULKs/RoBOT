package commands

import (
	"fmt"
	"log"

	"RoBOT/config"
	"RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

// saveEveryoneRoleID saves the ID for the @everyone role to the ServerConfig
func saveEveryoneRoleID(s *dg.Session, guildID string) {
	// TODO Add logging
	// Get roles for guild
	roles, err := s.GuildRoles(guildID)
	util.ErrCheck(err, "Error getting roles for guild ID: "+guildID)
	// Get ID for @everyone
	for _, role := range roles {
		if role.Name == "@everyone" {
			config.ServerConfig.EveryoneRoleID = role.ID
			return
		}
	}
	panic("Could not find a role named \"@everyone\"")
}

// clearServerAndReset reads all channels and roles for a guildID and deletes them, then resets @everyone
func clearServerAndReset(s *dg.Session, guildID string) {
	// TODO Add logging
	// Get all channels in server
	channels, err := s.GuildChannels(guildID)
	util.ErrCheck(err, "Failed to get channels for ID "+guildID)
	// Delete all channels
	for _, channel := range channels {
		_, err = s.ChannelDelete(channel.ID)
		util.ErrCheck(err, "Failed deleting channel "+channel.ID)
	}

	// Get all Roles and then delete them
	roles, err := s.GuildRoles(guildID)
	for _, role := range roles {
		// If role ID is @everyone, we can't delete this, also we can't delete
		// the bot's own role, which should have Administrator permissions
		if role.ID == config.ServerConfig.EveryoneRoleID ||
			((role.Name == s.State.User.Username) && role.Permissions == dg.PermissionAdministrator) {
			continue
		}
		// Delete Role
		err = s.GuildRoleDelete(guildID, role.ID)
		util.ErrCheck(err, "Failed deleting role "+role.ID)
	}

	// Set server-wide permissions for @everyone
	_, err = s.GuildRoleEdit(
		guildID, config.ServerConfig.EveryoneRoleID, &dg.RoleParams{
			Name:        "",
			Color:       util.PointyInt(0),
			Hoist:       util.PointyBool(false),
			Permissions: &config.ServerConfig.PermissionTemplates.Everyone,
			Mentionable: util.PointyBool(true),
		},
	)
	util.ErrCheck(err, "Failed setting permissions for @everyone role")
}

// createBasicChannels creates the basic channels (welcome, role-assignment, ...) with the desired config
func createBasicChannels(s *dg.Session, guildID, memberNick string) {
	// welcome
	chWelcome := util.CreateChannel(
		s, guildID, "welcome", "", "", dg.ChannelTypeGuildText, []*dg.PermissionOverwrite{
			// read-only for @everyone
			{
				ID:   config.ServerConfig.EveryoneRoleID,
				Type: dg.PermissionOverwriteTypeRole,
				Deny: 0,
				Allow: dg.PermissionViewChannel |
					dg.PermissionReadMessageHistory,
			},
		}, logCategory, memberNick,
	)
	config.ServerConfig.ProtectedChannels[chWelcome.ID] = nil
	// Save ID to ServerConfig
	config.ServerConfig.WelcomeChannelID = chWelcome.ID
	// TODO set welcome as system channel (Doesn't work, field in GuildParams missing)
	// s.GuildEdit(guildID, &dg.GuildParams{})

	// role-assignment
	chRoleAssignment := util.CreateChannel(
		s, guildID, "role-assignment", "", "", dg.ChannelTypeGuildText, []*dg.PermissionOverwrite{
			// write for @everyone
			{
				ID:   config.ServerConfig.EveryoneRoleID,
				Type: dg.PermissionOverwriteTypeRole,
				Deny: 0,
				Allow: dg.PermissionViewChannel |
					dg.PermissionReadMessageHistory |
					dg.PermissionSendMessages,
			},
		}, logCategory, memberNick,
	)
	config.ServerConfig.ProtectedChannels[chRoleAssignment.ID] = nil
	// Save ID to ServerConfig
	config.ServerConfig.RoleAssignmentChannelID = chRoleAssignment.ID

	// botcontrol
	chBotcontrol := util.CreateChannel(
		s, guildID, "botcontrol", "", "", dg.ChannelTypeGuildText, util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.RoBOTAdminRoleID,
		), logCategory, memberNick,
	)
	config.ServerConfig.ProtectedChannels[chBotcontrol.ID] = nil

	// INFORMATION: announcements, links
	catInformation := util.CreateCategory(
		s, guildID, "Information", "", util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.ParticipantRoleID,
		), logCategory, memberNick,
	)
	_ = util.CreateChannel(
		s, guildID, "announcements", "", catInformation.ID, dg.ChannelTypeGuildText, nil, logCategory, memberNick,
	)
	_ = util.CreateChannel(
		s, guildID, "links", "", catInformation.ID, dg.ChannelTypeGuildText, nil, logCategory, memberNick,
	)
	config.ServerConfig.ProtectedChannels[catInformation.ID] = nil

	// TODO more basic channels
	// META: help, feedback

	// PHOTOS: photo-wall, participant-selfie

	// GENERAL: town-hall, Voice: Town-Hall, Lounge 01-02, AFK

	// Create teamzone for each team
	log.Println("[Setup] Creating teamzones...")
	for _, t := range config.TeamList {
		// Create teamzone category
		catTeamzone := util.CreateCategory(
			s, guildID, t.Name, t.Name+" - Teamzone", []*dg.PermissionOverwrite{
				// Hide for @everyone
				{
					ID:   config.ServerConfig.EveryoneRoleID,
					Type: dg.PermissionOverwriteTypeRole,
					Deny: dg.PermissionViewChannel |
						dg.PermissionVoiceConnect,
					Allow: 0,
				},
				// All permissions for team
				{
					ID:   t.RoleID,
					Type: dg.PermissionOverwriteTypeRole,
					Deny: 0,
					Allow: dg.PermissionViewChannel |
						dg.PermissionVoiceConnect |
						dg.PermissionManageChannels,
				},
			}, logCategory, memberNick,
		)
		// Save ID to ProtectedChannels
		config.ServerConfig.ProtectedChannels[catTeamzone.ID] = nil
		// Save Teamzone ID to TeamConf entry in config.TeamList
		t.TeamzoneID = catTeamzone.ID
		// Create text channel
		util.CreateChannel(
			s, guildID, t.Name, t.Name+" - Teamzone", catTeamzone.ID, dg.ChannelTypeGuildText, nil,
			logCategory, memberNick,
		)
		// Create voice channels
		for i := 1; i < 4; i++ {
			util.CreateChannel(
				s, guildID, fmt.Sprintf("Teamzone %02d", i), t.Name+" - Teamzone", catTeamzone.ID,
				dg.ChannelTypeGuildVoice,
				nil, logCategory, memberNick,
			)
		}
	}

	// Create category for magic voice channels
	catMagic := util.CreateCategory(
		s, guildID, "Create Meetings", "", util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.ParticipantRoleID,
		), logCategory, memberNick,
	)
	config.ServerConfig.ProtectedChannels[catMagic.ID] = nil
	// Create channel to create moar channels
	chMagicVoice := util.CreateChannel(
		s, guildID, "Click to create room", "", catMagic.ID, dg.ChannelTypeGuildVoice, nil, logCategory,
		memberNick,
	)
	// Save ID to ServerConfig
	config.ServerConfig.VoiceChannelCreateID = chMagicVoice.ID

	// Create Archive category
	catArchive := util.CreateCategory(
		s, guildID, "Archive", "Archived channels", util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.ParticipantRoleID,
		), logCategory, memberNick,
	)
	// Save to config
	config.ServerConfig.ArchiveCategoryID = catArchive.ID
	config.ServerConfig.ProtectedChannels[catArchive.ID] = nil
}
