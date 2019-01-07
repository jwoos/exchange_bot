package main


import (
	"os"

	"github.com/op/go-logging"
)

var initialized = false


func initializeLogger(name string) *logging.Logger {
	if !initialized {
		format := logging.MustStringFormatter(
			`%{color}[%{time:01/02/06 15:04:05} %{level} %{shortfile}]%{color:reset} %{message}`,
		)
		backendStdin := logging.NewLogBackend(os.Stdin, "", 0)
		backendStdinFormatter := logging.NewBackendFormatter(backendStdin, format)
		logging.SetBackend(backendStdinFormatter)
		initialized = true
	}

	logger :=  logging.MustGetLogger(name)
	logging.SetLevel(LOG_LEVEL, name)

	return logger
}
