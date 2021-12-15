package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
	funk "github.com/thoas/go-funk"

	"github.com/neovim/go-client/nvim"
)

var nvc *nvim.Nvim
var appFs afero.Fs

func init() {
	appFs = afero.NewOsFs()
	if runtime.GOOS == "windows" {
		log.Fatalln("window is not supported yet")
	}

	socket := os.Getenv("NVIM_LISTEN_ADDRESS")
	var err error
	nvc, err = nvim.Dial(socket)
	if err != nil {
		log.Fatalf("connect to nvim failed, %v\n", err)
	}
}

func find_file(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("can not get current directory")
	}

	return filepath.Join(wd, name)
}

func findWindowToActive() nvim.Window {
	windows, err := nvc.Windows()
	if err != nil {
		log.Fatalf("failed to get windows")
	}

	for _, w := range windows {
		buffer, err := nvc.WindowBuffer(w)
		if err != nil {
			log.Fatalf("failed to get buffer from window")
		}
		name, err := nvc.BufferName(buffer)
		if err != nil {
			log.Fatalf("failed to get buffer name")
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

func main() {
	defer nvc.Close()

	window := findWindowToActive()
	if err := nvc.SetCurrentWindow(window); err != nil {
		log.Fatalf("failed to set current window, %v\n", err)
	}

	if err := nvc.Command("edit " + find_file(os.Args[1])); err != nil {
		log.Fatalf("open file failed, %v", err)
	}
}
