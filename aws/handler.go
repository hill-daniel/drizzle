package aws

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/hill-daniel/drizzle"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
)

const (
	tmpDirPattern = "_drizzle"
)

// LambdaHandler implements an AWS Lambda handler for an incoming SQSEvent Message.
type LambdaHandler struct {
	job     drizzle.Builder
	workDir string
}

// NewLambdaHandler creates a new LambdaHandler.
func NewLambdaHandler(job drizzle.Builder, workDir string) *LambdaHandler {
	return &LambdaHandler{
		job:     job,
		workDir: workDir,
	}
}

// Handle handles incoming SQSEvents.
func (lh *LambdaHandler) Handle(ctx context.Context, sqsEvent *events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		repository := &drizzle.Repository{}
		if err := json.Unmarshal([]byte(message.Body), repository); err != nil {
			log.Println(errors.Wrap(err, "failed to unmarshal json from SQS event body"))
			return nil
		}
		if err := lh.runPipeline(ctx, *repository); err != nil {
			log.Println(errors.Wrapf(err, "failed to run pipeline for repository %q with id %s", repository.FullName, repository.ID))
			return nil
		}
		// Do not return error, we want the sqs message to be processed, since to fix an error
		// we have to commit stuff again anyway, which will trigger the pipeline again.
	}
	return nil
}

func (lh *LambdaHandler) runPipeline(ctx context.Context, repository drizzle.Repository) error {
	pipelineDir, err := ioutil.TempDir(lh.workDir, tmpDirPattern)
	if err != nil {
		return errors.Wrapf(err, "failed to create temporary directory for pipeline")
	}
	defer func() {
		err := os.RemoveAll(pipelineDir)
		if err != nil {
			log.Println(errors.Wrap(err, "failed to cleanup temporary pipeline directory"))
		}
	}()

	if err = lh.job.Build(ctx, repository, pipelineDir); err != nil {
		return errors.Wrapf(err, "failed to execute job")
	}
	log.Println("Pipeline executed!")
	return nil
}
