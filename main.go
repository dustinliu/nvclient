package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/dustinliu/nvclient/internal/nvc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	Success = iota
	ClientError
	NotSupportError
)

var version string

var (
	debugFlag = &cli.BoolFlag{
		Name:    "debug",
		Usage:   "enable debug",
		Aliases: []string{"d"},
	}

	versionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "print version",
		Aliases: []string{"v"},
	}

	socketFlag = &cli.StringFlag{
		Name:    "socket",
		Usage:   "nvim socket file",
		Aliases: []string{"s"},
	}

	flags = []cli.Flag{
		socketFlag,
		debugFlag,
		versionFlag,
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
	if c.IsSet(versionFlag.Name) {
		fmt.Println(version)
		return nil
	}

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
