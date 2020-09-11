package main

import (
	"github.com/PSauerborn/go-get-git/pkg/api"
)

func main() {
	service := api.New()
	service.Run()
}