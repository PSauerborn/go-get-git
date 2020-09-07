package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

// function used to start new authentication service
func main() {
	// read environment variables from config into local variables
	ConfigureService()
	router := gin.New()

	// configure GET routes used for server
	router.GET("/go-get-git/health", HealthCheck)
	router.GET("/go-get-git/registry", Persistence.Middleware(), GetRegistryEntries)
	router.GET("/go-get-git/registry/:entryId", Persistence.Middleware(), GetRegistryEntry)

	// configure POST routes used for server
	router.POST("/go-get-git/registry", Persistence.Middleware(), CreateRegistryEntry)
	router.POST("/go-get-git/webhook", HandleGitWebHook)
	// configure DELETE routes used for server
	router.DELETE("/go-get-git/registry/:entryId", Persistence.Middleware(), RemoveRegistryEntry)

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
	entryId, err := createRepoEntry(Persistence.Persistence(ctx), getUser(ctx), requestBody)
	if err != nil {
		log.Error(fmt.Sprintf("received invalid request body"))
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
	_, err = createHookEntry(Persistence.Persistence(ctx), entryId, body)
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
	StandardHTTP.FeatureNotSupported(ctx)
}

// API Hander used to get user registry entries
func GetRegistryEntries(ctx *gin.Context) {
	StandardHTTP.FeatureNotSupported(ctx)
}

// API Handler used to remove registry entry
func RemoveRegistryEntry(ctx *gin.Context) {
	StandardHTTP.FeatureNotSupported(ctx)
}

// function used to extract raw request body from request
func getRequestBody(ctx *gin.Context) ([]byte, error) {
	buffer := make([]byte, 1024)
	index, err := ctx.Request.Body.Read(buffer)
	if err != nil {
		return []byte(""), err
	}
	return buffer[0:index], nil
}

// API route used to handle git hooks
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
	StandardHTTP.Success(ctx)
}