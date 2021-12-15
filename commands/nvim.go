package commands

import (
	"os"
	"runtime"

	"github.com/neovim/go-client/nvim"
	"github.com/sirupsen/logrus"
)

var nvc *nvim.Nvim

func InitNvim() {
	if runtime.GOOS == "windows" {
		logrus.Fatalln("window is not supported yet")
	}

	socket := os.Getenv("NVIM_LISTEN_ADDRESS")
	var err error
	nvc, err = nvim.Dial(socket)
	if err != nil {
		logrus.Fatalf("connect to nvim failed, %v\n", err)
	}
}

func GetNvim() *nvim.Nvim {
	return nvc
}

func CloseNvim() {
	nvc.Close()
}
