package commands

import (
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
}
