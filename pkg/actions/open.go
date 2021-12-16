package actions

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/dustinliu/nvclient/pkg/client"
	"github.com/neovim/go-client/nvim"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
)

func findWindowToActive() (nvim.Window, error) {
	nvc := client.Client()
	appFs := afero.NewOsFs()

	windows, err := nvc.Windows()
	if err != nil {
		return -1, cli.Exit(err, ClientError)
	}

	for _, w := range windows {
		buffer, err := nvc.WindowBuffer(w)
		if err != nil {
			return -1, cli.Exit(err, ClientError)
		}
		name, err := nvc.BufferName(buffer)
		if err != nil {
			return -1, cli.Exit(err, ClientError)
		}
		v, err := afero.Exists(appFs, name)
		if err != nil {
			continue
		}
		if v {
			return w, nil
		}
	}

	return windows[funk.MaxInt([]int{0, len(windows) - 1})], nil
}

func SpawnOpen(socket string, argv []string) error {
	e, err := exec.LookPath("nvim")
	if err != nil {
		logrus.Fatal("failed to locate then nvim  executable")
	}

	logrus.Debug("nvim executable: %v", e)

	err = syscall.Exec(e, append([]string{"nvim", "--listen", socket}, argv...),
		os.Environ())
	if err != nil {
		logrus.Fatalf("failed to start nvim, %v", err)
	}
	return nil
}

func RpcOpen(socket string, argv []string) error {
	nvc := client.Client()
	window, err := findWindowToActive()
	if err != nil {
		return cli.Exit(err, ClientError)
	}
	if err := nvc.SetCurrentWindow(window); err != nil {
		return cli.Exit(err, ClientError)
	}

	for i := 0; i < len(argv); i++ {
		file, err := filepath.Abs(argv[i])
		if err != nil {
			logrus.Fatal(err)
		}

		if err := nvc.Command("e " + file); err != nil {
			return cli.Exit(err, ClientError)
		}
	}

	return nil
}
