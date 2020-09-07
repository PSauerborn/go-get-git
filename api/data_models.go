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
	EntryId     uuid.UUID
	Uid 	    string
	RepoUrl     string
	AccessToken string
	CreatedAt   time.Time
}

type GitHookEntry struct {
	EntryId   uuid.UUID
	HookId    uuid.UUID
	CreatedAt time.Time
	Meta      interface{}
}