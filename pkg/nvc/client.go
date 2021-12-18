package nvc

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/neovim/go-client/nvim"
	"github.com/spf13/afero"
	"github.com/thoas/go-funk"
)

const (
	socket_env = "NVIM_LISTEN_ADDRESS"
	tmux_env   = "TMUX"
)

type Client struct {
	nvim   *nvim.Nvim
	once   sync.Once
	socket string
	fs     afero.Fs
}

func (nvc *Client) init() error {
	var e error
	nvc.once.Do(func() {
		nvc.fs = afero.NewOsFs()

		if err := nvc.initSocket(); err != nil {
			return
		}

		if nvc.socket == "" {
			LogEntry().Debug("no valid socket file, skip initializing the client")
			return
		}

		LogEntry().Debugf("connecting %v\n", nvc.socket)
		nvc.nvim, e = nvim.Dial(nvc.socket)
	})
	return e
}

func (nvc *Client) SetSocket(s string) {
	nvc.socket = s
}

func (nvc *Client) OpenFiles(files []string) error {
	if err := nvc.init(); err != nil {
		return err
	}

	LogEntry().Debugf("socket: %v\n", nvc.socket)

	v, _ := afero.Exists(nvc.fs, nvc.socket)
	if v {
		return nvc.rpcOpen(nvc.socket, files)
	}
	return nvc.spawnOpen(nvc.socket, files)
}

func (nvc *Client) Close() error {
	if nvc.nvim != nil {
		return nvc.nvim.Close()
	}
	return nil
}

func (nvc *Client) rpcOpen(socket string, argv []string) error {
	window, err := nvc.findWindowToActive()
	if err != nil {
		return err
	}
	if err := nvc.nvim.SetCurrentWindow(window); err != nil {
		return err
	}

	for i := 0; i < len(argv); i++ {
		file, err := filepath.Abs(argv[i])
		if err != nil {
			LogEntry().Fatal(err)
		}

		if err := nvc.nvim.Command("e " + file); err != nil {
			return err
		}
	}

	return nil
}

func (nvc *Client) spawnOpen(socket string, argv []string) error {
	e, err := exec.LookPath("nvim")
	if err != nil {
		LogEntry().Fatal("failed to locate then nvim  executable")
	}

	LogEntry().Debugf("nvim executable: %v", e)

	var args = []string{"nvim"}
	if len(socket) > 0 {
		args = append(args, "--listen", socket)
	}
	args = append(args, argv...)
	err = syscall.Exec(e, args, os.Environ())
	if err != nil {
		LogEntry().Fatalf("failed to start nvim, %v", err)
	}
	return nil
}

func (nvc *Client) findWindowToActive() (nvim.Window, error) {
	windows, err := nvc.nvim.Windows()
	if err != nil {
		return -1, err
	}

	for _, w := range windows {
		buffer, err := nvc.nvim.WindowBuffer(w)
		if err != nil {
			return -1, err
		}
		name, err := nvc.nvim.BufferName(buffer)
		if err != nil {
			return -1, err
		}
		v, err := afero.Exists(nvc.fs, name)
		if err != nil {
			continue
		}
		if v {
			return w, nil
		}
	}

	return windows[funk.MaxInt([]int{0, len(windows) - 1})], nil
}

func (nvc *Client) initSocket() error {
	var err error
	if nvc.socket != "" {
		LogEntry().Debug("socket has been set")
		return nil
	} else if InTmux() {
		nvc.socket, err = GetTmuxEnv(socket_env)
		if err != nil {
			LogEntry().Debugf("get tmux socket failed: [%v]", err)
			return err
		}
		LogEntry().Debugf("in tmux, socket: %s", nvc.socket)
	} else {
		nvc.socket = os.Getenv(socket_env)
		LogEntry().Debugf("env socket: [%v]\n", nvc.socket)
	}

	return err
}
