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
