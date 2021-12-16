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

const (
	SOCKET_ENV      = "NVIM_LISTEN_ADDRESS"
	SERVER_NAME_ENV = "NVC_SERVER_NAME"
	DEFAULT_SOCKET  = "nvclient"
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

func spawnNvim(name string) {
	e, err := exec.LookPath("nvim")
	if err != nil {
		logrus.Fatal("failed to locate then nvim  executable")
	}

	logrus.Debug("nvim executable: %v", e)

	err = syscall.Exec(e, []string{"nvim", "--listen", name, "go.mod"},
		os.Environ())
	if err != nil {
		logrus.Fatalf("failed to start nvim, %v", err)
	}
}

func socketFile(name string) string {
	return filepath.Join(os.TempDir(), "name.socket")
}

func OpenFile(c *cli.Context) error {
	socket := os.Getenv(SOCKET_ENV)
	if socket == "" {
		serverName := os.Getenv(SERVER_NAME_ENV)
		if serverName == "" {
			serverName = DEFAULT_SOCKET
		}
		spawnNvim(serverName)
		return nil
	}

	if c.Args().Len() < 1 {
		return nil
	}

	nvc := client.Client()
	window, err := findWindowToActive()
	if err != nil {
		return cli.Exit(err, ClientError)
	}
	if err := nvc.SetCurrentWindow(window); err != nil {
		return cli.Exit(err, ClientError)
	}

	for i := 0; i < c.Args().Len(); i++ {
		file, err := filepath.Abs(c.Args().Get(i))
		if err != nil {
			logrus.Fatal(err)
		}

		if err := nvc.Command("e " + file); err != nil {
			return cli.Exit(err, ClientError)
		}
	}

	return nil
}
