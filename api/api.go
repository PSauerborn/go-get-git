package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// function used to start new authentication service
func main() {
	// read environment variables from config into local variables and connect persistence
	ConfigureService()
	ConnectPersistence()

	router := gin.New()

	// configure GET routes used for server
	router.GET("/go-get-git/health", HealthCheck)
	router.GET("/go-get-git/registry", GetRegistryEntries)
	router.GET("/go-get-git/registry/:entryId", GetRegistryEntry)
	router.GET("/go-get-git/hooks", GetHookEntries)
	router.GET("/go-get-git/hooks/:entryId", GetHookEntriesById)
	router.GET("/go-get-git/hook/:hookId", GetHookEntry)

	// configure POST routes used for server
	router.POST("/go-get-git/registry", CreateRegistryEntry)
	router.POST("/go-get-git/webhook", HandleGitWebHook)
	// configure DELETE routes used for server
	router.DELETE("/go-get-git/registry/:entryId", RemoveRegistryEntry)

	log.Info(fmt.Sprintf("starting go-get-git service at %s:%d", ListenAddress, ListenPort))
	router.Run(fmt.Sprintf("%s:%d", ListenAddress, ListenPort))
}

func getUser(ctx *gin.Context) string {
	return ctx.Request.Header.Get("X-Authenticated-Userid")
}

// function used as basic health check
func HealthCheck(ctx *gin.Context) {
	StandardHTTP.Success(ctx)
}

// API Handler used to create new registry entries
func CreateRegistryEntry(ctx *gin.Context) {
	log.Debug(fmt.Sprintf("received request to create registry entry for user %s", getUser(ctx)))
	var requestBody NewRegistryEntry
	err := ctx.ShouldBind(&requestBody)
	if err != nil {
		log.Error(fmt.Sprintf("received invalid request body"))
		StandardHTTP.InvalidRequestBody(ctx)
		return
	}
	log.Debug(fmt.Sprintf("processing request with body %+v", requestBody))
	// create new repo entry in database
	entryId, err := persistence.createRepoEntry(getUser(ctx), requestBody)
	if err != nil {
		log.Error(fmt.Sprintf("received invalid request body"))
		StandardHTTP.InternalServerError(ctx)
		return
	}
	// create new applications and process event
	err = processNewApplicationEvent(ctx, entryId, getUser(ctx), requestBody.RepoName)
	if err != nil {
		log.Error(fmt.Errorf("unable to process new application: %v", err))
		StandardHTTP.InternalServerError(ctx)
		return
	}
	// create new git hook on git server
	body, err := createGitWebHook(requestBody.RepoOwner, requestBody.RepoName, requestBody.RepoAccessToken)
	if err != nil {
		log.Error(fmt.Errorf("unable to create git hook: %v", err))
		StandardHTTP.InvalidRequest(ctx)
		return
	}
	// create new hook entry in database
	_, err = persistence.createHookEntry(entryId, body)
	if err != nil {
		log.Error(fmt.Sprintf("received invalid request body"))
		StandardHTTP.InternalServerError(ctx)
		return
	}
	response := gin.H{"http_code": 200, "success": true, "message": "successfully registered new repo"}
	ctx.JSON(200, response)
}

// API Handler used to get specific registry entry
func GetRegistryEntry(ctx *gin.Context) {
	entryId, err := uuid.Parse(ctx.Param("entryId"))
	if err != nil {
		log.Error(fmt.Sprintf("received invalid uuid %s", ctx.Param("entryId")))
		StandardHTTP.InvalidRequest(ctx)
		return
	}
	// get repo entry from database
	entry, err := persistence.getRepoEntry(entryId)
	if err != nil {
		switch err {
			// return 404 if no data entry exists
		case pgx.ErrNoRows:
			StandardHTTP.NotFound(ctx)
			return
			// by default, return internal server error
		default:
			StandardHTTP.InternalServerError(ctx)
			return
		}
	}
	ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": entry})
}

