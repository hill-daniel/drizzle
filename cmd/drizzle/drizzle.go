package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hill-daniel/drizzle/aws"
	"github.com/hill-daniel/drizzle/core"
	"github.com/hill-daniel/drizzle/env"
	"github.com/hill-daniel/drizzle/shell"
	"github.com/pkg/errors"
	"log"
	"os"
)

const (
	envWorkDir = "WORK_DIR"
)

func main() {
	executor := shell.NewExecutor("")
	awsSession, err := session.NewSession()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create aws session"))
	}
	secretsManager := aws.NewSecretManager(secretsmanager.New(awsSession))
	git := shell.NewGit(executor, secretsManager)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize runner"))
	}
	configParser := &env.ConfigParser{}
	job := core.NewJob(git, configParser, executor)
	workDir := os.Getenv(envWorkDir)

	handler := aws.NewLambdaHandler(job, workDir)
	lambda.Start(handler.Handle)
}
