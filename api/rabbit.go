package main

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)


type Event struct {
	ApplicationId  uuid.UUID   `json:"application_id"`
	ParentId	   uuid.UUID   `json:"parent_id"`
	EventId		   uuid.UUID   `json:"event_id"`
	EventTimestamp time.Time   `json:"event_timestamp"`
	EventPayload   interface{} `json:"event_payload"`
}

type GitPushEvent struct {
	RepoUrl	             string `json:"repo_url"`
	Uid	                 string `json:"uid"`
	ApplicationDirectory string `json:"application_directory"`
}

func generateGitEvent(push GitPushEvent) Event {
	return Event{
		ApplicationId: ApplicationId,
		ParentId: uuid.New(),
		EventId: uuid.New(),
		EventTimestamp: time.Now(),
		EventPayload: push,
	}
}

// function used to process git event by sending message over rabbitmq server
func processGitPushEvent(ctx *gin.Context, e *github.PushEvent) {
	log.Info(fmt.Sprintf("received master push event for repo %s. sending message to worker", *e.Repo.URL))
	entry, err := getRepoEntryByRepoUrl(Persistence.Persistence(ctx), *e.Repo.URL)
	if err != nil {
		log.Error(fmt.Errorf("unable to get repo entry: %v", err))
	} else {
		log.Info(fmt.Sprintf("retrieved Repo Entry %+v", entry))
		event := generateGitEvent(GitPushEvent{ RepoUrl: entry.RepoUrl, Uid: entry.Uid, ApplicationDirectory: "" })
		sendRabbitPayload(event)
	}
}

// define function used to send message over rabbitmq server
func sendRabbitPayload(event Event) error {
	conn, err := amqp.Dial(RabbitQueueUrl)
	if err != nil {
		log.Error(fmt.Errorf("unable to connect to rabbitmq server: %s", err))
		return err
	}
	defer conn.Close()

	// create channel on rabbitmq server
	channel, err := conn.Channel()
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq channel: %s", err))
		return err
	}
	// declare events exchange with fanout type
	err = channel.ExchangeDeclare("events", "fanout", false, true, false, false, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to create rabbitmq exchange: %s", err))
		return err
	}
	// construct payload and send over rabbit server
	body, _ := json.Marshal(&event)
	payload := amqp.Publishing{ ContentType: "application/json", Body: []byte(body) }
	err = channel.Publish("events", "", false, false, payload)
	if err != nil {
		log.Error(fmt.Errorf("unable to send payload over rabbitmq server: %s", err))
		return err
	}
	log.Info(fmt.Sprintf("successfully sent payload %+v over rabbitMQ exchange", event))
	return nil
}


