package nvc

import (
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	debugFlag = &cli.BoolFlag{
		Name:    "debug",
		Usage:   "enable debug",
		Aliases: []string{"d"},
	}

	socketFlag = &cli.StringFlag{
		Name:    "socket",
		Usage:   "nvim socket file",
		Aliases: []string{"s"},
	}

	flags = []cli.Flag{
		socketFlag,
		debugFlag,
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
	return OpenFiles(c)
}

func initApp(c *cli.Context) error {
	if debugFlag.Value {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		return cli.Exit("Microsoft Windows is not supported yet",
			NotSupportError)
	}

	return nil
}

func finish(c *cli.Context) error {
	return client.Close()
}
