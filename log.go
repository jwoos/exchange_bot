package main


import (
	"os"

	"github.com/op/go-logging"
)


func initializeLogger(name string) *logging.Logger {
	format := logging.MustStringFormatter(
		`%{color}[%{time:01/02/06 15:04:05} %{level} %{shortfile}]%{color:reset} %{message}`,
	)
	backend1 := logging.NewLogBackend(os.Stdin, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)

	return logging.MustGetLogger(name)
}
