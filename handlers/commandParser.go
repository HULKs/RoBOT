package handlers

import (
	"strings"

	"RoBOT/config"
)

// ParseInput checks if the prefix is present, if yes, it returns a list of arguments
func ParseInput(input string) []string {
	if !strings.HasPrefix(input, config.RoBotConfig.Prefix) {
		return nil
	}
	input = strings.Replace(input, config.RoBotConfig.Prefix, "", 1)
	return strings.Split(input, " ")
}
