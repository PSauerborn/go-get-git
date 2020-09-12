#!/bin/bash

export GO_GET_GIT_RABBIT_QUEUE_URL="amqp://guest:guest@localhost:5672/"
export GO_GET_GIT_QUEUE_NAME="go-get-git-daemon-events"
export GO_GET_GIT_EVENT_EXCHANGE_NAME="events"
export GO_GET_GIT_EXCHANGE_TYPE="fanout"