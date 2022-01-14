//go:build integration
// +build integration

package core_test

import (
	"context"
	"fmt"
	"github.com/hill-daniel/drizzle"
	"github.com/hill-daniel/drizzle/core"
	"github.com/hill-daniel/drizzle/env"
	"github.com/hill-daniel/drizzle/mocks"
	"github.com/hill-daniel/drizzle/shell"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestTestRunner_Run_should_clone_repo_and_execute_script(t *testing.T) {
	workDir, err := ioutil.TempDir("", "git_")
	assert.NoError(t, err)
	defer func() {
		err := cleanUp(workDir)
		assert.NoError(t, err)
	}()
	configParser := &env.ConfigParser{}
	executor := shell.NewExecutor("")
	assert.NoError(t, err)
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		return "", nil
	}}
	git := shell.NewGit(executor, secretsManager)
	repository := &drizzle.Repository{}
	repository.BranchRef = "refs/heads/main"
	repository.Name = "finance-scraper"
	repository.CloneURL = "https://github.com/hill-daniel/finance-scraper.git"
	job := core.NewJob(git, configParser, executor)

	err = job.Build(context.Background(), *repository, workDir)

	assert.NoError(t, err)
	_, err = os.Stat(fmt.Sprintf("%s/finance-scraper/finance-scraper", workDir))
	assert.NoError(t, err)
}

func cleanUp(dirs ...string) error {
	var cleanupFailed []error
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			cleanupFailed = append(cleanupFailed, err)
		}
	}
	if len(cleanupFailed) > 0 {
		return fmt.Errorf("failed to cleanup directories %v", cleanupFailed)
	}
	return nil
}
