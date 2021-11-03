package commands

import (
	"log"
	"strings"

	"RoBOT/config"
	"RoBOT/errors"
	"RoBOT/helptexts"
	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

func setupRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {

	// Check for arguments
	if len(args) == 0 {
		msg, err := s.ChannelMessageSend(ev.ChannelID, "You didn't give any arguments!")
		errors.CheckMsgSend(err, msg.GuildID, msg.ChannelID)
		return
	}

	// Get guild ID and save to ServerConfig
	g, err := s.Guild(ev.GuildID)
	errors.Check(err, "Error getting guild for ID "+ev.GuildID)
	config.ServerConfig.GuildID = g.ID

	// Save own user ID to ServerConfig
	selfUser, err := s.User("@me")
	errors.Check(err, "Failed getting User struct for @me")
	config.ServerConfig.BotUserID = selfUser.ID

	getEveryoneRoleID(s, g)

	switch strings.ToLower(args[0]) {
	case "bootstrap":
		// Delete everything on server
		deleteChannelsAndRoles(s, g)

		// Create Participant role
		util.CreateRole(
			s, g, "Participant", "0x000000",
			config.ServerConfig.PermissionTemplates.Participant,
			false, false, &config.ServerConfig.ParticipantRoleID,
		)
		// Create Orga-Team role
		util.CreateRole(
			s, g, "Orga-Team", "0x9A58B4",
			config.ServerConfig.PermissionTemplates.OrgaTeam,
			false, false, &config.ServerConfig.OrgaTeamRoleID,
		)
		// Create RoBOT-Admin role
		util.CreateRole(
			s, g, "RoBOT-Admin", "0xFF0000",
			config.ServerConfig.PermissionTemplates.RoBOTAdmin,
			false, true, &config.ServerConfig.BotAdminRoleID,
		)
		// Create role for each team
		log.Println("Creating roles for teams...")
		for _, t := range config.TeamList {
			util.CreateRole(
				s, g, t.Name, t.TeamColor,
				config.ServerConfig.PermissionTemplates.TeamRole,
				true, true, &t.RoleID,
			)
		}

		createBasicChannels(s, g)

		// TODO Create Teamzones
		// TODO Create Magic Voice

		config.SaveServerConfig()
		config.SaveTeamConfig()
	case "add-team":
		// TODO
	case "add-channel":
		// TODO
	}
}

func setupHelp(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, helptexts.DB["ping"])
	errors.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}

// getEveryoneRoleID saves the ID for the @everyone role to the ServerConfig
func getEveryoneRoleID(s *discordgo.Session, g *discordgo.Guild) {
	// TODO Add logging
	// Get roles for guild
	roles, err := s.GuildRoles(g.ID)
	errors.Check(err, "Error getting roles for "+g.Name)
	// Get ID for @everyone
	for _, role := range roles {
		if role.Name == "@everyone" {
			config.ServerConfig.EveryoneRoleID = role.ID
			return
		}
	}
	panic("Could not find a role named \"@everyone\"")
}

func deleteChannelsAndRoles(s *discordgo.Session, g *discordgo.Guild) {
	// TODO Add logging
	// Get all channels in server
	channels, err := s.GuildChannels(g.ID)
	errors.Check(err, "Failed to get channels for ID "+g.ID)
	// Delete all channels
	for _, channel := range channels {
		_, err := s.ChannelDelete(channel.ID)
		errors.Check(err, "Failed deleting channel "+channel.ID)
	}
	// Get all Roles
	roles, err := s.GuildRoles(g.ID)
	for _, role := range roles {
		if role.ID == config.ServerConfig.EveryoneRoleID ||
			(strings.Contains(role.Name, "RoBOT") && role.Permissions == 8) {
			continue
		}
		// Delete Role
		err := s.GuildRoleDelete(g.ID, role.ID)
		errors.Check(err, "Failed deleting role "+role.ID)
	}
	// Set server-wide permissions for @everyone
	_, err = s.GuildRoleEdit(g.ID, config.ServerConfig.EveryoneRoleID, "", 0, false, 0, true)
	errors.Check(err, "Failed setting permissions for @everyone role")
}

func createBasicChannels(s *discordgo.Session, g *discordgo.Guild) {
	// welcome
	_ = util.CreateChannel(
		s, g, "welcome", "", "", discordgo.ChannelTypeGuildText,
		[]*discordgo.PermissionOverwrite{
			// read-only for @everyone
			{
				ID:   config.ServerConfig.EveryoneRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0,
				Allow: discordgo.PermissionViewChannel |
					discordgo.PermissionReadMessageHistory,
			},
		},
	)
	// TODO set welcome as system channel (Doesnt work, field in GuildParams missing)

	// role-assignment
	_ = util.CreateChannel(
		s, g, "role-assignment", "", "", discordgo.ChannelTypeGuildText,
		[]*discordgo.PermissionOverwrite{
			// write for @everyone
			{
				ID:   config.ServerConfig.EveryoneRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0,
				Allow: discordgo.PermissionViewChannel |
					discordgo.PermissionReadMessageHistory |
					discordgo.PermissionSendMessages,
			},
		},
	)

	// botcontrol
	_, err = s.GuildChannelCreate(g.ID, "botcontrol", 0)
	errors.Check(err, "Failed to create botcontrol channel")
	// TODO Only Bot Admins

	// INFORMATION: announcements, links

	// META: help, feedback

	// PHOTOS: photo-wall, participant-selfie

	// GENERAL: town-hall, Voice: Town-Hall, Lounge 01-02, AFK

}
