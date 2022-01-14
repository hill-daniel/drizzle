package shell

import (
	"bytes"
	"github.com/hill-daniel/drizzle"
	"os/exec"
)

const (
	standardShell = "/bin/sh"
)

// Executor executes commands on a shell.
type Executor struct {
	shellToUse string
}

// NewExecutor creates a new Executor. Shell string may be empty, in which case the default "/bin/sh" is used.
func NewExecutor(shell string) *Executor {
	if len(shell) == 0 {
		return &Executor{shellToUse: standardShell}
	}
	return &Executor{shellToUse: shell}
}

// Execute executes given command and returns stdout and stderr of command, if available.
func (e *Executor) Execute(command string, workDir string) (drizzle.ShellOut, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := &exec.Cmd{
		Path:   e.shellToUse,
		Args:   append([]string{e.shellToUse}, "-c", command),
		Dir:    workDir,
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	err := cmd.Run()

	return drizzle.ShellOut{
		StdErr: stderr.String(),
		StdOut: stdout.String(),
	}, err
}
