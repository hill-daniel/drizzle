package env

import (
	"fmt"
	"github.com/hill-daniel/drizzle"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

const (
	// PipelineConfigFilename is the pipeline config the parser is looking for.
	PipelineConfigFilename = ".drizzle.yml"
)

//ConfigParser parses given file from path to a Pipeline.
type ConfigParser struct {
	secretRetriever drizzle.SecretRetriever
}

// Parse parses yml file of given path and creates a Pipeline.
func (cp *ConfigParser) Parse(configDirPath string) (drizzle.Pipeline, error) {
	pipeline := &drizzle.Pipeline{}
	configPath := fmt.Sprintf("%s/%s", configDirPath, PipelineConfigFilename)
	err := cleanenv.ReadConfig(configPath, pipeline)
	if err != nil {
		return *pipeline, errors.Wrapf(err, "failed to parse config")
	}
	resolvedKeyValues, err := ExpandVars(pipeline.Variables, cp.secretRetriever)
	if err != nil {
		return *pipeline, errors.Wrapf(err, "failed to resolve variables")
	}
	pipeline.Variables = resolvedKeyValues
	pipeline.Path = configDirPath
	return *pipeline, nil
}
