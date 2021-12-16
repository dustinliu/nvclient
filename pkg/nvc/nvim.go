package nvc

import (
	"os"
	"sync"

	"github.com/neovim/go-client/nvim"
	"github.com/sirupsen/logrus"
)

const (
	socket_env = "NVIM_LISTEN_ADDRESS"
)

var client *nvim.Nvim
var once sync.Once

func Client() *nvim.Nvim {
	once.Do(func() {
		socket := SocketFile()
		if socket == "" {
			logrus.Debug("no socket file, skip initializing the client")
			return
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
	if client != nil {
		return client.Close()
	}
	return nil
}

func SocketFile() string {
	var socket string
	if socketFlag.IsSet() {
		socket = socketFlag.Value
	} else {
		socket = os.Getenv(socket_env)
	}
	return socket
}
