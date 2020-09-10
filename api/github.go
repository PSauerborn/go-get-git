package main

import (
	"fmt"
	"errors"
	"strings"
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

// function used to generate git ghook configuration
func getGitHookConfig() GitHookConfig {
	config := GitHookConfig{
		Url: GitHookUrl,
		ContentType: "json",
		InsecureSSL: 0,
		Secret: GitHookSecret,
	}
	return config
}

// function used to check if git repo is valid
func isValidGitRepo(repo string) bool {
	return true
}

// function used to create new git webhook when request is made
func createGitWebHook(owner, repo, token string) (NewGitHookRequest, error) {
	// create new git hook request object
	requestBody := NewGitHookRequest{
		Active: true,
		Events: []string{ "push" },
		Name: "web",
		Config: getGitHookConfig(),
	}
	requestBytes, _ := json.Marshal(&requestBody)

	// create instance of HTTP Client and add required headers
	client := &http.Client{}
	log.Debug(fmt.Sprintf("creating new hook for user %s with repo %s", owner, repo))
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/hooks", owner, repo)

	// format request with headers and set basic auth with given token
	request, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	request.Header.Add("accept", "application/vnd.github.v3+json")
	request.SetBasicAuth(owner, token)

	// make request to Github API to create new git hook
	resp, err := client.Do(request)
	if err != nil {
		log.Error(fmt.Errorf("unable to great new git hook: %v", err))
		return requestBody, err
	}
	defer resp.Body.Close()

	// parse request body if status is not 200 and return error
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Error(fmt.Sprintf("unable to create git hook: API returned code %d and body %s", resp.StatusCode, body))
		return requestBody, errors.New("unable to create new Githook")
	}
	return requestBody, nil
}

// function used check if git events are pushes to master branch
func isMasterPushEvent(e *github.PushEvent) bool {
	if e.Ref != nil {
		return strings.HasPrefix(*e.Ref, "refs/heads/") && strings.HasSuffix(*e.Ref, "/master")
	}
	return false
}
