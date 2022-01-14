package shell

import (
	"fmt"
	"github.com/hill-daniel/drizzle"
	"github.com/pkg/errors"
	"log"
	"strings"
)

const (
	gitHubSecretsPrefix     = "GITHUB"
	repositorySecretsPrefix = "REPOSITORY"
	defaultGitUser          = "drizzleGitUser"
)

type credentials struct {
	User     string
	Password string
}

// Git represents an interface to executing git commands.
type Git struct {
	executor       drizzle.Executor
	secretsManager drizzle.SecretRetriever
}

// NewGit creates a new Git.
func NewGit(executor drizzle.Executor, secretsManager drizzle.SecretRetriever) *Git {
	return &Git{
		executor:       executor,
		secretsManager: secretsManager,
	}
}

// Clone git clones the branch of the repository into given workDir and returns the absolute path.
// E.g. workdir is /tmp, repository name is cool-project the resulting clonePath will be /tmp/cool-project.
func (g *Git) Clone(repository drizzle.Repository, workDir string) (string, error) {
	cloneURL := repository.CloneURL
	if repository.Private {
		credentials, err := createCredentials(repository, g.secretsManager)
		if err != nil {
			return "", errors.Wrapf(err, "failed to clone private repository; no credentials given")
		}
		cloneURL = createPrivateCloneURL(cloneURL, credentials)
	}
	branchName := extractBranchName(repository.BranchRef)
	cloneCmd := fmt.Sprintf("git clone -b %s --single-branch %s", branchName, cloneURL)

	cloneOut, err := g.executor.Execute(cloneCmd, workDir)
	if err != nil {
		log.Println(cloneOut.StdErr)
		return "", errors.Wrap(err, "failed to clone repository")
	}
	return fmt.Sprintf("%s/%s", workDir, repository.Name), nil
}

func createCredentials(repository drizzle.Repository, secretsManager drizzle.SecretRetriever) (credentials, error) {
	secret := fmt.Sprintf("%s_%s_%s", gitHubSecretsPrefix, repositorySecretsPrefix, repository.ID)
	secretValue, err := secretsManager.RetrieveSecret(secret)
	if err != nil {
		return credentials{}, errors.Wrapf(err, "failed to retrieve value for secret %q", secret)
	}
	return credentials{
		User:     defaultGitUser,
		Password: secretValue,
	}, nil
}

// https://username:password@github.com/username/repository.git
func createPrivateCloneURL(cloneURL string, cred credentials) string {
	protocolSep := strings.Index(cloneURL, "://") + 3
	return fmt.Sprintf("%s%s:%s@%s", cloneURL[:protocolSep], cred.User, cred.Password, cloneURL[protocolSep:])
}

func extractBranchName(branchRef string) string {
	return branchRef[strings.LastIndex(branchRef, "/")+1:]
}
