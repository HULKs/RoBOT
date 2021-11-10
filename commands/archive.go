package commands

import (
	"log"

	"RoBOT/config"
	"RoBOT/helptexts"
	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

func archiveRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	// ErrCheck if user has permission
	userPerms, err := s.State.MessagePermissions(ev.Message)
	util.ErrCheck(err, "[Archive] Failed getting permissions for user "+ev.Author.Username)
	if userPerms&discordgo.PermissionManageChannels != discordgo.PermissionManageChannels {
		log.Printf("[Archive] User %s has no permissions to delete channel!", ev.Author.Username)
		// TODO Send message to notify user
		return
	}
	// Get channel
	channel, err := s.Channel(ev.ChannelID)
	util.ErrCheck(err, "[Archive] Failed getting channel for ID "+ev.ChannelID)

	// Don't continue if we are already archived
	if channel.ParentID == config.ServerConfig.ArchiveCategoryID {
		log.Printf(
			"[Archive] User %s tried archiving channel %s but it's already archived", ev.Author.Username, channel.Name,
		)
		return
	}

	// Get other channels in category (catChs)
	var catChs []*string
	gChs, err := s.GuildChannels(ev.GuildID)
	util.ErrCheck(err, "[Archive] Failed getting guild channels")
	for _, gch := range gChs {
		if gch.ParentID == channel.ParentID && gch.ID != channel.ID {
			catChs = append(catChs, &gch.ID)
		}
	}
	// Also append parent category to remove later
	catChs = append(catChs, &channel.ParentID)

	// Move to archive
	_, err = s.ChannelEditComplex(
		channel.ID, &discordgo.ChannelEdit{
			ParentID: config.ServerConfig.ArchiveCategoryID,
		},
	)
	util.ErrCheck(err, "[Archive] Failed setting Archive category as parent for channel")
	log.Printf("[Archive] User %s moved channel %s to the archive", ev.Author.Username, channel.Name)

	// Delete the rest
	for _, ch := range catChs {
		_, err = s.ChannelDelete(*ch)
		util.ErrCheck(err, "[Archive] Failed deleting channel "+*ch)
		log.Printf("[Archive] User %s deleted channel %d", ev.Author.Username, ch)
	}
}

func archiveHelp(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, helptexts.DB["ping"])
	util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}
