package util

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"RoBOT/config"
	"RoBOT/errors"
)

// LoadJSON loads the json file into target
func LoadJSON(path string, target interface{}) {
	// Read []byte from file
	dat, err := ioutil.ReadFile(path)
	errors.Check(err, "Error reading from file")
	// Parse json
	err = json.Unmarshal(dat, &target)
	errors.Check(err, "Failed to parse json")
}

func ParseTeamColor(t *config.TeamConf) int {
	tc, err := strconv.ParseInt(t.TeamColor, 0, 0)
	errors.Check(err, "Failed to parse TeamColor for team "+t.Name)
	return int(tc)
}

func ParseHexColor(hex string) int {
	i, err := strconv.ParseInt(hex, 0, 0)
	errors.Check(err, "Failed to parse color from "+hex)
	return int(i)
}
