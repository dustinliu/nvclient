package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/neovim/go-client/nvim"
)

func find_file(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("can not get current directory")
	}

	return filepath.Join(wd, name)
}

func main() {
	if runtime.GOOS == "windows" {
		log.Fatalln("window is not supported yet")
	}

	socket := os.Getenv("NVIM_LISTEN_ADDRESS")
	nvclient, err := nvim.Dial(socket)
	if err != nil {
		log.Fatalf("connect to nvim failed, %v", err)
	}
	defer nvclient.Close()

	if err := nvclient.Call("edit", nil, find_file("go.mod")); err != nil {
		log.Fatalf("open file failed, %v", err)
	}
}
