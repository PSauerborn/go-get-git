package main

import (
	"github.com/PSauerborn/go-get-git/pkg/daemon"
)

func main() {

	service := daemon.New()
	service.Run()
}