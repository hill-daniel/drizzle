package shell_test

import (
	"github.com/hill-daniel/drizzle/shell"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	executor := shell.NewExecutor("")

	out, err := executor.Execute("ls -lah", "/tmp")

	assert.NoError(t, err)
	assert.NotNil(t, out)
}
