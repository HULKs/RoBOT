package config

import (
	"github.com/bwmarrin/discordgo"
)

// GetPermOverParticipantDefault returns the PermissionOverwrites for channels
// that should be visible for @Participant, but not @everyone
func GetPermOverParticipantDefault() []*discordgo.PermissionOverwrite {
	// View/Connect for @Participant, None for @everyone
	return []*discordgo.PermissionOverwrite{
		{
			ID:   ServerConfig.EveryoneRoleID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel |
				discordgo.PermissionVoiceConnect,
			Allow: 0,
		},
		{
			ID:   ServerConfig.ParticipantRoleID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: 0,
			Allow: discordgo.PermissionViewChannel |
				discordgo.PermissionVoiceConnect,
		},
	}
}
