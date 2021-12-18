package nvc

import (
	"os"
	"os/exec"
	"strings"
)

const tmux_cmd = "tmux"

type ICmd interface {
	Output() ([]byte, error)
}

func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

func GetTmuxEnv(key string) (string, error) {
	cmd := _getCmd(key)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	r := _splitEnvValue(string(out))
	return string(r), nil
}

func _getCmd(key string) *exec.Cmd {
	return exec.Command(tmux_cmd, "show-environment", key)
}

func _splitEnvValue(s string) string {
	r := strings.SplitN(string(s), "=", 2)
	if len(r) < 2 {
		LogEntry().Debug("split value: []")
		return ""
	}

	LogEntry().Debugf("split value: [%v]", r[1])
	return strings.TrimSpace(r[1])
}
