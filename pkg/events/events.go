package events

import (
	"fmt"
	"errors"
	"time"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)

var (
	InvalidEventError = errors.New("invalid go-get-git event payload")
	validate = validator.New()
	parser EventParser
)

func ParseEvent(payload []byte) (*Event, error) {
	log.Info(fmt.Sprintf("received event payload %s", payload))
	parser := DefaultParser{}
	return parser.ParseEvent(payload)
}

type Event struct {
	ApplicationId  string      `json:"application_id" validate:"required"`
	ParentId	   uuid.UUID   `json:"parent_id" validate:"required"`
	EventId		   uuid.UUID   `json:"event_id" validate:"required"`
	EventTimestamp time.Time   `json:"event_timestamp" validate:"required"`
	EventType	   string      `json:"event_type" validate:"required"`
	EventPayload   interface{} `json:"event_payload" validate:"required"`
}

type GitPushEvent struct {
	RepoUrl	             string `json:"repo_url" validate:"required"`
	ApplicationDirectory string `json:"application_directory" validate:"required"`
}

type NewGitRepoEvent struct {
	RepoUrl 			 string	`json:"repo_url" validate:"required"`
	ApplicationDirectory string `json:"application_directory" validate:"required"`
}

type BuildTriggeredEvent struct {
	EntryId uuid.UUID `json:"entry_id" validate:"required"`
	RepoUrl string	  `json:"repo_url" validate:"required"`
}

type BuildFailedEvent struct {
	EntryId uuid.UUID `json:"entry_id" validate:"required"`
	RepoUrl string	  `json:"repo_url" validate:"required"`
}

type BuildCompletedEvent struct {
	EntryId     uuid.UUID `json:"entry_id" validate:"required"`
	RepoUrl     string	  `json:"repo_url" validate:"required"`
	ContainerId string    `json:"container_id" validate:"required"`
}

type ContainerCrashedEvent struct {
	ContainerId string `json:"container_id" validate:"required"`
}

type ContainerRestartEvent struct {
	ContainerId string `json:"container_id" validate:"required"`
}

// #######################################
// # Define interface used to parse events
// #######################################

type EventParser interface {
	ParseEvent(payload []byte) interface{}
}

type DefaultParser struct {}

// basic function used to parse an event into Event structs. Note that
// Event objects contaim a specific EventPayload struct, which must
// first be parsed. All Event JSON messages must contain an event_type
// which is parsed and returned
func(parser DefaultParser) ParseEvent(payload []byte) (*Event, error) {
	var e Event
	// parse generic JSON into Event struct and return error if event cannot be parsed
	err := json.Unmarshal(payload, &e)
	if err != nil || e.EventPayload == nil {
		log.Error("unable to parse event from JSON format")
		return &Event{}, InvalidEventError
	}

	log.Info(fmt.Sprintf("parsing event type '%s' with payload '%s'", e.EventType, e.EventPayload))
	var event interface{}

	// parse original payload back to JSON format to parse manually
	eventPayload, _ := json.Marshal(e.EventPayload)

	switch e.EventType {
	case "NewGitRepoEvent":
		event, err= parser.ParseNewGitRepoEvent(eventPayload)
	case "GitPushEvent":
		event, err = parser.ParseGitPushEvent(eventPayload)
	case "BuildTriggeredEvent":
		event, err = parser.ParseBuildTriggeredEvent(eventPayload)
	case "BuildFailedEvent":
		event, err = parser.ParseBuildFailedEvent(eventPayload)
	case "BuildCompletedEvent":
		event, err = parser.ParseBuildCompletedEvent(eventPayload)
	case "ContainerCrashedEvent":
		event, err = parser.ParseContainerCrashedEvent(eventPayload)
	case "ContainerRestartEvent":
		event, err = parser.ParseContainerRestartEvent(eventPayload)
	default:
		event, err = nil, InvalidEventError
	}

	if err != nil {
		log.Error(fmt.Errorf("unable to parse event: %+v", err))
		return nil, err
	}

	// assign parsed event payload as attribute of event
	log.Info(fmt.Sprintf("successfully parsed event %+v", event))
	e.EventPayload = event
	return &e, validate.Struct(event)
}

func(parser DefaultParser) ParseGitPushEvent(eventPayload []byte) (GitPushEvent, error) {
	var event GitPushEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}

func(parser DefaultParser) ParseBuildTriggeredEvent(eventPayload []byte) (BuildTriggeredEvent, error) {
	var event BuildTriggeredEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}

func(parser DefaultParser) ParseBuildFailedEvent(eventPayload []byte) (BuildFailedEvent, error) {
	var event BuildFailedEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}

func(parser DefaultParser) ParseBuildCompletedEvent(eventPayload []byte) (BuildCompletedEvent, error) {
	var event BuildCompletedEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}

func(parser DefaultParser) ParseContainerCrashedEvent(eventPayload []byte) (ContainerCrashedEvent, error) {
	var event ContainerCrashedEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}

func(parser DefaultParser) ParseContainerRestartEvent(eventPayload []byte) (ContainerRestartEvent, error) {
	var event ContainerRestartEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}

func(parser DefaultParser) ParseNewGitRepoEvent(eventPayload []byte) (NewGitRepoEvent, error) {
	var event NewGitRepoEvent
	err := json.Unmarshal(eventPayload, &event)
	return event, err
}