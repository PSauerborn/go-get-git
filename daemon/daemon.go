package main

import (
	log "github.com/sirupsen/logrus"
)

var (
	daemon *GoGetGitDaemon
	logLevels = map[string]log.Level {"DEBUG": log.DebugLevel, "INFO": log.InfoLevel, "WARN": log.WarnLevel }
)

func main() {
	ConfigureService()

	// create new instance of daemon and run
	daemon := New()
	daemon.Run()
}