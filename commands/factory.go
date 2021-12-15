package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{}

func RegisterCommand(command *cli.Command) {
	logrus.Debugf("register command: ", command.Name)
	commands = append(commands, command)
}

func GetCommand() []*cli.Command {
	return commands
}