// API Hander used to get user registry entries
func GetRegistryEntries(ctx *gin.Context) {
	log.Debug(fmt.Sprintf("received request for registry entries from user %s", getUser(ctx)))
	// retrieve repo entries from database
	entries, err := persistence.getAllRepoEntries()
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": []GitRepoEntry{}})
			return
		default:
			StandardHTTP.InternalServerError(ctx)
			return
		}
	}
	ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": entries})
}

// API Handler used to remove registry entry
func RemoveRegistryEntry(ctx *gin.Context) {
	StandardHTTP.FeatureNotSupported(ctx)
}

// API route used to handle git hooks. Note that only Git Hooks
// that contain pushes to the master repositrory are handled and
// sent over the message bus
func HandleGitWebHook(ctx *gin.Context) {
	log.Info("received new git hook trigger")
	// validate git hook request
	payload, err := github.ValidatePayload(ctx.Request, []byte(GitHookSecret))
	if err != nil {
		log.Error(fmt.Errorf("unable to validate hook signature: %v", err))
		StandardHTTP.Forbidden(ctx)
		return
	}
	// parse event
	event, err := github.ParseWebHook(github.WebHookType(ctx.Request), payload)
	if err != nil {
		log.Error(fmt.Errorf("unable to parse webhook: %v", err))
		StandardHTTP.InternalServerError(ctx)
		return
	}
	log.Info(fmt.Sprintf("received event hook %+v", event))

	// check event type matches Push Event, and then check that push referse to master branch
	switch e := event.(type) {
	case *github.PushEvent:
		if isMasterPushEvent(e) {
			processGitPushEvent(ctx, e)
		} else {
			log.Info(fmt.Sprintf("received push event to non-master branch %s", *e.Ref))
		}
	default:
		log.Info(fmt.Sprintf("received non-push type event %v", e))
	}
	StandardHTTP.Success(ctx)
}

// API Route used to retrieve a particular hook entry by Hook ID
func GetHookEntry(ctx *gin.Context) {
	// parse hook ID into UUID format and return 400 if invalid
	hookId, err := uuid.Parse(ctx.Param("hookId"))
	if err != nil {
		log.Error(fmt.Sprintf("received invalid uuid %s", ctx.Param("entryId")))
		StandardHTTP.InvalidRequest(ctx)
		return
	}
	log.Debug(fmt.Sprintf("received request for hook with hook ID %s", hookId))
	// get hook entry from database. note that 404 if found if entry doesnt exist
	entry, err := persistence.getHookEntry(hookId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			StandardHTTP.NotFound(ctx)
			return
		default:
			StandardHTTP.InternalServerError(ctx)
			return
		}
	}
	ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": entry})
}

// API Route used to retrieve all hook entries currently stored in database
func GetHookEntries(ctx *gin.Context) {
	log.Debug(fmt.Sprintf("received request to fetch all hook entries from user %s", getUser(ctx)))
	// retrieve list of hook entries from database and return
	entries, err := persistence.getAllHookEntries()
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": []GitHookEntry{}})
			return
		default:
			StandardHTTP.InternalServerError(ctx)
			return
		}
	}
	ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": entries})
}

// API route used to retrieve all git hook entries that belong
// to a particular parent ID
func GetHookEntriesById(ctx *gin.Context) {
	entryId, err := uuid.Parse(ctx.Param("entryId"))
	if err != nil {
		log.Error(fmt.Sprintf("received invalid uuid %s", ctx.Param("entryId")))
		StandardHTTP.InvalidRequest(ctx)
		return
	}

	log.Debug(fmt.Sprintf("received request for hook entries for entry ID %s", entryId))
	entries, err := persistence.getAllHookEntriesByEntryId(entryId)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": []GitHookEntry{}})
			return
		default:
			StandardHTTP.InternalServerError(ctx)
			return
		}
	}
	ctx.JSON(200, gin.H{ "http_code": 200, "success": true, "payload": entries})
}