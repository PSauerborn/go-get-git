package main

import (
	"fmt"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
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
	log.Debug(fmt.Sprintf("making request to %s", url))

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

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

// function used to check signature set
func isValidHookPayload(signature string, body []byte) bool {
	if (!strings.HasPrefix(signature, "sha1=")) {
		log.Error(fmt.Sprintf("received invalid hash format %s", signature))
		return false
	}
	actualHash := make([]byte, 20)
	hex.Decode(actualHash, []byte(signature[5:]))

	if hmac.Equal(signBody([]byte(GitHookSecret), body), actualHash) {
		log.Info("successfully validated git signature")
		return true
	} else {
		log.Error("unable to verify git signature")
		return false
	}
}

// function used to check git checks
func isMasterPushEvent(e *github.PushEvent) bool {
	if e.Ref != nil {
		return strings.HasPrefix(*e.Ref, "refs/heads/") && strings.HasSuffix(*e.Ref, "/master")
	}
	return false
}
