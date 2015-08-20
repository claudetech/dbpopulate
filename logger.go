package main

import (
	"github.com/claudetech/loggo"
	"os"
)

var logger = getLogger()

func getLogger() *loggo.Logger {
	logger := loggo.New("dbpopulate")
	logger.SetLevel(loggo.Info)
	logger.AddAppenderWithFilter(loggo.NewStdoutAppender(),
		&loggo.MaxLogLevelFilter{MaxLevel: loggo.Info}, loggo.Color)
	logger.AddAppenderWithFilter(loggo.NewStderrAppender(),
		&loggo.MinLogLevelFilter{MinLevel: loggo.Warning}, loggo.Color)
	if os.Getenv("DEBUG") != "" {
		logger.SetLevel(loggo.Debug)
	}
	return logger
}
