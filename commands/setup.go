package commands

import (
	"fmt"
	"log"
	"strings"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/helptexts"
	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

const logCategory string = "Setup"

func setupRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {

	// ErrCheck for arguments
	if len(args) == 0 {
		msg, err := s.ChannelMessageSend(ev.ChannelID, "You didn't give any arguments!")
		util.CheckMsgSend(err, msg.GuildID, msg.ChannelID)
		return
	}

	// Get guild ID and save to ServerConfig
	g, err := s.Guild(ev.GuildID)
	util.ErrCheck(err, "Error getting guild for ID "+ev.GuildID)
	config.ServerConfig.GuildID = g.ID

	getEveryoneRoleID(s, g)

	switch strings.ToLower(args[0]) {
	case "bootstrap":
		// Delete everything on server
		deleteChannelsAndRoles(s, g)

		// Reset list of protected channels
		config.ServerConfig.ProtectedChannels = make(map[string]interface{})

		// Rename server
		_, err := s.GuildEdit(g.ID, discordgo.GuildParams{Name: config.ServerConfig.EventName})
		util.ErrCheck(err, "Failed renaming Server")

		// Create Orga-Team role
		util.CreateRole(
			s, g, "Orga-Team", "0x9A58B4",
			config.ServerConfig.PermissionTemplates.OrgaTeam,
			true, true, &config.ServerConfig.OrgaTeamRoleID,
		)
		// Create RoBOT-Admin role
		util.CreateRole(
			s, g, "RoBOT-Admin", "0xFF0000",
			config.ServerConfig.PermissionTemplates.RoBOTAdmin,
			true, true, &config.ServerConfig.RoBOTAdminRoleID,
		)
		// Create role for each team
		log.Println("[Setup] Creating roles for teams...")
		for _, t := range config.TeamList {
			util.CreateRole(
				s, g, t.Name, t.TeamColor,
				config.ServerConfig.PermissionTemplates.TeamRole,
				true, true, &t.RoleID,
			)
		}
		// Create Participant role
		util.CreateRole(
			s, g, "Participant", "0x000000",
			config.ServerConfig.PermissionTemplates.Participant,
			false, false, &config.ServerConfig.ParticipantRoleID,
		)

		createBasicChannels(s, g, ev)
		sendRoleAssignmentMessage(s, g)

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
	util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}

// getEveryoneRoleID saves the ID for the @everyone role to the ServerConfig
func getEveryoneRoleID(s *discordgo.Session, g *discordgo.Guild) {
	// TODO Add logging
	// Get roles for guild
	roles, err := s.GuildRoles(g.ID)
	util.ErrCheck(err, "Error getting roles for "+g.Name)
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
	util.ErrCheck(err, "Failed to get channels for ID "+g.ID)
	// Delete all channels
	for _, channel := range channels {
		_, err = s.ChannelDelete(channel.ID)
		util.ErrCheck(err, "Failed deleting channel "+channel.ID)
	}
	// Get all Roles
	roles, err := s.GuildRoles(g.ID)
	for _, role := range roles {
		// If role ID is @everyone, we can't delete this, also we can't delete
		// the bot's own role, which should have Administrator permissions
		if role.ID == config.ServerConfig.EveryoneRoleID ||
			((role.Name == s.State.User.Username) && role.Permissions == 8) {
			continue
		}
		// Delete Role
		err = s.GuildRoleDelete(g.ID, role.ID)
		util.ErrCheck(err, "Failed deleting role "+role.ID)
	}
	// Set server-wide permissions for @everyone
	_, err = s.GuildRoleEdit(
		g.ID, config.ServerConfig.EveryoneRoleID, "",
		0, false, config.ServerConfig.PermissionTemplates.Everyone, true,
	)
	util.ErrCheck(err, "Failed setting permissions for @everyone role")
}

func createBasicChannels(s *discordgo.Session, g *discordgo.Guild, ev *discordgo.MessageCreate) {
	// welcome
	chWelcome := util.CreateChannel(
		s, g, "welcome", "", "", discordgo.ChannelTypeGuildText, []*discordgo.PermissionOverwrite{
			// read-only for @everyone
			{
				ID:   config.ServerConfig.EveryoneRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0,
				Allow: discordgo.PermissionViewChannel |
					discordgo.PermissionReadMessageHistory,
			},
		}, logCategory, ev.Author.Username,
	)
	config.ServerConfig.ProtectedChannels[chWelcome.ID] = nil
	// Save ID to ServerConfig
	config.ServerConfig.WelcomeChannelID = chWelcome.ID
	// TODO set welcome as system channel (Doesn't work, field in GuildParams missing)

	// role-assignment
	chRoleAssignment := util.CreateChannel(
		s, g, "role-assignment", "", "", discordgo.ChannelTypeGuildText, []*discordgo.PermissionOverwrite{
			// write for @everyone
			{
				ID:   config.ServerConfig.EveryoneRoleID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0,
				Allow: discordgo.PermissionViewChannel |
					discordgo.PermissionReadMessageHistory |
					discordgo.PermissionSendMessages,
			},
		}, logCategory, ev.Author.Username,
	)
	config.ServerConfig.ProtectedChannels[chRoleAssignment.ID] = nil
	// Save ID to ServerConfig
	config.ServerConfig.RoleAssignmentChannelID = chRoleAssignment.ID

	// botcontrol
	chBotcontrol := util.CreateChannel(
		s, g, "botcontrol", "", "", discordgo.ChannelTypeGuildText, util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.RoBOTAdminRoleID,
		), logCategory, ev.Author.Username,
	)
	config.ServerConfig.ProtectedChannels[chBotcontrol.ID] = nil

	// INFORMATION: announcements, links
	catInformation := util.CreateCategory(
		s, g, "Information", "", util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.ParticipantRoleID,
		), logCategory, ev.Author.Username,
	)
	_ = util.CreateChannel(
		s, g, "announcements", "", catInformation.ID, discordgo.ChannelTypeGuildText, nil, logCategory, ev.Author.Username,
	)
	_ = util.CreateChannel(
		s, g, "links", "", catInformation.ID, discordgo.ChannelTypeGuildText, nil, logCategory, ev.Author.Username,
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
			s, g, t.Name, t.Name+" - Teamzone", []*discordgo.PermissionOverwrite{
				// Hide for @everyone
				{
					ID:   config.ServerConfig.EveryoneRoleID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel |
						discordgo.PermissionVoiceConnect,
					Allow: 0,
				},
				// All permissions for team
				{
					ID:   t.RoleID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: 0,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionVoiceConnect |
						discordgo.PermissionManageChannels,
				},
			}, logCategory, ev.Author.Username,
		)
		// Save ID to ProtectedChannels
		config.ServerConfig.ProtectedChannels[catTeamzone.ID] = nil
		// Save Teamzone ID to TeamConf entry in config.TeamList
		t.TeamzoneID = catTeamzone.ID
		// Create text channel
		util.CreateChannel(
			s, g, t.Name, t.Name+" - Teamzone", catTeamzone.ID, discordgo.ChannelTypeGuildText, nil,
			logCategory, ev.Author.Username,
		)
		// Create voice channels
		for i := 1; i < 4; i++ {
			util.CreateChannel(
				s, g, fmt.Sprintf("Teamzone %02d", i), t.Name+" - Teamzone", catTeamzone.ID, discordgo.ChannelTypeGuildVoice,
				nil, logCategory, ev.Author.Username,
			)
		}
	}

	// Create category for magic voice channels
	catMagic := util.CreateCategory(
		s, g, "Create Meetings", "", util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.ParticipantRoleID,
		), logCategory, ev.Author.Username,
	)
	config.ServerConfig.ProtectedChannels[catMagic.ID] = nil
	// Create channel to create moar channels
	chMagicVoice := util.CreateChannel(
		s, g, "Click to create room", "", catMagic.ID, discordgo.ChannelTypeGuildVoice, nil, logCategory,
		ev.Author.Username,
	)
	// Save ID to ServerConfig
	config.ServerConfig.VoiceChannelCreateID = chMagicVoice.ID

	// Create Archive category
	catArchive := util.CreateCategory(
		s, g, "Archive", "Archived channels", util.PermOverwriteHideForAShowForB(
			config.ServerConfig.EveryoneRoleID,
			config.ServerConfig.ParticipantRoleID,
		), logCategory, ev.Author.Username,
	)
	// Save to config
	config.ServerConfig.ArchiveCategoryID = catArchive.ID
	config.ServerConfig.ProtectedChannels[catArchive.ID] = nil
}

func sendRoleAssignmentMessage(s *discordgo.Session, g *discordgo.Guild) {
	// Generate Embed with all teams
	var desc strings.Builder
	desc.WriteString("```text")
	for i, team := range config.TeamList {
		desc.WriteString(fmt.Sprintf("\n%2d: %s", i, team.Name))
	}
	desc.WriteString("\n```")

	embed := discordgo.MessageEmbed{
		Title: "Team Assignment",
		Description: fmt.Sprintf(
			"Greetings! You have reached the role assignment channel! "+
				"To get access to the digital venue, assign yourself a role using the `%steam $TEAMNUMBER` command. "+
				"The Bot will then assign you to your team and grant you access to the rest of the server.",
			config.RoBotConfig.Prefix,
		),
		Color:  colors.GREEN,
		Footer: util.HelpEmbedFooter(),
		Image:  nil, // TODO There should be an image here
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Team List",
				Value: desc.String(),
			},
		},
	}

	_, err := s.ChannelMessageSendEmbed(config.ServerConfig.RoleAssignmentChannelID, &embed)
	util.ErrCheck(err, "[Setup] Failed sending role assignment message")
}
