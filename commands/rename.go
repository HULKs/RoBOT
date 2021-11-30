package commands

import (
	"log"
	"strings"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/helptexts"
	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

// TODO For some reason this only works twice per channel. Don't make too many typos I guess!

func renameRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	var err error

	// Get channel
	channel, err := s.Channel(ev.ChannelID)
	util.ErrCheck(err, "[Rename] Failed getting channel for "+ev.ChannelID)

	// Check if inside protected channel
	if config.IsProtected(channel) {
		log.Printf("[Rename] User %s tried rename in protected channel %s !", ev.Author.Username, channel.Name)
		err = util.SendProtectedCommandEmbed(s, ev.ChannelID)
		util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
		return
	}

	// ErrCheck if user has permission
	userPerms, err := s.State.MessagePermissions(ev.Message)
	util.ErrCheck(err, "[Rename] Failed getting permissions for user "+ev.Author.Username)
	if userPerms&discordgo.PermissionManageChannels != discordgo.PermissionManageChannels {
		log.Printf("[Rename] User %s has no permissions to rename channel!", ev.Author.Username)
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

	// Reassemble Title
	var builder strings.Builder
	for _, arg := range args {
		builder.WriteString(arg)
		builder.WriteRune(' ')
	}
	newName := strings.TrimSpace(builder.String())
	newNameText := strings.ToLower(strings.ReplaceAll(newName, " ", "_"))

	// Get Archive category for position
	catArchive, err := s.Channel(config.ServerConfig.ArchiveCategoryID)
	util.ErrCheck(err, "[Rename] Failed getting Archive category")

	// Get parent category
	parent, err := s.Channel(channel.ParentID)
	util.ErrCheck(err, "[Rename] Failed getting parent category")

	// Rename parent category
	_, err = s.ChannelEditComplex(
		parent.ID, &discordgo.ChannelEdit{
			Name:     newName,
			Position: catArchive.Position - 1,
		},
	)
	util.ErrCheck(err, "[Rename] Failed renaming category")

	// Get all channels in guild
	guildChannels, err := s.GuildChannels(ev.GuildID)
	util.ErrCheck(err, "[Rename] Failed getting channels for guild")
	// Rename all channels with this category as parent
	for _, gch := range guildChannels {
		if gch.ParentID == parent.ID {
			switch gch.Type {
			case discordgo.ChannelTypeGuildText:
				_, err = s.ChannelEdit(gch.ID, newNameText)
				util.ErrCheck(err, "[Rename] Failed renaming text channel "+gch.Name)
			case discordgo.ChannelTypeGuildVoice:
				_, err = s.ChannelEdit(gch.ID, newName)
				util.ErrCheck(err, "[Rename] Failed renaming voice channel "+gch.Name)
			}
		}
	}
}

func renameHelp(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, helptexts.DB["rename"])
	util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}
