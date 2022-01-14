package core

import (
	"context"
	"github.com/hill-daniel/drizzle"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

// Job executes preparing steps before the runner.
type Job struct {
	git          drizzle.Cloner
	configParser drizzle.Parser
	executor     drizzle.Executor
}

// NewJob creates a new job.
func NewJob(git drizzle.Cloner, configParser drizzle.Parser, executor drizzle.Executor) *Job {
	return &Job{
		git:          git,
		configParser: configParser,
		executor:     executor,
	}
}

// Build executes the build steps for pipeline execution.
func (j *Job) Build(ctx context.Context, repository drizzle.Repository, workDir string) error {
	clonePath, err := j.git.Clone(repository, workDir)
	if err != nil {
		return err
	}

	pipeline, err := j.configParser.Parse(clonePath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse config in %q", clonePath)
	}

	if err = Run(pipeline, repository.BranchRef, j.executor); err != nil {
		return errors.Wrapf(err, "failed to run pipeline")
	}
	return nil
}

// Run executes a build pipeline from a GitHub web hook request.
func Run(pipeline drizzle.Pipeline, branchRef string, executor drizzle.Executor) error {
	for _, stage := range pipeline.Stages {
		if !shouldExecuteOnBranch(stage, branchRef) {
			continue
		}
		log.Printf("Executing stage %q ", stage.Name)

		for _, command := range stage.Execute {
			cmd := os.Expand(command, func(s string) string {
				return pipeline.Variables[s]
			})
			out, err := executor.Execute(cmd, pipeline.Path)

			if stage.Debug {
				log.Println(command)
				log.Println(out.StdOut)
				log.Println(out.StdErr)
			}
			if err != nil {
				return errors.Wrapf(err, "failed to execute command %q of stage %q", command, stage.Name)
			}
		}
	}
	return nil
}

func shouldExecuteOnBranch(stage *drizzle.Stage, branchRef string) bool {
	if len(stage.Branches) == 0 {
		return true
	}
	for _, branch := range stage.Branches {
		if strings.HasSuffix(branchRef, branch) {
			return true
		}
	}
	return false
}
