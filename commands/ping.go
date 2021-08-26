package commands

import (
	"RoBOT/errors"
	"RoBOT/helptexts"

	"github.com/bwmarrin/discordgo"
)

func pingRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, "Pong!")
	errors.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}

func pingHelp(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, helptexts.DB["ping"])
	errors.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}
