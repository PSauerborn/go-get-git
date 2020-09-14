package events

import (
    "time"
    "github.com/google/uuid"
)

// function used to generate new event with a given payload
func New(EventType, ApplicationId string, ParentId uuid.UUID, payload interface{}) Event {
    event := Event{
        ApplicationId: ApplicationId,
        ParentId: ParentId,
        EventId: uuid.New(),
        EventTimestamp: time.Now(),
        EventType: EventType,
        EventPayload: payload,
    }
    return event
}