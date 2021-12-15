package main

import (
	"os"

	"github.com/dustinliu/nvclient/commands"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var (
	debug bool

	flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Destination: &debug,
		},
	}
)

func Init(c *cli.Context) error {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	commands.InitNvim()
	return nil
}

func Finish(c *cli.Context) error {
	commands.CloseNvim()
	return nil
}

func main() {
	app := &cli.App{
		Usage:  "neovim remote client",
		Flags:  flags,
		Before: Init,
		After:  Finish,
	}

	app.Setup()
	app.Commands = commands.GetCommand()

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
