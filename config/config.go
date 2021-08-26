package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"RoBOT/errors"
	"RoBOT/util"
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
	loadTeamConfigs("db/teams")
}

func loadTeamConfigs(dir string) {
	dirEntries, err := os.ReadDir(dir)
	errors.Check(err, "Failed to ls db")
	for _, file := range dirEntries {
		if !file.Type().IsRegular() {
			continue
		}
		tc := new(TeamConf)
		util.LoadJSON(path.Join(dir, file.Name()), tc)
		TeamList = append(TeamList, tc)
	}
}

func SaveServerConfig() {
	conf, err := json.Marshal(ServerConfig)
	errors.Check(err, "Failed to marshal ServerConfig")
	err = ioutil.WriteFile("db/server.json", conf, 0600)
	errors.Check(err, "Error writing db/server.json")
}

func SaveTeamConfig() {
	for _, team := range TeamList {
		conf, err := json.Marshal(team)
		errors.Check(err, "Failed marshaling team "+team.Name)

		filename := "db/teams/" + team.Name + ".json"
		err = ioutil.WriteFile(filename, conf, 0600)
		errors.Check(err, "Error writing "+filename)
	}
}
