package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"RoBOT/util"

	"github.com/bwmarrin/discordgo"
)

var (
	// RoBotConfig contains the currently active configuration options
	RoBotConfig BotConf
	// ServerConfig contains metadata about the current event on this server
	ServerConfig ServerConf
	// TeamList represents the teams for this server
	TeamList []*TeamConf
)

// TODO Add logging
// TODO the path should be a cmdline flag

func init() {
	// Load the config db
	util.LoadJSON("db/config.json", &RoBotConfig)
	util.LoadJSON("db/server.json", &ServerConfig)
	// Ensure ProtectedChannels is not nil
	if ServerConfig.ProtectedChannels == nil {
		ServerConfig.ProtectedChannels = make(map[string]interface{})
	}
	loadTeamConfigs("db/teams")
}

func loadTeamConfigs(dir string) {
	dirEntries, err := os.ReadDir(dir)
	util.ErrCheck(err, "Failed to ls db")
	for _, file := range dirEntries {
		if !file.Type().IsRegular() {
			continue
		}
		tc := new(TeamConf)
		util.LoadJSON(path.Join(dir, file.Name()), tc)
		TeamList = append(TeamList, tc)
	}
	sort.Slice(
		TeamList, func(i, j int) bool {
			return TeamList[i].Name < TeamList[j].Name
		},
	)
}

func SaveServerConfig() {
	conf, err := json.Marshal(ServerConfig)
	util.ErrCheck(err, "Failed to marshal ServerConfig")
	err = ioutil.WriteFile("db/server.json", conf, 0600)
	util.ErrCheck(err, "Error writing db/server.json")
}

func SaveTeamConfig() {
	for _, team := range TeamList {
		conf, err := json.Marshal(team)
		util.ErrCheck(err, "Failed marshaling team "+team.Name)

		filename := "db/teams/" + team.Name + ".json"
		err = ioutil.WriteFile(filename, conf, 0600)
		util.ErrCheck(err, "Error writing "+filename)
	}
}

// IsProtected returns true if channel or its parent is
// contained in ServerConfig.ProtectedChannels
func IsProtected(channel *discordgo.Channel) bool {
	for ID, _ := range ServerConfig.ProtectedChannels {
		if ID == channel.ID || ID == channel.ParentID {
			return true
		}
	}
	return false
}

// SanityCheck checks if the Token and Prefix members have been set (!= "").
// This does not mean the Token or Prefix are valid!
func (conf *BotConf) SanityCheck() bool {
	if conf.Token == "" || conf.Prefix == "" {
		return false
	}
	return true
}

// SanityCheck checks if the members of ServerConf have been set to valid values.
func (conf *ServerConf) SanityCheck(s *discordgo.Session) bool {
	var err error

	// Event Name
	if conf.EventName == "" {
		// TODO
		return false
	}

	// GuildID
	_, err = s.Guild(conf.GuildID)
	if err != nil {
		// TODO
		return false
	}

	// Special channels
	channels := []string{
		conf.VoiceChannelCreateID,
		conf.ArchiveCategoryID,
		conf.WelcomeChannelID,
		conf.RoleAssignmentChannelID,
	}
	for _, chID := range channels {
		_, err = s.Channel(chID)
		if err != nil {
			// TODO
			return false
		}
	}

	// Protected channels
	for chID, _ := range conf.ProtectedChannels {
		_, err = s.Channel(chID)
		if err != nil {
			// TODO
			return false
		}
	}

	// Roles
	guildRoles, err := s.GuildRoles(conf.GuildID)
	util.ErrCheck(err, "") // TODO

	roles := []string{
		conf.EveryoneRoleID,
		conf.ParticipantRoleID,
		conf.OrgaTeamRoleID,
		conf.RoBOTAdminRoleID,
	}
	// Generate slice of all guild roles
	var guildRoleIDs []string
	for _, guildRole := range guildRoles {
		guildRoleIDs = append(guildRoleIDs, guildRole.ID)
	}
	// Check if all roles are in GuildRoleIDs
	for _, role := range roles {
		if util.ContainsStr(&guildRoleIDs, &role) {
			continue
		}
		// TODO
		return false
	}

	// TODO PermissionTemplates?

	return true
}
