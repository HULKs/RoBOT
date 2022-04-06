package commands

import (
	"RoBOT/config"

	dg "github.com/bwmarrin/discordgo"
)

// func init() {
// 	commandMap["ping"] = Command{pingRun, pingHelp}
// 	commandMap["setup"] = Command{setupRun, setupHelp}
// 	commandMap["rename"] = Command{renameRun, renameHelp}
// 	commandMap["archive"] = Command{archiveRun, archiveHelp}
// 	commandMap["delete"] = Command{deleteRun, deleteHelp}
// 	commandMap["team"] = Command{teamRun, teamHelp}
// }

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
				Type:        dg.ApplicationCommandOptionInteger,
				Name:        "team-number",
				Description: "Number of the team you want to assign yourself to",
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
		},
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
		default:
		}
	},
}
