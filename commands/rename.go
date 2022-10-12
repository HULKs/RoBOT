package commands

import (
	"log"
	"strings"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

// TODO For some reason this only works twice per channel. Don't make too many typos I guess!

func renameChannel(s *dg.Session, i *dg.Interaction, newName string) {
	var err error

	// Get channel
	channel, err := s.Channel(i.ChannelID)
	util.ErrCheck(err, "[Rename] Failed getting channel for "+i.ChannelID)

	// Check if inside protected channel
	if config.IsProtected(channel) {
		log.Printf("[Rename] User %s tried rename in protected channel %s !", i.Member.User.String(), channel.Name)
		err = util.SendProtectedCommandEmbed(s, channel.ID)
		util.CheckMsgSend(err, channel.ID)
		return
	}

	// ErrCheck if user has permission
	userPerms, err := s.State.UserChannelPermissions(i.Member.User.ID, channel.ID)
	util.ErrCheck(err, "[Rename] Failed getting permissions for user "+i.Member.User.String())
	if userPerms&dg.PermissionManageChannels != dg.PermissionManageChannels {
		log.Printf("[Rename] User %s has no permissions to rename channel!", i.Member.User.String())
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

	// Adjust string for text channels
	newNameText := strings.ToLower(strings.ReplaceAll(newName, " ", "_"))

	// Get Archive category for position
	catArchive, err := s.Channel(config.ServerConfig.ArchiveCategoryID)
	util.ErrCheck(err, "[Rename] Failed getting Archive category")

	// Get parent category
	parent, err := s.Channel(channel.ParentID)
	util.ErrCheck(err, "[Rename] Failed getting parent category")

	// Rename parent category
	_, err = s.ChannelEditComplex(
		parent.ID, &dg.ChannelEdit{
			Name:     newName,
			Position: catArchive.Position - 1,
		},
	)
	util.ErrCheck(err, "[Rename] Failed renaming category")

	// Get all channels in guild
	guildChannels, err := s.GuildChannels(i.GuildID)
	util.ErrCheck(err, "[Rename] Failed getting channels for guild")
	// Rename all channels with this category as parent
	for _, gch := range guildChannels {
		if gch.ParentID == parent.ID {
			switch gch.Type {
			case dg.ChannelTypeGuildText:
				_, err = s.ChannelEdit(gch.ID, &dg.ChannelEdit{Name: newNameText})
				util.ErrCheck(err, "[Rename] Failed renaming text channel "+gch.Name)
			case dg.ChannelTypeGuildVoice:
				_, err = s.ChannelEdit(gch.ID, &dg.ChannelEdit{Name: newName})
				util.ErrCheck(err, "[Rename] Failed renaming voice channel "+gch.Name)
			}
		}
	}
}
