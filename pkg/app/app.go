package app

import (
	"os"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/dustinliu/nvclient/pkg/actions"
	"github.com/dustinliu/nvclient/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const (
	socket_env          = "NVIM_LISTEN_ADDRESS"
	server_name_env     = "NVC_SERVER_NAME"
	default_server_name = "nvim"
)

var (
	DebugFlag = &cli.BoolFlag{
		Name:    "debug",
		Usage:   "enable debug",
		Aliases: []string{"d"},
	}

	ServerFlag = &cli.StringFlag{
		Name:    "server",
		Usage:   "nvim server name",
		Aliases: []string{"s"},
	}

	flags = []cli.Flag{
		ServerFlag,
		DebugFlag,
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
	socket := os.Getenv(socket_env)
	server := os.Getenv(server_name_env)
	if socket == "" && server == "" {
		var err error
		socket, err = socketFile(default_server_name)
		if err != nil {
			return cli.Exit(err, actions.IOError)
		}
	} else if socket == "" {
		var err error
		socket, err = socketFile(server)
		if err != nil {
			return cli.Exit(err, actions.IOError)
		}
	}
	logrus.Debug("socket: %v\n", socket)

	fs := afero.NewOsFs()
	v, _ := afero.Exists(fs, socket)
	if v {
		return actions.RpcOpen(socket, c.Args().Slice())
	}
	return actions.SpawnOpen(socket, c.Args().Slice())
}

func socketFile(name string) (string, error) {
	socket, err := xdg.StateFile("nvclient/" + name + ".socket")
	if err != nil {
		return "", err
	}
	return socket, nil
}

func initApp(c *cli.Context) error {
	if DebugFlag.Value {
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
