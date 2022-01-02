package nvc

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func TestSocketWithEnvNoFile(t *testing.T) {
	t.Setenv(socket_env, "/tmp/aaa")
	t.Setenv("TMUX", "")
	client := &Client{}
	client.init()

	assert.Equal(t, "", client.socket)
}

func TestSocketWithEnvWithFile(t *testing.T) {
	t.Setenv(socket_env, "/tmp/aaa")
	t.Setenv("TMUX", "")
	client := &Client{}
	client.fs = createFile(t, "/tmp/aaa")
	client.init()

	assert.Equal(t, "/tmp/aaa", client.socket)
}

func TestSocketWithEnvFileWithoutTmuxEnv(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := fs.Mkdir("/tmp", 0755); err != nil {
		t.Error("failed to create tmp dir")
	}
	if err := afero.WriteFile(fs, "/tmp/aaa", []byte("aaa"), 0644); err != nil {
		t.Error("failed to create socket file")
	}

	t.Setenv(socket_env, "/tmp/aaa")
	t.Setenv("TMUX", "tmux")
	client := &Client{}
	client.fs = fs
	client.init()

	assert.Equal(t, "/tmp/aaa", client.socket)
}

func TestSocketWithEnvNoFileWithTmuxFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tmux := NewMockTmux(ctrl)
	tmux.EXPECT().
		InTmux().
		Return(true)
	tmux.EXPECT().
		GetTmuxEnv(gomock.Any()).
		Return("/tmp/bbb", nil)

	t.Setenv(socket_env, "/tmp/aaa")
	t.Setenv("TMUX", "tmux")
	client := &Client{}
	client.fs = createFile(t, "/tmp/bbb")
	client.tmux = tmux
	client.init()

	assert.Equal(t, "/tmp/bbb", client.socket)
}

func TestSocketWithArg(t *testing.T) {
	t.Setenv(socket_env, "aaa")
	client := &Client{}
	client.SetSocket("bbb")
	client.init()

	assert.Equal(t, "bbb", client.socket)
}

func createFile(t *testing.T, file string) afero.Fs {
	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Errorf("failed to create dir")
		return nil
	}
	if err := afero.WriteFile(fs, file, []byte("aaa"), 0644); err != nil {
		t.Error("failed to create socket file")
		return nil
	}
	return fs
}
