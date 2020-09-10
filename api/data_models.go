package main


import (
	"time"
	"github.com/google/uuid"
)

type NewRegistryEntry struct {
	RepoName		string `json:"repo_name" binding:"required"`
	RepoUrl         string `json:"repo_url" binding:"required"`
	RepoOwner       string `json:"repo_owner" binding:"required"`
	RepoAccessToken string `json:"repo_access_token" binding:"required"`
}

type GitHookConfig struct {
	Url  		string `json:"url"`
	ContentType string `json:"content_type"`
	InsecureSSL int    `json:"insecure_ssl"`
	Secret 		string `json:"secret"`
}

type NewGitHookRequest struct {
	Name   string   	 `json:"name"`
	Active bool     	 `json:"active"`
	Events []string 	 `json:"events"`
	Config GitHookConfig `json:"config"`
}

type GitEventHookResponse struct {
	Ref string `json:"ref" binding:"required"`
}

// #########################################
// # Define data models used for persistence
// #########################################

type GitRepoEntry struct {
	EntryId     uuid.UUID `json:"entryId"`
	Uid 	    string    `json:"uid"`
	RepoUrl     string    `json:"repoUrl"`
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GitHookEntry struct {
	EntryId   uuid.UUID   `json:"entryId"`
	HookId    uuid.UUID   `json:"hookId"`
	CreatedAt time.Time   `json:"createdAt"`
	Meta      interface{} `json:"meta"`
}

type Event struct {
	ApplicationId  string      `json:"application_id"`
	ParentId	   uuid.UUID   `json:"parent_id"`
	EventId		   uuid.UUID   `json:"event_id"`
	EventTimestamp time.Time   `json:"event_timestamp"`
	EventPayload   interface{} `json:"event_payload"`
}

type GitPushEvent struct {
	RepoUrl	             string `json:"repo_url"`
	Uid	                 string `json:"uid"`
	ApplicationDirectory string `json:"application_directory"`
}

type NewGitRepoEvent struct {
	Uid	                 string `json:"uid"`
	ApplicationDirectory string `json:"application_directory"`
}