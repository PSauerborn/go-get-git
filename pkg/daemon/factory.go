package daemon

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

// define struct used to control daemon
type GoGetGitDaemon struct {
	RabbitMQConnection string
	ApiUrl			   string
}

func (daemon GoGetGitDaemon) Run() {
	log.Info("starting new instance of GoGetGit Daemon")
	err := RabbitListener(daemon.ProcessRabbitMessage)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create rabbitmq worker: %v", err))
	}
}

func (daemon GoGetGitDaemon) ProcessRabbitMessage(payload []byte) {
	log.Info(fmt.Sprintf("received rabbitmq message %v", payload))
}

// function used to create new daemon
func New() *GoGetGitDaemon {
	return &GoGetGitDaemon{}
}

