package daemon

import (
	"fmt"
	"github.com/PSauerborn/go-get-git/pkg/events"
	log "github.com/sirupsen/logrus"
)

// define struct used to control daemon
type GoGetGitDaemon struct {
	RabbitConfig *RabbitMQConfig
}

// function used to create go-get-git daemon
func (daemon GoGetGitDaemon) Run() {
	log.Info("starting new instance of GoGetGit Daemon")
	// start listening on rabbitMQ queue for events
	err := ListenOnQueue(daemon.RabbitConfig, daemon.ProcessRabbitMessage)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create rabbitmq listener: %v", err))
	}
}

// rabbitMQ listener function
func (daemon GoGetGitDaemon) ProcessRabbitMessage(payload []byte) {
	log.Info(fmt.Sprintf("received rabbitmq message %v", string(payload)))
	event, err := events.ParseEvent(payload)
	if err != nil {
		log.Error(fmt.Errorf("unable to parse event: %s", err))
	} else {
		switch e := event.EventPayload.(type) {
		case events.GitPushEvent:
			log.Debug(fmt.Sprintf("processing new GitPushEvent %v", e))
		default:
			log.Debug(fmt.Sprintf("received event type '%+v'", e))
		}
	}
}

// function used to create new daemon
func New() *GoGetGitDaemon {
	ConfigureService()
	config := RabbitMQConfig{}.FromEnvironment()
	return &GoGetGitDaemon{config}
}

