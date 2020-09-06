package main


import (

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