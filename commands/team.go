package commands

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"RoBOT/colors"
	"RoBOT/config"
	"RoBOT/helptexts"
	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

// TODO Add logging
func teamRun(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	var (
		err         error
		msg         *discordgo.Message
		teamNum     int64
		member      *discordgo.Member
		team        config.TeamConf
		newUsername string
	)

	// Check if in role-assignment channel
	if ev.ChannelID != config.ServerConfig.RoleAssignmentChannelID {
		// Get role-assignment channel
		chRoleAs, err := s.Channel(config.ServerConfig.RoleAssignmentChannelID)
		util.ErrCheck(err, "[Team] Failed getting channel for config.ServerConfig.RoRoleAssignmentChannelID")

		_, err = s.ChannelMessageSendEmbed(
			ev.ChannelID, &discordgo.MessageEmbed{
				Title: fmt.Sprintf(
					"You can only use the `%srole` command in the %s channel!\nThis incident will be reported",
					config.RoBotConfig.Prefix,
					chRoleAs.Mention(),
				),
				Footer: util.HelpEmbedFooter(),
			},
		)
		util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
		return
	}

	// Check if args empty
	if len(args) == 0 {
		goto noValidNumber
	}

	// Check if team number is an int
	teamNum, err = strconv.ParseInt(args[0], 0, 32)
	if err != nil {
		goto noValidNumber
	}

	// Check if team number is in the list
	if int(teamNum) > len(config.TeamList) {
		goto noValidNumber
	}

	// Get member
	member, err = s.GuildMember(ev.GuildID, ev.Author.ID)
	util.ErrCheck(err, "[Team] Failed getting Member struct for User "+ev.Author.Username)

	// Remove all Roles from member
	for _, role := range member.Roles {
		err = s.GuildMemberRoleRemove(ev.GuildID, member.User.ID, role)
		util.ErrCheck(err, "[Team] Failed removing role from user "+member.User.ID)
	}

	// Change Username
	team = *config.TeamList[teamNum]
	newUsername = fmt.Sprintf(
		"%.32s", fmt.Sprintf("[%.17s] %s", team.Name, member.User.Username),
	)
	err = s.GuildMemberNickname(ev.GuildID, member.User.ID, newUsername)
	if err != nil {
		log.Printf("[Team] Failed changing Username for user %s", member.User.Username)
	}

	// Assign Participant and Teamrole
	err = s.GuildMemberEdit(
		ev.GuildID, member.User.ID, []string{
			config.ServerConfig.ParticipantRoleID,
			team.RoleID,
		},
	)
	util.ErrCheck(err, "[Team] Failed assigning roles for user "+member.User.Username)

	msg, err = s.ChannelMessageSendEmbed(
		ev.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Success! You have been assigned to team %s.", team.Name),
			Description: fmt.Sprintf(
				"Your nickname has been changed to include your team tag.\nYour new username is: `%s`",
				newUsername,
			),
			Color:  colors.GREEN,
			Footer: util.HelpEmbedFooter(),
		},
	)

	goto deleteBoth

noValidNumber:
	// Send error message to channel
	msg, err = s.ChannelMessageSendEmbed(
		ev.ChannelID, &discordgo.MessageEmbed{
			Title:       "You didn't specify a valid team number!",
			Description: fmt.Sprintf("The command should look something like `%steam 42`", config.RoBotConfig.Prefix),
			Color:       colors.RED,
			Footer:      util.HelpEmbedFooter(),
		},
	)
	util.CheckMsgSend(err, msg.GuildID, msg.ChannelID)

deleteBoth:
	// Delete ev and error message after 5 seconds
	go func(event *discordgo.MessageCreate, reply *discordgo.Message) {
		time.Sleep(time.Duration(time.Second * 10))
		err := s.ChannelMessagesBulkDelete(event.ChannelID, []string{event.Message.ID, reply.ID})
		util.CheckMsgSend(err, event.GuildID, event.ChannelID)
	}(ev, msg)
}

func teamHelp(s *discordgo.Session, ev *discordgo.MessageCreate, args []string) {
	_, err := s.ChannelMessageSend(ev.ChannelID, helptexts.DB["ping"])
	util.CheckMsgSend(err, ev.GuildID, ev.ChannelID)
}
