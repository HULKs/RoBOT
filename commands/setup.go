package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HULKs/RoBOT/colors"
	"github.com/HULKs/RoBOT/config"
	"github.com/HULKs/RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

// TODO This is stupid
const logCategory string = "Setup"

func setupRun(s *dg.Session, ev *dg.MessageCreate, args []string) {
	// Get guild ID and save to ServerConfig
	g, err := s.Guild(ev.GuildID)
	util.ErrCheck(err, "Error getting guild for ID "+ev.GuildID)
	config.ServerConfig.GuildID = g.ID

	switch strings.ToLower(args[0]) {
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
	}
}

func setupBootstrap(s *dg.Session, guildID string, i *dg.InteractionCreate) {
	// Get role ID for @everyone
	saveEveryoneRoleID(s, guildID)
	// Delete all channels/roles and reset @everyone
	clearServerAndReset(s, guildID)

	// Reset list of protected channels
	config.ServerConfig.ProtectedChannels = make(map[string]interface{})

	// Rename server
	_, err := s.GuildEdit(guildID, &dg.GuildParams{Name: config.ServerConfig.EventName})
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

	createBasicChannels(s, guildID, i.Member.User.String())
	setupRoleAssignmentMessage(s, config.ServerConfig.RoleAssignmentChannelID)

	config.SaveServerConfig()
	config.SaveTeamConfig()
}

func setupRepairRoles(s *dg.Session, guildID string, i *dg.InteractionCreate) {
	var err error

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
		guildID, config.ServerConfig.ParticipantRoleID, &dg.RoleParams{
			Name:        "Participant",
			Color:       util.PointyInt(0),
			Hoist:       util.PointyBool(false),
			Permissions: &config.ServerConfig.PermissionTemplates.Participant,
			Mentionable: util.PointyBool(false),
		},
	)
	util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting Participant role!")

	// Everyone
	_, err = s.GuildRoleEdit(
		guildID, config.ServerConfig.EveryoneRoleID, &dg.RoleParams{
			Name:        "",
			Color:       util.PointyInt(0),
			Hoist:       util.PointyBool(false),
			Permissions: &config.ServerConfig.PermissionTemplates.Everyone,
			Mentionable: util.PointyBool(true),
		},
	)
	util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting @everyone role!")

	// RoBOT-Admin
	_, err = s.GuildRoleEdit(
		guildID, config.ServerConfig.RoBOTAdminRoleID, &dg.RoleParams{
			Name:        "RoBOT-Admin",
			Color:       util.PointyInt(0xFF0000),
			Hoist:       util.PointyBool(false),
			Permissions: &config.ServerConfig.PermissionTemplates.RoBOTAdmin,
			Mentionable: util.PointyBool(true),
		},
	)
	util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting RoBOT-Admin role!")

	// Orga-Team
	_, err = s.GuildRoleEdit(
		guildID, config.ServerConfig.RoBOTAdminRoleID, &dg.RoleParams{
			Name:        "Orga-Team",
			Color:       util.PointyInt(0x9A58B4),
			Hoist:       util.PointyBool(true),
			Permissions: &config.ServerConfig.PermissionTemplates.OrgaTeam,
			Mentionable: util.PointyBool(true),
		},
	)
	util.ErrCheck(err, "[Setup/Repair-Roles] Failed resetting @everyone role!")
}

func setupRoleAssignmentMessage(s *dg.Session, channelID string) {
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
			"Greetings! You have reached the team assignment channel! " +
				"To get access to the digital venue, select a team using the `/team` command. " +
				"The Bot will then assign you and grant you access to the rest of the server.",
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
