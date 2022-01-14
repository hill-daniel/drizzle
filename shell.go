package drizzle

// ShellOut represents the standard output and the standard error output of a shell command execution.
type ShellOut struct {
	StdErr string
	StdOut string
}

// Executor executes commands.
type Executor interface {
	Execute(command string, workDir string) (ShellOut, error)
}
