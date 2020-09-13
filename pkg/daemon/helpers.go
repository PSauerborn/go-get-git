package daemon

import (
	"os"
	"os/exec"
	"fmt"
	"github.com/PSauerborn/go-get-git/pkg/events"
	log "github.com/sirupsen/logrus"
)

// helper function used to create new directory for application
func handleNewApplicationEvent(event events.NewGitRepoEvent) error {
	// create directory for new application
	err := os.Mkdir(event.ApplicationDirectory, 0775)
	if err != nil {
		log.Error(fmt.Errorf("unable to create new application directory: %v", err))
		return err
	}
	// clone git repository into given directory
	cmd := exec.Command(fmt.Sprintf("git clone %s %s", event.RepoUrl, event.ApplicationDirectory))
	stdout, err := cmd.Output()
	log.Info(stdout)
	if err != nil {
		log.Error(fmt.Errorf("unable to clone git repo %s into directory %s: %v", event.RepoUrl, event.ApplicationDirectory, err))
		return err
	}
	return nil
}

// helper function used to handle new git push event
func handleGitPushEvent(event events.GitPushEvent) error {
	// clone git repository into given directory
	cmd := exec.Command(fmt.Sprintf("git clone %s %s", event.RepoUrl, event.ApplicationDirectory))
	stdout, err := cmd.Output()
	log.Info(stdout)
	if err != nil {
		log.Error(fmt.Errorf("unable to clone git repo %s into directory %s: %v", event.RepoUrl, event.ApplicationDirectory, err))
		return err
	}
	return nil
}

