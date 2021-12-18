package main

import (
	"os"
	"runtime"

	"github.com/dustinliu/nvclient/pkg/nvc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	Success = iota
	ClientError
	NotSupportError
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

	client = &nvc.Client{}
)

func initApp(c *cli.Context) error {
	if c.Bool(debugFlag.Name) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		return cli.Exit("Microsoft Windows is not supported yet",
			NotSupportError)
	}

	return nil
}

func run(c *cli.Context) error {
	if c.IsSet(socketFlag.Name) {
		client.SetSocket(c.String(socketFlag.Name))
	}

	if err := client.OpenFiles(c.Args().Slice()); err != nil {
		return cli.Exit(err, ClientError)
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
		Action:          run,
		HideHelpCommand: true,
	}
	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		cli.HandleExitCoder(err)
	}
}
