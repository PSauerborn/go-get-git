package main

import (
    "fmt"
    "os"
    "os/signal"
    "github.com/PSauerborn/go-get-git/pkg/daemon"
    log "github.com/sirupsen/logrus"
)

func main() {
    // create channel used for signal catching
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs)

    go func() {
        s := <-sigs
        log.Info(fmt.Sprintf("received signal %s", s))
        PostServiceHook()
        os.Exit(1)
    }()
    // create new service and start listening on rabbit queue
    service := daemon.New()
    service.Run()
}

func PostServiceHook() {
    log.Info("shutting down go-get-git deamon...")
}