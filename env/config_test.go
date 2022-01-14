package env_test

import (
	"github.com/hill-daniel/drizzle/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig_should_parse_basic_config(t *testing.T) {
	parser := &env.ConfigParser{}

	config, err := parser.Parse("../testdata/yml")

	assert.NoError(t, err)
	assert.Equal(t, 3, len(config.Stages))
	assert.Equal(t, 1, len(config.Variables))
}
