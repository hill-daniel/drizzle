package drizzle

import "context"

// Repository represents a repository for build pipeline execution.
type Repository struct {
	ID        string `json:"id"`
	BranchRef string `json:"branch_ref"`
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Private   bool   `json:"private"`
	URL       string `json:"url"`
	CloneURL  string `json:"clone_url"`
}

// Cloner represents an interface to clone repositories.
type Cloner interface {
	Clone(repository Repository, workDir string) (string, error)
}

// Builder executes preparing steps before running a pipeline.
type Builder interface {
	Build(ctx context.Context, repository Repository, workDir string) error
}
