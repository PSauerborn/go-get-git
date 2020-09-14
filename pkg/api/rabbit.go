package api

import (
    "fmt"
    "encoding/json"
    "github.com/PSauerborn/go-get-git/pkg/events"
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
    "github.com/google/go-github/github"
    log "github.com/sirupsen/logrus"
    rabbit "github.com/PSauerborn/go-jackrabbit"
)


// function used to process git event by sending message over rabbitmq server
func processGitPushEvent(ctx *gin.Context, e *github.PushEvent) {
    log.Info(fmt.Sprintf("received master push event for repo %s. sending message to worker", *e.Repo.URL))
    // get repo entry from database
    entry, err := persistence.getRepoEntryByRepoUrl(*e.Repo.URL)
    if err != nil {
        log.Error(fmt.Errorf("unable to get repo entry: %v", err))
    } else {
        log.Info(fmt.Sprintf("retrieved Repo Entry %+v", entry))
        // get file directory of application from database
        dir, err := persistence.getEntryDirectory(entry.EntryId)
        if err != nil {
            log.Error(fmt.Errorf("unable to fetch application directory: %s", err))
        } else {
            // generate rabbitMQ event and send over rabbit server to daemon
            payload := events.GitPushEvent{RepoUrl: *e.Repo.URL, ApplicationDirectory: dir}
            event := events.New("GitPushEvent", ApplicationId, uuid.New(), payload)
            sendRabbitPayload(event)
        }
    }
}

func processNewApplicationEvent(ctx *gin.Context, entryId uuid.UUID, user, application, url string) error {
    err := persistence.createEntryDirectory(entryId, BaseApplicationDirectory + application)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new application directory entry: %v", err))
        return err
    } else {
        // generate rabbitMQ event and send over rabbit server to daemon
        payload := events.NewGitRepoEvent{RepoUrl: url, ApplicationDirectory: BaseApplicationDirectory + application}
        event := events.New("NewGitRepoEvent", ApplicationId, uuid.New(), payload)
        sendRabbitPayload(event)
        return nil
    }
}

// define function used to send message over rabbitmq server
func sendRabbitPayload(event events.Event) error {

    RabbitConfig := rabbit.RabbitConnectionConfig{
        QueueURL: RabbitQueueUrl,
        ExchangeName: "events",
        ExchangeType: "fanout",
    }
    // construct payload and send over rabbit server
    body, _ := json.Marshal(&event)
    return rabbit.ConnectAndDeliverOverExchange(RabbitConfig, body)
}


