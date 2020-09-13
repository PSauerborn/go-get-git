package daemon

import (
	"fmt"
	"github.com/PSauerborn/go-get-git/pkg/events"
	rabbit "github.com/PSauerborn/go-jackrabbit"
	log "github.com/sirupsen/logrus"
)

// define struct used to control daemon
type GoGetGitDaemon struct {}

// function used to create go-get-git daemon
func (daemon GoGetGitDaemon) Run() {
	log.Info("starting new instance of GoGetGit Daemon")
	config := rabbit.RabbitConnectionConfig{
		QueueURL: RabbitQueueUrl,
		QueueName: QueueName,
		ExchangeName: EventExchangeName,
		ExchangeType: ExchangeType,
	}
	// start listening on rabbitMQ queue for events
	err := rabbit.ListenOnQueueWithExchange(config, daemon.ProcessRabbitMessage)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create rabbitmq listener: %v", err))
	}
}

// function used to define how rabbitMQ messages are handled
func (daemon GoGetGitDaemon) ProcessRabbitMessage(payload []byte) {
	log.Info(fmt.Sprintf("received rabbitmq message %v", string(payload)))
	event, err := events.ParseEvent(payload)
	if err != nil {
		log.Error(fmt.Errorf("unable to parse event: %s", err))
	} else {

		// handle incoming event based on event type
		switch e := event.EventPayload.(type) {
			// handle event triggered when new master push is triggered on git repo
		case events.GitPushEvent:
			log.Debug(fmt.Sprintf("processing new GitPushEvent %+v", e))
			err := handleGitPushEvent(e)
			if err != nil {
				log.Error(fmt.Errorf("unable to process eventL %v", err))
			}
			// handle event triggered when new application is registered
		case events.NewGitRepoEvent:
			log.Debug(fmt.Sprintf("processing new Git Application event %+v", e))
			err := handleNewApplicationEvent(e)
			if err != nil {
				log.Error(fmt.Errorf("unable to process eventL %v", err))
			}
			// handle default case
		default:
			log.Debug(fmt.Sprintf("received event type '%+v'", e))
		}
	}
}

// function used to create new daemon
func New() *GoGetGitDaemon {
	ConfigureService()
	return &GoGetGitDaemon{}
}

