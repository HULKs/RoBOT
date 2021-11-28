package config

// BotConf is the struct representing db/config.json
type BotConf struct {
	Prefix string `json:"PREFIX"`
	Token  string `json:"TOKEN"`
}

// TeamConf represents a json in db/teams
type TeamConf struct {
	Name       string `json:"Name"`
	RoleID     string `json:"RoleID"`
	TeamzoneID string `json:"TeamzoneID"`
	TeamColor  string `json:"TeamColor"`
}

// ServerConf represents db/server.json
type ServerConf struct {
	EventName string `json:"EventName"`
	GuildID   string `json:"GuildID"`
	// Channels
	VoiceChannelCreateID    string                 `json:"VoiceChannelCreateID"`
	ArchiveCategoryID       string                 `json:"ArchiveCategoryID"`
	WelcomeChannelID        string                 `json:"WelcomeChannelID"`
	RoleAssignmentChannelID string                 `json:"RoleAssignmentChannelID"`
	ProtectedChannels       map[string]interface{} `json:"ProtectedChannels"`
	// Roles
	EveryoneRoleID    string `json:"EveryoneRoleID"`
	ParticipantRoleID string `json:"ParticipantRoleID"`
	OrgaTeamRoleID    string `json:"Orga-TeamRoleID"`
	RoBOTAdminRoleID  string `json:"RoBOT-AdminRoleID"`
	// Permissions
	PermissionTemplates struct {
		Everyone    int64 `json:"Everyone"`
		OrgaTeam    int64 `json:"Orga-Team"`
		Participant int64 `json:"Participant"`
		RoBOTAdmin  int64 `json:"RoBOT-Admin"`
		TeamRole    int64 `json:"TeamRole"`
	} `json:"PermissionTemplates"`
}
