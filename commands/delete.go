package commands

import (
	"log"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/helptexts"
	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

func deleteRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	// ErrCheck if user has permission
	userPerms, err := s.State.MessagePermissions(ev.Message)
	util.ErrCheck(err, "[Delete] Failed getting permissions for user "+ev.Author.Mention())
	if userPerms&discordgo.PermissionManageChannels != discordgo.PermissionManageChannels {
		log.Printf("[Delete] User %s has no permissions to delete channel!", ev.Author.Mention())
		_, err = s.ChannelMessageSendEmbed(
			ev.ChannelID, &discordgo.MessageEmbed{
				Title: "You don't have the necessary permissions to use this command. This incident will be reported.",
				Color: colors.RED,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "If you think this is an error, contact the RoBOT-Admins",
				},
			},
		)
		util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
		return
	}

	// Get channel
	channel, err := s.Channel(ev.ChannelID)
	util.ErrCheck(err, "[Archive] Failed getting channel for ID "+ev.ChannelID)

	// Don't continue if we are inside the archive
	if channel.ParentID == config.ServerConfig.ArchiveCategoryID {
		log.Printf(
			"[Delete] User %s tried deleting channel %s but it's in the archive!", ev.Author.Mention(), channel.Mention(),
		)
		return
	}

	log.Printf("[Delete] User %s invoked DELETE in channel %s", ev.Author.Mention(), channel.Mention())

	// Get channels in category (catChs)
	var catChs []*string
	gChs, err := s.GuildChannels(ev.GuildID)
	util.ErrCheck(err, "[Archive] Failed getting guild channels")
	for _, gch := range gChs {
		if gch.ParentID == channel.ParentID {
			catChs = append(catChs, &gch.ID)
		}
	}
	// Also append parent category to remove later
	catChs = append(catChs, &channel.ParentID)

	// Delete all channels and the category
	for _, ch := range catChs {
		_, err = s.ChannelDelete(*ch)
		util.ErrCheck(err, "[Delete] Failed deleting channel "+*ch)
		log.Printf("[Delete] User %s deleted channel %d", ev.Author.Mention(), ch)
	}
}

func deleteHelp(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, helptexts.DB["ping"])
	util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}
