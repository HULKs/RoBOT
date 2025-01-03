package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/HULKs/RoBOT/util"

	dg "github.com/bwmarrin/discordgo"
)

var (
	// RoBotConfig contains the currently active configuration options
	RoBotConfig BotConf
	// ServerConfig contains metadata about the current event on this server
	ServerConfig ServerConf
	// TeamList represents the teams for this server
	TeamList []*TeamConf

	// CLI flags
	configPath string
	help       bool

	// DB paths, filled in config.init()
	dbPath, dbConfigjson, dbServerjson, dbTeamsPath string
)

// TODO Add logging

func init() {
	// Parse CLI flags
	flag.StringVar(&configPath, "c", "", "Path pointing to the configuration db")
	flag.BoolVar(&help, "h", false, "Show help message")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Fill DB paths
	dbPath = configPath
	dbConfigjson = path.Join(dbPath, "config.json")
	dbServerjson = path.Join(dbPath, "server.json")
	dbTeamsPath = path.Join(dbPath, "teams")

	// Load the config db
	util.LoadJSON(dbConfigjson, &RoBotConfig)
	util.LoadJSON(dbServerjson, &ServerConfig)
	// Ensure ProtectedChannels is not nil
	if ServerConfig.ProtectedChannels == nil {
		ServerConfig.ProtectedChannels = make(map[string]interface{})
	}
	loadTeamConfigs(dbTeamsPath)
}

func loadTeamConfigs(dir string) {
	dirEntries, err := os.ReadDir(dir)
	util.ErrCheck(err, "Failed to ls "+dir)
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
	err = os.WriteFile(dbServerjson, conf, 0600)
	util.ErrCheck(err, fmt.Sprintf("Error writing %s", dbServerjson))
}

func SaveTeamConfig() {
	for _, team := range TeamList {
		conf, err := json.Marshal(team)
		util.ErrCheck(err, "Failed marshaling team "+team.Name)

		filepath := path.Join(dbTeamsPath, team.Name+".json")
		err = os.WriteFile(filepath, conf, 0600)
		util.ErrCheck(err, "Error writing "+filepath)
	}
}

// IsProtected returns true if channel or its parent is
// contained in ServerConfig.ProtectedChannels
func IsProtected(channel *dg.Channel) bool {
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
func (conf *ServerConf) SanityCheck(s *dg.Session) bool {
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
