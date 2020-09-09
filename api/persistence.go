package main

import (
	"fmt"
	"time"
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

func getRepoEntry(db *pgx.Conn, entryId uuid.UUID) (GitRepoEntry, error) {
	log.Debug(fmt.Sprintf("retrieving repo entry with ID %s", entryId))
	var (uid, repoUrl, accessToken string; createdAt time.Time)
	results := db.QueryRow(context.Background(), "SELECT entry_id,uid,repo_url,access_token,created_at FROM repo_entries WHERE entry_id=$1", entryId)
	err := results.Scan(&uid, &repoUrl, &accessToken, &createdAt)
	if err != nil {
		log.Error(fmt.Errorf("unable to fetch repo entries from database: %v", err))
		return GitRepoEntry{}, err
	}
	return GitRepoEntry{ EntryId: entryId, Uid: uid, RepoUrl: repoUrl, AccessToken: accessToken, CreatedAt: createdAt }, nil
}

func getRepoEntryByRepoUrl(db *pgx.Conn, url string) (GitRepoEntry, error) {
	log.Debug(fmt.Sprintf("retrieving repo entry for url %s", url))
	var (entryId uuid.UUID; uid, repoUrl, accessToken string; createdAt time.Time)
	results := db.QueryRow(context.Background(), "SELECT entry_id,uid,repo_url,access_token,created_at FROM repo_entries WHERE repo_url=$1", url)
	err := results.Scan(&entryId, &uid, &repoUrl, &accessToken, &createdAt)
	if err != nil {
		log.Error(fmt.Errorf("unable to fetch repo entries from database: %v", err))
		return GitRepoEntry{}, err
	}
	return GitRepoEntry{ EntryId: entryId, Uid: uid, RepoUrl: repoUrl, AccessToken: accessToken, CreatedAt: createdAt }, nil
}

func getAllRepoEntries(db *pgx.Conn) ([]GitRepoEntry, error) {
	log.Debug("retrieving all repo entries")
	values := []GitRepoEntry{}
	rows, err := db.Query(context.Background(), "SELECT entry_id,uid,repo_url,access_token,created_at FROM repo_entries")
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve repo entries: %v", err))
		return values, err
	}
	for rows.Next() {
		var (entryId uuid.UUID; uid, repoUrl, accessToken string; createdAt time.Time)
		err := rows.Scan(&entryId, &uid, &repoUrl, &accessToken, &createdAt)
		if err != nil {
			log.Error(fmt.Errorf("unable to process row: %v", err))
		} else {
			entry := GitRepoEntry{ EntryId: entryId, Uid: uid, RepoUrl: repoUrl, AccessToken: accessToken, CreatedAt: createdAt }
			values = append(values, entry)
		}
	}
	return values, nil
}

func getUserRepoEntries(db *pgx.Conn, uid string) ([]GitRepoEntry, error) {
	log.Debug(fmt.Sprintf("retrieving all repo entries for user %s", uid))
	values := []GitRepoEntry{}
	rows, err := db.Query(context.Background(), "SELECT entry_id,uid,repo_url,access_token,created_at FROM repo_entries WHERE uid=$1", uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve repo entries: %v", err))
		return values, err
	}
	for rows.Next() {
		var (entryId uuid.UUID; uid, repoUrl, accessToken string; createdAt time.Time)
		err := rows.Scan(&entryId, &uid, &repoUrl, &accessToken, &createdAt)
		if err != nil {
			log.Error(fmt.Errorf("unable to process row: %v", err))
		} else {
			entry := GitRepoEntry{ EntryId: entryId, Uid: uid, RepoUrl: repoUrl, AccessToken: accessToken, CreatedAt: createdAt }
			values = append(values, entry)
		}
	}
	return values, nil
}

func deleteRepoEntry(db *pgx.Conn, entryId uuid.UUID) error {
	log.Debug(fmt.Sprintf("deleting repo entry with ID %s", entryId))
	_, err := db.Exec(context.Background(), "DELETE FROM repo_entries WHERE entry_id = $1", entryId)
	if err != nil {
		log.Error(fmt.Errorf("unable to delete repo entry %s: %v", entryId, err))
		return err
	}
	return nil
}

func getHookEntry(db *pgx.Conn, hookId uuid.UUID) (GitHookEntry, error) {
	log.Debug(fmt.Sprintf("retrieving hook entry with ID %s", hookId))
	var (entryId uuid.UUID; created time.Time; meta interface{})
	hook := db.QueryRow(context.Background(), "SELECT entry_id,created_at,meta FROM git_hooks WHERE hook_id = $1", hookId)
	err := hook.Scan(&entryId, &created, &meta)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve git hook %s: %v", hookId, err))
		return GitHookEntry{}, err
	}
	return GitHookEntry{ EntryId: entryId, HookId: hookId, CreatedAt: created, Meta: meta }, nil
}

func getAllHookEntries(db *pgx.Conn) ([]GitHookEntry, error) {
	log.Debug("retrieving all hook entries")
	values := []GitHookEntry{}
	rows, err := db.Query(context.Background(), "SELECT entry_id,hook_id,created_at,meta FROM git_hooks")
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve repo entries: %v", err))
		return values, err
	}
	for rows.Next() {
		var (entryId, hookId uuid.UUID; created time.Time; meta interface{})
		err := rows.Scan(&entryId, &hookId, &created, &meta)
		if err != nil {
			log.Error(fmt.Errorf("unable to process row: %v", err))
		} else {
			entry := GitHookEntry{ EntryId: entryId, HookId: hookId, CreatedAt: created, Meta: meta }
			values = append(values, entry)
		}
	}
	return values, nil
}

func getAllHookEntriesByEntryId(db *pgx.Conn, entryId uuid.UUID) ([]GitHookEntry, error) {
	log.Debug("retrieving all hook entries")
	values := []GitHookEntry{}
	rows, err := db.Query(context.Background(), "SELECT entry_id,hook_id,created_at,meta FROM git_hooks WHERE entry_id = $1", entryId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve repo entries: %v", err))
		return values, err
	}
	for rows.Next() {
		var (entryId, hookId uuid.UUID; created time.Time; meta interface{})
		err := rows.Scan(&entryId, &hookId, &created, &meta)
		if err != nil {
			log.Error(fmt.Errorf("unable to process row: %v", err))
		} else {
			entry := GitHookEntry{ EntryId: entryId, HookId: hookId, CreatedAt: created, Meta: meta }
			values = append(values, entry)
		}
	}
	return values, nil
}

func deleteHookEntry(db *pgx.Conn, hookId uuid.UUID) error {
	log.Debug(fmt.Sprintf("deleting Git Hook with ID %s", hookId))
	_, err := db.Exec(context.Background(), "DELETE FROM git_hooks WHERE hook_id = $1", hookId)
	if err != nil {
		log.Error(fmt.Errorf("unable to delete repo entry %s: %v", hookId, err))
		return err
	}
	return nil
}