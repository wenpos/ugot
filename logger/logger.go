package logger

import (
	"os"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("ugot")

var format = logging.MustStringFormatter(
	`[%{level:.8s}] [ugot] %{time:2006/01/02 15:04:05} %{callpath} %{shortfile}: %{message}`)

func GetLogger() *logging.Logger {
	backendFormat := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backendFormat, format)
	logging.SetBackend(backend2Formatter)
	return log
}
func CheckError(err error, msg string) {
	if err != nil {
		GetLogger().Error(msg + err.Error())
	}
}