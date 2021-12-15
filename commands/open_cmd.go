package commands

import (
	"os"
	"path/filepath"

	"github.com/neovim/go-client/nvim"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
)

func init() {
	RegisterCommand(&cli.Command{
		Name:   "open",
		Usage:  "open file in nvim server",
		Action: open,
	})
}

func open(c *cli.Context) error {
	window := findWindowToActive()
	if err := nvc.SetCurrentWindow(window); err != nil {
		logrus.Fatalf("failed to set current window, %v\n", err)
	}

	if err := nvc.Command("e " + find_file(c.Args().Get(0))); err != nil {
		logrus.Fatalf("open file failed, %v", err)
	}
	return nil
}

func find_file(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("can not get current directory")
	}

	return filepath.Join(wd, name)
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
