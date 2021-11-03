package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"RoBOT/errors"

	"github.com/bwmarrin/discordgo"
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

func ParseTeamColor(name, hex string) int {
	tc, err := strconv.ParseInt(hex, 0, 0)
	errors.Check(err, "Failed to parse TeamColor for team "+name)
	return int(tc)
}

// ParseHexColor converts a hex RGB color of format "0x000000" code to an int
func ParseHexColor(hex string) (int, error) {
	i, err := strconv.ParseInt(hex, 0, 0)
	return int(i), err
}

// CreateRole creates a new role with the given properties
func CreateRole(
	s *discordgo.Session, g *discordgo.Guild, name, hexColor string, perm int64,
	hoist, mention bool, configRef *string,
) *discordgo.Role {
	// Create Role
	role, err := s.GuildRoleCreate(g.ID)
	errors.Check(err, "Failed to create Role \""+name+"\"")
	log.Printf("[%s] Created Role (%s)", role.ID, name)

	// Add Role ID to config
	if configRef != nil {
		*configRef = role.ID
	}
	log.Printf("[%s] Saved ID to config (%s)", role.ID, name)

	// Parse color
	col, err := ParseHexColor(hexColor)
	errors.Check(err, "Failed to parse color from "+hexColor+" for role "+name)

	// Edit Role, set name and permissions
	_, err = s.GuildRoleEdit(g.ID, role.ID, name, col, hoist, perm, mention)
	errors.Check(err, "Failed to edit Role \""+name+"\"")
	log.Printf(
		"[%s] Edited Role (Name: %s, Color: %s, Perm: %d, Hoist: %t, Mention: %t)",
		role.ID, name, hexColor, perm, hoist, mention,
	)

	return role
}

// CreateCategory creates a category with the given properties and returns the channel struct
func CreateCategory(
	s *discordgo.Session, g *discordgo.Guild, name, topic string, permissionOverwrites []*discordgo.PermissionOverwrite,
) *discordgo.Channel {
	// Create category
	category, err := s.GuildChannelCreateComplex(
		g.ID, discordgo.GuildChannelCreateData{
			Name:                 name,
			Type:                 discordgo.ChannelTypeGuildCategory,
			Topic:                topic,
			PermissionOverwrites: permissionOverwrites,
		},
	)
	errors.Check(err, "Failed creating category "+name)
	log.Printf("[%s] Created category: %s", category.ID, category.Name)

	return category
}

// CreateChannel creates a channel with the given properties and returns the channel
func CreateChannel(
	s *discordgo.Session, g *discordgo.Guild, name, topic, parentID string, chtype discordgo.ChannelType,
	permissionOverwrites []*discordgo.PermissionOverwrite,
) *discordgo.Channel {
	channel, err := s.GuildChannelCreateComplex(
		g.ID, discordgo.GuildChannelCreateData{
			Name:                 name,
			Topic:                topic,
			Type:                 chtype,
			ParentID:             parentID,
			PermissionOverwrites: permissionOverwrites,
		},
	)
	errors.Check(err, "Failed creating channel "+name)
	log.Printf("[%s] Created channel: %s", channel.ID, channel.Name)

	return channel
}
