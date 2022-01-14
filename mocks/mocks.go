package mocks

import (
	"github.com/hill-daniel/drizzle"
)

//ConfigParser mock.
type ConfigParser struct {
	ParseFunc func(path string) (drizzle.Pipeline, error)
}

// Parse mock.
func (p *ConfigParser) Parse(path string) (drizzle.Pipeline, error) {
	return p.ParseFunc(path)
}

// Git mock.
type Git struct {
	CloneFunc func(repository drizzle.Repository, workDir string) (string, error)
}

// Clone mock.
func (g *Git) Clone(repository drizzle.Repository, workDir string) (string, error) {
	return g.CloneFunc(repository, workDir)
}

//Executor mock.
type Executor struct {
	ExecuteFunc func(command string, workDir string) (drizzle.ShellOut, error)
}

// Execute mock.
func (e *Executor) Execute(command string, workDir string) (drizzle.ShellOut, error) {
	return e.ExecuteFunc(command, workDir)
}

// SecretsManager mock.
type SecretsManager struct {
	GetSecretFunc func(secretID string) (string, error)
}

// RetrieveSecret mock.
func (sm *SecretsManager) RetrieveSecret(secretID string) (string, error) {
	return sm.GetSecretFunc(secretID)
}
