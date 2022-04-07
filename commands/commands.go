package commands

import (
	"RoBOT/config"

	dg "github.com/bwmarrin/discordgo"
)

var Commands = []*dg.ApplicationCommand{
	{
		Name:        "ping",
		Description: "If online, the bot answers with: Pong!",
	},
	{
		Name:        "team",
		Description: "Assign yourself to a team",
		Options: []*dg.ApplicationCommandOption{
			{
				Name:        "team-number",
				Description: "Number of the team you want to assign yourself to",
				Type:        dg.ApplicationCommandOptionInteger,
				MinValue:    new(float64),
				MaxValue:    float64(len(config.TeamList)),
				Required:    true,
			},
		},
	},
	{
		Name:        "setup",
		Description: "Manage the server",
		Options: []*dg.ApplicationCommandOption{
			{
				Name:        "bootstrap",
				Description: "!DESTRUCTIVE! Bootstrap the server from scratch",
				Type:        dg.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "repair-roles",
				Description: "Reset role permissions to values specified in config",
				Type:        dg.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "role-assignment-message",
				Description: "Send the role assignment message to the role-assignment channel",
				Type:        dg.ApplicationCommandOptionSubCommand,
			},
		},
	},
	{
		Name:        "rename",
		Description: "Rename your event category, text- and voicechannels",
		Options: []*dg.ApplicationCommandOption{
			{
				Name:        "new-name",
				Description: "The new name for your channels",
				Type:        dg.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
	{
		Name:        "archive",
		Description: "Archive your event channel",
	},
	{
		Name:        "delete",
		Description: "Delete your event channel",
	},
}

var CommandHandlers = map[string]func(s *dg.Session, i *dg.InteractionCreate){
	"ping": func(s *dg.Session, i *dg.InteractionCreate) {
		s.InteractionRespond(
			i.Interaction, &dg.InteractionResponse{
				Type: dg.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: "Pong!",
				},
			},
		)
	},
	"team": func(s *dg.Session, i *dg.InteractionCreate) {
		team := config.TeamList[i.ApplicationCommandData().Options[0].IntValue()]
		assignMemberToTeam(s, i.Member, team, i.GuildID)
		s.InteractionRespond(
			i.Interaction, &dg.InteractionResponse{
				Type: dg.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: "Assigned you to team " + team.Name,
				},
			},
		)
	},
	"setup": func(s *dg.Session, i *dg.InteractionCreate) {
		switch i.ApplicationCommandData().Options[0].Name {
		case "bootstrap":
			s.InteractionRespond(
				i.Interaction, &dg.InteractionResponse{
					Type: dg.InteractionResponseChannelMessageWithSource,
					Data: &dg.InteractionResponseData{
						Content: "Bootstrapping server...",
					},
				},
			)
			setupBootstrap(s, i.GuildID, i)
		case "repair-roles":
			s.InteractionRespond(
				i.Interaction, &dg.InteractionResponse{
					Type: dg.InteractionResponseChannelMessageWithSource,
					Data: &dg.InteractionResponseData{
						Content: "Repairing roles...",
					},
				},
			)
			setupRepairRoles(s, i.GuildID, i)
		case "role-assignment-message":
			s.InteractionRespond(
				i.Interaction, &dg.InteractionResponse{
					Type: dg.InteractionResponseChannelMessageWithSource,
					Data: &dg.InteractionResponseData{
						Content: "Sending role-assignment message...",
					},
				},
			)
			sendRoleAssignmentMessage(s, config.ServerConfig.RoleAssignmentChannelID)
		default:
		}
	},
	"rename": func(s *dg.Session, i *dg.InteractionCreate) {
		newName := i.ApplicationCommandData().Options[0].StringValue()
		renameChannel(s, i.Interaction, newName)
		s.InteractionRespond(
			i.Interaction, &dg.InteractionResponse{
				Type: dg.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: "Renamed to: " + newName,
				},
			},
		)
	},
	"archive": func(s *dg.Session, i *dg.InteractionCreate) {
		s.InteractionRespond(
			i.Interaction, &dg.InteractionResponse{
				Type: dg.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: "Archiving channel...",
				},
			},
		)
		archiveChannel(s, i.Interaction)
	},
	"delete": func(s *dg.Session, i *dg.InteractionCreate) {
		s.InteractionRespond(
			i.Interaction, &dg.InteractionResponse{
				Type: dg.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Content: "Deleting channel...",
				},
			},
		)
		deleteChannel(s, i.Interaction)
	},
}
