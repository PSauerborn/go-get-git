package events

import (
	"errors"
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

var (
	InvalidEventError = errors.New("invalid go-get-git event payload")
)

type Event struct {
	ApplicationId  string      `json:"application_id"`
	ParentId	   uuid.UUID   `json:"parent_id"`
	EventId		   uuid.UUID   `json:"event_id"`
	EventTimestamp time.Time   `json:"event_timestamp"`
	EventType	   string      `json:"event_type"`
	EventPayload   interface{} `json:"event_payload"`
}

type GitPushEvent struct {
	RepoUrl	             string `json:"repo_url"`
	Uid	                 string `json:"uid"`
	ApplicationDirectory string `json:"application_directory"`
}

type BuildTriggeredEvent struct {
	EntryId uuid.UUID `json:"entry_id"`
	RepoUrl string	  `json:"repo_url"`
}

type BuildFailedEvent struct {
	EntryId uuid.UUID `json:"entry_id"`
	RepoUrl string	  `json:"repo_url"`
}

type BuildCompletedEvent struct {
	EntryId     uuid.UUID `json:"entry_id"`
	RepoUrl     string	  `json:"repo_url"`
	ContainerId string    `json:"container_id"`
}

type ContainerCrashedEvent struct {
	ContainerId string `json:"container_id"`
}

type ContainerRestartEvent struct {
	ContainerId string `json:"container_id"`
}

// function used ti parse event into Event objects
func ParseEvent(payload []byte) (*Event, error) {
	var e Event
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return &Event{}, InvalidEventError
	}

	switch e.EventType {
	case "GitPushEvent":
		var event GitPushEvent
		err := json.Unmarshal([]byte(e.EventPayload.(string)), &event)
		if err != nil {
			return &Event{}, InvalidEventError
		}
		e.EventPayload = event
		return &e, nil
	case "BuildTriggeredEvent":
		var event BuildTriggeredEvent
		err := json.Unmarshal([]byte(e.EventPayload.(string)), &event)
		if err != nil {
			return &Event{}, InvalidEventError
		}
		e.EventPayload = event
		return &e, nil
	case "BuildFailedEvent":
		var event BuildFailedEvent
		err := json.Unmarshal([]byte(e.EventPayload.(string)), &event)
		if err != nil {
			return &Event{}, InvalidEventError
		}
		e.EventPayload = event
		return &e, nil
	case "BuildCompletedEvent":
		var event BuildCompletedEvent
		err := json.Unmarshal([]byte(e.EventPayload.(string)), &event)
		if err != nil {
			return &Event{}, InvalidEventError
		}
		e.EventPayload = event
		return &e, nil
	case "ContainerCrashedEvent":
		var event ContainerCrashedEvent
		err := json.Unmarshal([]byte(e.EventPayload.(string)), &event)
		if err != nil {
			return &Event{}, InvalidEventError
		}
		e.EventPayload = event
		return &e, nil
	case "ContainerRestartEvent":
		var event ContainerRestartEvent
		err := json.Unmarshal([]byte(e.EventPayload.(string)), &event)
		if err != nil {
			return &Event{}, InvalidEventError
		}
		e.EventPayload = event
		return &e, nil
	default:
		return &Event{}, InvalidEventError
	}
}