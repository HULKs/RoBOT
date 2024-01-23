package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/HULKs/RoBOT/colors"

	dg "github.com/bwmarrin/discordgo"
)

// LoadJSON loads the json file into target
func LoadJSON(path string, target interface{}) {
	// Read []byte from file
	dat, err := os.ReadFile(path)
	ErrCheck(err, "Error reading from file")
	// Parse json
	err = json.Unmarshal(dat, &target)
	ErrCheck(err, "Failed to parse json")
}

// ParseHexColor converts a hex RGB color of format "0x000000" code to an int
func ParseHexColor(hex string) (int, error) {
	i, err := strconv.ParseInt(hex, 0, 0)
	return int(i), err
}

// CreateRole creates a new role with the given properties.
func CreateRole(
	s *dg.Session, guildID, name, hexColor string, permissions int64,
	hoist, mentionable bool, configRef *string,
) *dg.Role {
	// Parse color
	color, err := ParseHexColor(hexColor)
	ErrCheck(err, "Failed to parse color from "+hexColor+" for role "+name)

	// Create Role
	role, err := s.GuildRoleCreate(
		guildID, &dg.RoleParams{
			Name:        name,
			Color:       &color,
			Hoist:       &hoist,
			Permissions: &permissions,
			Mentionable: &mentionable,
		},
	)
	ErrCheck(err, "Failed to create Role \""+name+"\"")
	log.Printf(
		"[%s] Created Role (Name: %s, Color: %s, Perm: %d, Hoist: %t, Mention: %t)",
		role.ID, name, hexColor, permissions, hoist, mentionable,
	)

	// Add Role ID to config
	if configRef != nil {
		*configRef = role.ID
	}
	log.Printf("[%s] Saved ID to config (%s)", role.ID, name)

	return role
}

// CreateCategory creates a category with the given properties and returns the channel struct
func CreateCategory(
	s *dg.Session, guildID, name, topic string, permissionOverwrites []*dg.PermissionOverwrite,
	logCategory, logUser string,
) *dg.Channel {
	// Create category
	category, err := s.GuildChannelCreateComplex(
		guildID, dg.GuildChannelCreateData{
			Name:                 name,
			Type:                 dg.ChannelTypeGuildCategory,
			Topic:                topic,
			PermissionOverwrites: permissionOverwrites,
		},
	)
	ErrCheck(err, "Failed creating category "+name)
	log.Printf("[%s] Created category: %s for %s", logCategory, category.Name, logUser)

	return category
}

// CreateChannel creates a channel with the given properties and returns the channel
func CreateChannel(
	s *dg.Session, guildID, name, topic, parentID string, channelType dg.ChannelType,
	permissionOverwrites []*dg.PermissionOverwrite, logCategory, logUser string,
) *dg.Channel {
	channel, err := s.GuildChannelCreateComplex(
		guildID, dg.GuildChannelCreateData{
			Name:                 name,
			Topic:                topic,
			Type:                 channelType,
			ParentID:             parentID,
			PermissionOverwrites: permissionOverwrites,
		},
	)
	ErrCheck(err, "Failed creating channel "+name)
	log.Printf("[%s] Created channel: %s for %s", logCategory, channel.Name, logUser)

	return channel
}

// PermOverwriteHideForAShowForB returns the PermissionOverwrites for channels
// that should be visible for A, but not B
func PermOverwriteHideForAShowForB(A, B string) []*dg.PermissionOverwrite {
	return []*dg.PermissionOverwrite{
		{
			ID:   A,
			Type: dg.PermissionOverwriteTypeRole,
			Deny: dg.PermissionViewChannel |
				dg.PermissionVoiceConnect,
			Allow: 0,
		},
		{
			ID:   B,
			Type: dg.PermissionOverwriteTypeRole,
			Deny: 0,
			Allow: dg.PermissionViewChannel |
				dg.PermissionVoiceConnect,
		},
	}
}

// ErrCheck panics with an error message if err != nil
func ErrCheck(err error, logMsg string) {
	if err != nil {
		log.Panicf("[ERROR] %s ERROR: %s", logMsg, err)
	}
}

// CheckMsgSend is a wrapper for util.ErrCheck() with a msg prefilled for ChannelMessageSend errors
func CheckMsgSend(err error, chName string) {
	ErrCheck(err, fmt.Sprintf("Failed sending message in channel: %s", chName))
}

// SendProtectedCommandEmbed sends a red Embed to the channel saying a user is using a
// command in a protected channel
func SendProtectedCommandEmbed(s *dg.Session, chID string) error {
	_, err := s.ChannelMessageSendEmbed(
		chID, &dg.MessageEmbed{
			Title: "You are trying to use a command in a protected channel. This command is reserved for meeting organizers. " +
				"This incident will be reported.",
			Color: colors.RED,
			Footer: &dg.MessageEmbedFooter{
				Text: "If you think this is an error, contact the RoBOT-Admins",
			},
		},
	)
	return err
}

// HelpEmbedFooter returns a reference to the default "if you need help" embed footer
func HelpEmbedFooter() *dg.MessageEmbedFooter {
	return &dg.MessageEmbedFooter{
		Text: "If you need help, contact the @Orga-Team or @RoBOT-Admin",
	}
}

// ContainsStr checks if a []string contains searchStr
func ContainsStr(s *[]string, searchStr *string) bool {
	for _, str := range *s {
		if *searchStr == str {
			return true
		}
	}
	return false
}

// PointyInt returns a pointer to an int with value i
func PointyInt(i int) *int {
	x := new(int)
	*x = i
	return x
}

// PointyBool returns a pointer to a bool with value b
func PointyBool(b bool) *bool {
	x := new(bool)
	*x = b
	return x
}
