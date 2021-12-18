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

func Test_getCmd(t *testing.T) {
	cmd := _getCmd("fdf")
	tmux, err := exec.LookPath(tmux_cmd)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, tmux+" show-environment fdf", cmd.String())
}

func Test_getEnvValue(t *testing.T) {
	s := "aaa=b"
	assert.Equal(t, "b", _splitEnvValue(s))

	s = "af=b=fdsf"
	assert.Equal(t, "b=fdsf", _splitEnvValue(s))

	s = "af=b =fdsf"
	assert.Equal(t, "b =fdsf", _splitEnvValue(s))

	s = ""
	assert.Equal(t, "", _splitEnvValue(s))
}
