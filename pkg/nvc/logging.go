package nvc

import (
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

func LogEntry() *logrus.Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		logrus.Warningln("Could not get context info for logger!")
	}
	filename := filepath.Base(file) + ":" + strconv.Itoa(line)
	return logrus.WithField("file", filename)
}
