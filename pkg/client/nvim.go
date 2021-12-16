package client

import (
	"os"
	"sync"

	"github.com/neovim/go-client/nvim"
	"github.com/sirupsen/logrus"
)

var client *nvim.Nvim
var once sync.Once

func Client() *nvim.Nvim {
	once.Do(func() {
		socket := os.Getenv("NVIM_LISTEN_ADDRESS")
		if socket == "" {
			logrus.Fatal("not running nvim instance found")
		}

		logrus.Debugf("connect to %v\n", socket)
		var err error
		client, err = nvim.Dial(socket)
		if err != nil {
			logrus.Fatalf("failed to connect to nvim, %v", err)
		}
	})
	return client
}

func Close() error {
	return client.Close()
}
