package drizzle

// Pipeline represents a build pipeline.
type Pipeline struct {
	Variables map[string]string `yaml:"variables"`
	Stages    Stages            `yaml:"stages"`
	Path      string
}

// Stages represents the stages of the build.
type Stages []*Stage

// Stage represents the stage of the build.
type Stage struct {
	Name     string   `yaml:"name"`
	Execute  []string `yaml:"execute"`
	Branches []string `yaml:"branches"`
	Debug    bool     `yaml:"debug"`
}

// Parser parses input and creates a Pipeline
type Parser interface {
	Parse(path string) (Pipeline, error)
}
