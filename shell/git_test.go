package shell_test

import (
	"github.com/hill-daniel/drizzle"
	"github.com/hill-daniel/drizzle/mocks"
	"github.com/hill-daniel/drizzle/shell"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGit_Clone_should_return_absolute_clone_path(t *testing.T) {
	executor := &mocks.Executor{ExecuteFunc: func(command string, workDir string) (drizzle.ShellOut, error) {
		assert.True(t, strings.Contains(command, "https://github.com/user/cool-project.git"))
		return drizzle.ShellOut{}, nil
	}}
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		return "", nil
	}}
	git := shell.NewGit(executor, secretsManager)
	repository := drizzle.Repository{
		Name:     "cool-project",
		CloneURL: "https://github.com/user/cool-project.git",
	}

	clonePath, err := git.Clone(repository, "/tmp")

	assert.NoError(t, err)
	assert.Equal(t, "/tmp/cool-project", clonePath)
}

func TestGit_Clone_should_clone_private_repo(t *testing.T) {
	executor := &mocks.Executor{ExecuteFunc: func(command string, workDir string) (drizzle.ShellOut, error) {
		assert.True(t, strings.Contains(command, " https://drizzleGitUser:myPassword12345@github.com/user/cool-project.git"))
		return drizzle.ShellOut{}, nil
	}}
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		assert.Equal(t, "GITHUB_REPOSITORY_12345", secretID)
		return "myPassword12345", nil
	}}
	git := shell.NewGit(executor, secretsManager)
	repository := drizzle.Repository{
		ID:       "12345",
		Name:     "cool-project",
		CloneURL: "https://github.com/user/cool-project.git",
		Private:  true,
	}

	clonePath, err := git.Clone(repository, "/tmp")

	assert.NoError(t, err)
	assert.Equal(t, "/tmp/cool-project", clonePath)
}
