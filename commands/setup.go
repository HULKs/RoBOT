package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

func setupBootstrap(s *dg.Session, guildID string, i *dg.InteractionCreate) {
	// Delete everything on server
	deleteChannelsAndRoles(s, guildID)

	// Reset list of protected channels
	config.ServerConfig.ProtectedChannels = make(map[string]interface{})

	// Rename server
	_, err := s.GuildEdit(guildID, dg.GuildParams{Name: config.ServerConfig.EventName})
	util.ErrCheck(err, "Failed renaming Server")

	// Create Orga-Team role
	util.CreateRole(
		s, guildID, "Orga-Team", "0x9A58B4",
		config.ServerConfig.PermissionTemplates.OrgaTeam,
		true, true, &config.ServerConfig.OrgaTeamRoleID,
	)
	// Create RoBOT-Admin role
	util.CreateRole(
		s, guildID, "RoBOT-Admin", "0xFF0000",
		config.ServerConfig.PermissionTemplates.RoBOTAdmin,
		true, true, &config.ServerConfig.RoBOTAdminRoleID,
	)
	// Create role for each team
	log.Println("[Setup] Creating roles for teams...")
	for _, t := range config.TeamList {
		util.CreateRole(
			s, guildID, t.Name, t.TeamColor,
			config.ServerConfig.PermissionTemplates.TeamRole,
			true, true, &t.RoleID,
		)
	}
	// Create Participant role
	util.CreateRole(
		s, guildID, "Participant", "0x000000",
		config.ServerConfig.PermissionTemplates.Participant,
		false, false, &config.ServerConfig.ParticipantRoleID,
	)

	createBasicChannels(s, guildID, i.Member.Nick)
	sendRoleAssignmentMessage(s, config.ServerConfig.RoleAssignmentChannelID)

	config.SaveServerConfig()
	config.SaveTeamConfig()
}

// TODO This is stupid
const logCategory string = "Setup"

func setupRun(s *dg.Session, ev *dg.MessageCreate, args []string) {

	// Get guild ID and save to ServerConfig
	g, err := s.Guild(ev.GuildID)
	util.ErrCheck(err, "Error getting guild for ID "+ev.GuildID)
	config.ServerConfig.GuildID = g.ID

	getEveryoneRoleID(s, g)

	switch strings.ToLower(args[0]) {
	case "repair-roles":
		// Sanity checks for the config to see if manual entries are any good
		if !config.RoBotConfig.SanityCheck() {
			log.Println("[Setup/Repair-Roles] RoBotConfig.SanityCheck failed! Exiting...")
			_ = s.Close()
			os.Exit(1)
		}
		if !config.ServerConfig.SanityCheck(s) {
			log.Println("[Setup/Repair-Roles] ServerConfig.SanityCheck failed! Exiting...")
			_ = s.Close()
			os.Exit(1)
		}

		// Participant
		_, err = s.GuildRoleEdit(
			g.ID, config.ServerConfig.ParticipantRoleID, "Participant",
			0, false, config.ServerConfig.PermissionTemplates.Participant,
			false,
		)
		util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting Participant role!")

		// Everyone
		_, err = s.GuildRoleEdit(
			g.ID, config.ServerConfig.EveryoneRoleID, "",
			0, false, config.ServerConfig.PermissionTemplates.Everyone,
			true,
		)
		util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting @everyone role!")

		// RoBOT-Admin
		_, err = s.GuildRoleEdit(
			g.ID, config.ServerConfig.RoBOTAdminRoleID, "RoBOT-Admin",
			0xFF0000, false, config.ServerConfig.PermissionTemplates.RoBOTAdmin,
			true,
		)
		util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting RoBOT-Admin role!")

		// Orga-Team
		_, err = s.GuildRoleEdit(
			g.ID, config.ServerConfig.RoBOTAdminRoleID, "Orga-Team",
			0x9A58B4, true, config.ServerConfig.PermissionTemplates.OrgaTeam,
			true,
		)
		util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting @everyone role!")
	case "repair-channels":
		// TODO
		// Here:
		//   PER TEAM:
		//   - teamrole check
		//     - Name, RoleID, TeamzoneID present
		//     - Permission reset to zero
		//   - Teamzone permission check
		// - channel check
		//   - Name
		//   - Permission check
	case "add-team":
		// TODO
	case "add-channel":
		// TODO
	case "ramsg":
		sendRoleAssignmentMessage(s, config.ServerConfig.RoleAssignmentChannelID)
	}
}

// getEveryoneRoleID saves the ID for the @everyone role to the ServerConfig
func getEveryoneRoleID(s *dg.Session, g *dg.Guild) {
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

func deleteChannelsAndRoles(s *dg.Session, guildID string) {
	// TODO Add logging
	// Get all channels in server
	channels, err := s.GuildChannels(guildID)
	util.ErrCheck(err, "Failed to get channels for ID "+guildID)
	// Delete all channels
	for _, channel := range channels {
		_, err = s.ChannelDelete(channel.ID)
		util.ErrCheck(err, "Failed deleting channel "+channel.ID)
	}
	// Get all Roles
	roles, err := s.GuildRoles(guildID)
	for _, role := range roles {
		// If role ID is @everyone, we can't delete this, also we can't delete
		// the bot's own role, which should have Administrator permissions
		if role.ID == config.ServerConfig.EveryoneRoleID ||
			((role.Name == s.State.User.Username) && role.Permissions == 8) {
			continue
		}
		// Delete Role
		err = s.GuildRoleDelete(guildID, role.ID)
		util.ErrCheck(err, "Failed deleting role "+role.ID)
	}
	// Set server-wide permissions for @everyone
	_, err = s.GuildRoleEdit(
		guildID, config.ServerConfig.EveryoneRoleID, "",
		0, false, config.ServerConfig.PermissionTemplates.Everyone, true,
	)
	util.ErrCheck(err, "Failed setting permissions for @everyone role")
}

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
				s, guildID, fmt.Sprintf("Teamzone %02d", i), t.Name+" - Teamzone", catTeamzone.ID, dg.ChannelTypeGuildVoice,
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

func sendRoleAssignmentMessage(s *dg.Session, channelID string) {
	// Generate Embed with all teams
	var desc strings.Builder
	desc.WriteString("```text")
	for i, team := range config.TeamList {
		desc.WriteString(fmt.Sprintf("\n%2d: %s", i, team.Name))
	}
	desc.WriteString("\n```")

	embed := dg.MessageEmbed{
		Title: "Team Assignment",
		Description: fmt.Sprintf(
			"Greetings! You have reached the role assignment channel! "+
				"To get access to the digital venue, assign yourself a role using the `%steam TEAMNUMBER` command. "+
				"The Bot will then assign you to your team and grant you access to the rest of the server.",
			config.RoBotConfig.Prefix,
		),
		Color:  colors.GREEN,
		Footer: util.HelpEmbedFooter(),
		Image:  nil, // TODO There should be an image here
		Fields: []*dg.MessageEmbedField{
			{
				Name:  "Team List",
				Value: desc.String(),
			},
		},
	}

	_, err := s.ChannelMessageSendEmbed(channelID, &embed)
	util.ErrCheck(err, "[Setup] Failed sending role assignment message")
}
