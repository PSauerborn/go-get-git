package daemon

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "path/filepath"
    "github.com/PSauerborn/go-get-git/pkg/events"
    rabbit "github.com/PSauerborn/go-jackrabbit"
    log "github.com/sirupsen/logrus"
)

// function used to create new daemon
func New() *GoGetGitDaemon {
    ConfigureService()
    return &GoGetGitDaemon{}
}

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
                log.Error(fmt.Errorf("unable to process NewGitPush event: %v", err))
            }
            // handle event triggered when new application is registered
        case events.NewGitRepoEvent:
            log.Debug(fmt.Sprintf("processing new Git Application event %+v", e))
            err := handleNewApplicationEvent(e)
            if err != nil {
                log.Error(fmt.Errorf("unable to process NewGitRepo event: %v", err))
            }
            // handle default case
        default:
            log.Debug(fmt.Sprintf("received event type '%+v'", e))
        }
    }
}

// helper function used to create new directory for application
func handleNewApplicationEvent(event events.NewGitRepoEvent) error {
    log.Info(fmt.Sprintf("processing new application directory for %s", event.ApplicationDirectory))
    // create directory for new application
    err := os.Mkdir(event.ApplicationDirectory, 0775)
    if err != nil {
        log.Error(fmt.Errorf("unable to create new application directory: %v", err))
        return err
    }
    // clone git repository into given directory
    cmd := exec.Command(fmt.Sprintf("git clone %s.%s %s", event.RepoUrl, "git", event.ApplicationDirectory))
    stdout, err := cmd.Output()
    if len(stdout) > 0 {
        log.Info(stdout)
    }
    if err != nil {
        log.Error(fmt.Errorf("unable to clone git repo %s into directory %s: %v", event.RepoUrl, event.ApplicationDirectory, err))
        return err
    }
    return nil
}

// helper function used to handle new git push event
func handleGitPushEvent(event events.GitPushEvent) error {
    log.Info(fmt.Sprintf("processing new git push event for directory %s", event.ApplicationDirectory))
    // clone git repository into given directory
    cmd := exec.Command(fmt.Sprintf("git clone %s.%s %s", event.RepoUrl, "git", event.ApplicationDirectory))
    stdout, err := cmd.Output()
    if len(stdout) > 0 {
        log.Info(stdout)
    }
    if err != nil {
        log.Error(fmt.Errorf("unable to clone git repo %s into directory %s: %v", event.RepoUrl, event.ApplicationDirectory, err))
        return err
    }

    // find path of docker compose files in directory
    paths, err := findDockerCompose(event.ApplicationDirectory)
    if err != nil {
        log.Error(fmt.Errorf("unable to find docker-compose in directory %s: %v", event.ApplicationDirectory, err))
        return err
    }

    log.Debug(fmt.Sprintf("found %d docker-compose files to build", len(paths)))
    // iterate over path(s) of docker compose files and build docker files
    for _, path := range(paths) {
        log.Debug(fmt.Sprintf("building new docker compose file at %s", path))
        err := buildDockerComposeFile(path)
        if err != nil {
            log.Error(fmt.Errorf("unable to build docker-compose file at %s: %v", path, err))
        }
    }
    return nil
}

// helper function used to build new docker compose file
func buildDockerComposeFile(path string) error {
    // build docker compose file
    cmd := exec.Command(fmt.Sprintf("docker-compose run -f %s --build --remove-orphans -d", path))
    stdout, err := cmd.Output()
    if len(stdout) > 0 {
        log.Info(stdout)
    }
    return err
}

// helper function used to travers directory and find all docker compose files
func findDockerCompose(directory string) ([]string, error) {
    composeFiles := []string{}
    // walk through directory and find all docker compose files
    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }
        // add file to list of wanted files if it has the correct suffix
        if strings.HasSuffix(path, "docker-compose.yaml" ) || strings.HasSuffix(path, "docker-compose.yml") {
            composeFiles = append(composeFiles, path)
        }
        return nil
    })

    if err != nil {
        return composeFiles, err
    }
    return composeFiles, nil
}
