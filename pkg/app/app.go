package app

import (
	"runtime"

	"github.com/dustinliu/nvclient/pkg/actions"
	"github.com/dustinliu/nvclient/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	server_flag = &cli.BoolFlag{
		Name:    "debug",
		Usage:   "enable debug",
		Aliases: []string{"d"},
	}

	debug_flag = &cli.StringFlag{
		Name:    "server",
		Usage:   "nvim server name",
		Aliases: []string{"s"},
	}

	flags = []cli.Flag{
		server_flag,
		debug_flag,
	}
)

func NewApp() *cli.App {
	app := &cli.App{
		Usage:           "neovim remote client",
		Flags:           flags,
		Before:          initApp,
		After:           finish,
		Action:          run,
		HideHelpCommand: true,
	}
	app.Setup()
	return app
}

func run(c *cli.Context) error {
	return actions.OpenFile(c)
}

func initApp(c *cli.Context) error {
	if c.Bool(debug_flag.Name) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		return cli.Exit("Microsoft Windows is not supported yet",
			actions.NotSupportError)
	}

	return nil
}

func finish(c *cli.Context) error {
	return client.Close()
}
