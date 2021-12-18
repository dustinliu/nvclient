package nvc

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func TestSocketWithEnv(t *testing.T) {
	t.Setenv(socket_env, "aaa")
	t.Setenv("TMUX", "")
	client := &Client{}
	client.init()

	assert.Equal(t, "aaa", client.socket)
}

func TestSocketWithArg(t *testing.T) {
	t.Setenv(socket_env, "aaa")
	client := &Client{}
	client.SetSocket("bbb")
	client.init()

	assert.Equal(t, "bbb", client.socket)
}
