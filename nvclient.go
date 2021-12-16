package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/neovim/go-client/nvim"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"

	"github.com/urfave/cli/v2"
)

var nvc *nvim.Nvim

var debug bool

var flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "debug",
		Aliases:     []string{"d"},
		Destination: &debug,
	},
}

func initApp(c *cli.Context) error {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if runtime.GOOS == "windows" {
		logrus.Fatalln("window is not supported yet")
	}

	InitNvim()

	return nil
}

func finish(c *cli.Context) error {
	nvc.Close()
	return nil
}

func findWindowToActive() nvim.Window {
	appFs := afero.NewOsFs()
	windows, err := nvc.Windows()
	if err != nil {
		logrus.Fatalf("failed to get windows")
	}

	for _, w := range windows {
		buffer, err := nvc.WindowBuffer(w)
		if err != nil {
			logrus.Fatalf("failed to get buffer from window")
		}
		name, err := nvc.BufferName(buffer)
		if err != nil {
			logrus.Fatalf("failed to get buffer name")
		}
		v, err := afero.Exists(appFs, name)
		if err != nil {
			continue
		}
		if v {
			return w
		}
	}

	return windows[funk.MaxInt([]int{0, len(windows) - 1})]
}

func open(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return cli.Exit("error, no filename provided", 1)
	}

	window := findWindowToActive()
	if err := nvc.SetCurrentWindow(window); err != nil {
		logrus.Fatalf("failed to set current window, %v\n", err)
	}

	for i := 0; i < c.Args().Len(); i++ {
		file, err := filepath.Abs(c.Args().Get(i))
		if err != nil {
			logrus.Fatal(err)
		}

		if err := nvc.Command("e " + file); err != nil {
			logrus.Fatalf("open file failed, %v", err)
		}
	}

	return nil
}

func main() {
	app := &cli.App{
		Usage:           "neovim remote client",
		Flags:           flags,
		Before:          initApp,
		After:           finish,
		Action:          open,
		HideHelpCommand: true,
	}

	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
