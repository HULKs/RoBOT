package config

// BotConf is the struct representing db/config.json
type BotConf struct {
	Prefix string `json:"PREFIX"`
	Token  string `json:"TOKEN"`
}

// TeamConf represents a json in db/teams
type TeamConf struct {
	Name      string `json:"Name"`
	RoleID    string `json:"RoleID"`
	TeamColor string `json:"TeamColor"`
}

// ServerConf represents db/server.json
type ServerConf struct {
	EventName            string `json:"EventName"`
	GuildID              string `json:"GuildID"`
	VoiceChannelCreateID string `json:"VoiceChannelCreateID"`
	ArchiveCategoryID    string `json:"ArchiveCategoryID"`
	EveryoneRoleID       string `json:"EveryoneRoleID"`
	ParticipantRoleID    string `json:"ParticipantRoleID"`
	OrgaTeamRoleID       string `json:"OrgaTeamRoleID"`
	BotAdminRoleID       string `json:"BotAdminRoleID"`
}
