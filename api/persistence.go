package main

import (
	"fmt"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

var (
	PostgresConnection = OverrideStringVariable("POSTGRES_CONNECTION", "postgres://postgres:postgres-dev@localhost:5432/postgres")
)

// function used to create new repository entry in database
func createRepoEntry(db *pgx.Conn, user string, body NewRegistryEntry) (uuid.UUID, error) {
	log.Debug(fmt.Sprintf("creating new registry entry %+v", body))
	entryId := uuid.New()
	// insert entry into database
	_, err := db.Exec(context.Background(), "INSERT INTO repo_entries(entry_id,uid,repo_url,access_token) VALUES($1,$2,$3,$4)", entryId, user, body.RepoUrl, body.RepoAccessToken)
	if err != nil {
		log.Error(fmt.Errorf("unable to insert values into users table: %v", err))
		return entryId, err
	}
	return entryId, nil
}

// function used to create new hook entry in database
func createHookEntry(db *pgx.Conn, entryId uuid.UUID, config NewGitHookRequest) (uuid.UUID, error) {
	log.Debug(fmt.Sprintf("creating new hook entry for entry %s", entryId))
	hookId := uuid.New()

	meta, _ := json.Marshal(&config)
	// insert entry into database
	_, err := db.Exec(context.Background(), "INSERT INTO git_hooks(hook_id,entry_id,meta) VALUES($1,$2,$3)", hookId, entryId, string(meta))
	if err != nil {
		log.Error(fmt.Errorf("unable to insert values into git hooks table: %v", err))
		return hookId, err
	}
	return hookId, nil
}