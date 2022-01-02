package nvc

import (
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func TestInTmux(t *testing.T) {
	tmux := &TmuxImpl{}
	assert.False(t, tmux.InTmux())

	t.Setenv("TMUX", "tmux")
	assert.True(t, tmux.InTmux())
}

func Test_getCmd(t *testing.T) {
	full_tmux_path, err := exec.LookPath(tmux_cmd)
	if err != nil {
		t.Error("tmux executable not found")
	}
	tmux := &TmuxImpl{}
	cmd := tmux.getCmd("fdf")
	assert.Equal(t, full_tmux_path+" show-environment fdf", cmd.String())
}

func Test_getEnvValue(t *testing.T) {
	tmux := &TmuxImpl{}
	s := "aaa=b"
	assert.Equal(t, "b", tmux.splitEnvValue(s))

	s = "af=b=fdsf"
	assert.Equal(t, "b=fdsf", tmux.splitEnvValue(s))

	s = "af=b =fdsf"
	assert.Equal(t, "b =fdsf", tmux.splitEnvValue(s))

	s = ""
	assert.Equal(t, "", tmux.splitEnvValue(s))
}
