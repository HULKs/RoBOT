package commands

import (
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

	// Get Guild ID
	g, err := s.Guild(ev.GuildID)
	errors.Check(err, "Error getting guild for ID "+ev.GuildID)

	// Save to ServerConfig
	config.ServerConfig.GuildID = g.ID

	getEveryoneRoleID(s, g)

	switch strings.ToLower(args[0]) {
	case "bootstrap":
		deleteChannelsAndRoles(s, g)

		util.CreateRole(
			s, g, "Participant", "0x000000", 242769972800, false, false, &config.ServerConfig.ParticipantRoleID,
		)
		util.CreateRole(s, g, "Orga-Team", "0x9A58B4", 0, false, false, &config.ServerConfig.OrgaTeamRoleID)
		util.CreateRole(s, g, "RoBOT-Admin", "0xFF0000", 0, false, true, &config.ServerConfig.BotAdminRoleID)
		createTeamRoles(s, g)

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
		if role.Name == "@everyone" || (strings.Contains(role.Name, "RoBOT") && role.Permissions == 8) {
			continue
		}
		// Delete Role
		err := s.GuildRoleDelete(g.ID, role.ID)
		errors.Check(err, "Failed deleting role "+role.ID)
	}
}

// createTeamRoles creates the roles for the participating teams
func createTeamRoles(s *discordgo.Session, g *discordgo.Guild) {
	// TODO Add logging
	for _, t := range config.TeamList {
		util.CreateRole(s, g, t.Name, t.TeamColor, 0, true, true, &t.RoleID)
	}
}

func createBasicChannels(s *discordgo.Session, g *discordgo.Guild) {
	var err error
	// welcome
	_, err = s.GuildChannelCreate(g.ID, "welcome", 0)
	errors.Check(err, "Failed to create welcome channel")
	// TODO set as system channel

	// botcontrol
	_, err = s.GuildChannelCreate(g.ID, "botcontrol", 0)
	errors.Check(err, "Failed to create botcontrol channel")
	// TODO Only Bot Admins

	// INFORMATION: announcements, links

	// META: help, feedback

	// PHOTOS: photo-wall, participant-selfie

	// GENERAL: town-hall, Voice: Town-Hall, Lounge 01-02, AFK

}
