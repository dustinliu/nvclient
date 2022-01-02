package nvc

import (
	"os"
	"os/exec"
	"strings"
)

const tmux_cmd = "tmux"

//go:generate  mockgen -source=tmux.go -destination=mock_tmux.go -package=nvc
type Tmux interface {
	InTmux() bool
	GetTmuxEnv(key string) (string, error)
}

type TmuxImpl struct {
}

func (tmux *TmuxImpl) InTmux() bool {
	in := os.Getenv("TMUX")
	LogEntry().Debugf("TMUX env: [%v]", in)
	return in != ""
}

func (tmux *TmuxImpl) GetTmuxEnv(key string) (string, error) {
	cmd := tmux.getCmd(key)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	r := tmux.splitEnvValue(string(out))
	return string(r), nil
}

func (tmux *TmuxImpl) getCmd(key string) *exec.Cmd {
	return exec.Command(tmux_cmd, "show-environment", key)
}

func (tmux *TmuxImpl) splitEnvValue(s string) string {
	r := strings.SplitN(string(s), "=", 2)
	if len(r) < 2 {
		LogEntry().Debug("split value: []")
		return ""
	}

	LogEntry().Debugf("split value: [%v]", r[1])
	return strings.TrimSpace(r[1])
}
