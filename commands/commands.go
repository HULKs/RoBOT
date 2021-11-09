package commands

var (
	// commandMap maps strings to associated commands
	commandMap = make(map[string]Command)
)

func init() {
	commandMap["ping"] = Command{pingRun, pingHelp}
	commandMap["setup"] = Command{setupRun, setupRun}
	commandMap["rename"] = Command{renameRun, renameHelp}
	commandMap["archive"] = Command{archiveRun, archiveHelp}
	commandMap["delete"] = Command{deleteRun, deleteHelp}
}

func GetCommand(str string) Command {
	return commandMap[str]
}
