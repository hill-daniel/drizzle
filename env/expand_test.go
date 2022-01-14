package env_test

import (
	"github.com/hill-daniel/drizzle/env"
	"github.com/hill-daniel/drizzle/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_should_return_input_if_nothing_to_expand(t *testing.T) {
	input := map[string]string{"key": "value"}

	expanded, err := env.ExpandVars(input, &mocks.SecretsManager{})

	assert.NoError(t, err)
	assert.Equal(t, input, expanded)
}

func Test_should_expand_input_from_environment_variables(t *testing.T) {
	input := map[string]string{"path": "$PATH"}
	err := os.Setenv("PATH", "/tmp/somePath")
	assert.NoError(t, err)

	expanded, err := env.ExpandVars(input, &mocks.SecretsManager{})

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"path": "/tmp/somePath"}, expanded)
}

func Test_should_expand_input_from_secrets_manager(t *testing.T) {
	input := map[string]string{"secret": "$SECRET"}
	secretsManager := &mocks.SecretsManager{GetSecretFunc: func(secretID string) (string, error) {
		assert.Equal(t, secretID, "SECRET")
		return "MySecretValue1234", nil
	}}

	expanded, err := env.ExpandVars(input, secretsManager)

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"secret": "MySecretValue1234"}, expanded)
}
