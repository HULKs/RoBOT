package commands

import (
	"log"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

func deleteChannel(s *dg.Session, i *dg.Interaction) {
	// TODO This should only delete other channels when it's just one text channel left

	var err error

	// Get channel
	channel, err := s.Channel(i.ChannelID)
	util.ErrCheck(err, "[Archive] Failed getting channel for ID "+i.ChannelID)

	// Check if inside protected channel
	if config.IsProtected(channel) {
		log.Printf("[Delete] User %s tried delete in protected channel %s !", i.Member.User.String(), channel.Name)
		err = util.SendProtectedCommandEmbed(s, i.ChannelID)
		util.CheckMsgSend(err, i.GuildID, i.ChannelID)
		return
	}

	// ErrCheck if user has permission
	userPerms, err := s.State.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	util.ErrCheck(err, "[Delete] Failed getting permissions for user "+i.Member.User.String())
	if userPerms&dg.PermissionManageChannels != dg.PermissionManageChannels {
		log.Printf("[Delete] User %s has no permissions to delete channel!", i.Member.User.String())
		_, err = s.ChannelMessageSendEmbed(
			i.ChannelID, &dg.MessageEmbed{
				Title: "You don't have the necessary permissions to use this command. This incident will be reported.",
				Color: colors.RED,
				Footer: &dg.MessageEmbedFooter{
					Text: "If you think this is an error, contact the RoBOT-Admins",
				},
			},
		)
		util.CheckMsgSend(err, i.GuildID, i.ChannelID)
		return
	}

	// Don't continue if we are inside the archive
	if channel.ParentID == config.ServerConfig.ArchiveCategoryID {
		log.Printf(
			"[Delete] User %s tried deleting channel %s but it's in the archive!", i.Member.User.String(), channel.Name,
		)
		return
	}

	log.Printf("[Delete] User %s invoked DELETE in channel %s", i.Member.User.String(), channel.Name)

	// Get channels in category (catChs)
	var catChs []*string
	gChs, err := s.GuildChannels(i.GuildID)
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
		log.Printf("[Delete] User %s deleted channel %d", i.Member.User.String(), ch)
	}
}
