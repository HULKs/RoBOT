package commands

import (
	"log"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

func archiveChannel(s *dg.Session, i *dg.Interaction) {
	// TODO This should only archive one and only delete other channels when it's just one text channel left

	var err error

	// Get channel
	channel, err := s.Channel(i.ChannelID)
	util.ErrCheck(err, "[Archive] Failed getting channel for ID "+i.ChannelID)

	// Check if inside protected channel
	if config.IsProtected(channel) {
		log.Printf("[Archive] User %s tried archive in a protected channel %s !", i.Member.User.String(), channel.Name)
		err = util.SendProtectedCommandEmbed(s, channel.ID)
		util.CheckMsgSend(err, channel.Name)
		return
	}

	// ErrCheck if user has permission
	userPerms, err := s.State.UserChannelPermissions(i.Member.User.ID, channel.ID)
	util.ErrCheck(err, "[Archive] Failed getting permissions for user "+i.Member.User.String())
	if userPerms&dg.PermissionManageChannels != dg.PermissionManageChannels {
		log.Printf("[Archive] User %s has no permissions to archive channel!", i.Member.User.String())
		_, err = s.ChannelMessageSendEmbed(
			channel.ID, &dg.MessageEmbed{
				Title: "You don't have the necessary permissions to use this command. This incident will be reported.",
				Color: colors.RED,
				Footer: &dg.MessageEmbedFooter{
					Text: "If you think this is an error, contact the RoBOT-Admins",
				},
			},
		)
		util.CheckMsgSend(err, channel.ID)
		return
	}

	// Don't continue if we are already archived
	if channel.ParentID == config.ServerConfig.ArchiveCategoryID {
		log.Printf(
			"[Archive] User %s tried archiving channel %s but it's already archived", i.Member.User.String(),
			channel.Name,
		)
		return
	}

	// Get other channels in category (catChs)
	var catChs []*string
	gChs, err := s.GuildChannels(i.GuildID)
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
		channel.ID, &dg.ChannelEdit{
			ParentID: config.ServerConfig.ArchiveCategoryID,
		},
	)
	util.ErrCheck(err, "[Archive] Failed setting Archive category as parent for channel")
	log.Printf("[Archive] User %s moved channel %s to the archive", i.Member.User.String(), channel.Name)

	// Delete the rest
	for _, ch := range catChs {
		_, err = s.ChannelDelete(*ch)
		util.ErrCheck(err, "[Archive] Failed deleting channel "+*ch)
		log.Printf("[Archive] User %s deleted channel %d", i.Member.User.String(), ch)
	}
}
