package commands

import (
	"fmt"
	"log"

	"RoBOT/config"
	"RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

func assignMemberToTeam(s *dg.Session, member *dg.Member, team *config.TeamConf, guildID string) {
	var err error

	// Remove all Roles from member
	for _, role := range member.Roles {
		err = s.GuildMemberRoleRemove(guildID, member.User.ID, role)
		util.ErrCheck(err, "[Team] Failed removing role from user "+member.User.ID)
	}

	// Reset Member Nickname
	err = s.GuildMemberNickname(guildID, member.User.ID, "")
	if err != nil {
		log.Printf("[Team] Failed resetting nickname for member: %s", member.User.Username)
	}

	// Change Nickname
	newUsername := fmt.Sprintf(
		"%.32s", fmt.Sprintf("[%.17s] %s", team.Name, member.User.Username),
	)
	err = s.GuildMemberNickname(guildID, member.User.ID, newUsername)
	if err != nil {
		log.Printf("[Team] Failed changing nickname for member: %s", member.User.Username)
	}

	// Assign Participant and Teamrole
	for _, roleID := range []string{
		config.ServerConfig.ParticipantRoleID,
		team.RoleID,
	} {
		err = s.GuildMemberRoleAdd(guildID, member.User.ID, roleID)
		util.ErrCheck(err, "[Team] Failed assigning roles for user "+member.User.Username)
	}
}
