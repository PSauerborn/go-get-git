package main

import (
	"github.com/PSauerborn/go-get-git/pkg/daemon"
)

func main() {
	daemon := New()
	daemon.Run()
}