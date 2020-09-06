package main


import (
	log "github.com/sirupsen/logrus"
)

// define struct used to control daemon
type GoGetGitDaemon struct {
	RabbitMQConnection string
	ApiUrl			   string
}

func(daemon GoGetGitDaemon) Run() {
	log.Info("starting new instance of GoGetGit Daemon")

}

// function used to create new daemon
func New() *GoGetGitDaemon {
	return &GoGetGitDaemon{}
}

