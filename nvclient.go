package main

import (
	"os"
	"runtime"

	"github.com/dustinliu/nvclient/client"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "debug",
		Usage:   "enable debug",
		Aliases: []string{"d"},
	},
	&cli.StringFlag{
		Name:    "server",
		Usage:   "nvim server name",
		Aliases: []string{"s"},
	},
}

func initApp(c *cli.Context) error {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		return cli.Exit("Microsoft Windows is not supported yet",
			client.CodeNotSupportError)
	}

	return nil
}

func finish(c *cli.Context) error {
	return client.Close()
}

func main() {
	app := &cli.App{
		Usage:           "neovim remote client",
		Flags:           flags,
		Before:          initApp,
		After:           finish,
		Action:          client.OpenFile,
		HideHelpCommand: true,
	}

	app.Setup()
	err := app.Run(os.Args)
	if err != nil {
		cli.HandleExitCoder(err)
	}
}
