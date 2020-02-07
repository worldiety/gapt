package gapt

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	Include []string `yaml:"include",envconfig:"GAPT_INCLUDE"`
	Exclude []string `yaml:"exclude",envconfig:"GAPT_EXCLUDE"`
}

func (c *Config) Default() {
	c.Include = []string{"*"}
	c.Exclude = []string{"makefile", "*.md", "*.go", "*.mod", "license", "*.sum"}
}

// ParseFile loads from a local yaml file and puts environment GAPT_* values on top.
func (c *Config) ParseFile(ymlFile string) (err error) {
	f, err := os.Open(ymlFile)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", ymlFile, err)
	}
	defer func() { err = f.Close() }()
	return c.Parse(f)
}

// Parse loads from a yaml reader and puts environment GAPT_* values on top.
func (c *Config) Parse(reader io.Reader) error {
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("failed to decode yml: %w", err)
	}

	// mix in the env
	err = envconfig.Process("", c)
	if err != nil {
		return fmt.Errorf("failed to decode env: %w", err)
	}

	return nil
}

// Write serializes the config into a yaml writer.
func (c *Config) Write(writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	err := encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("failed to encode yaml: %w", err)
	}
	return nil
}

// WriteFile serializes the config into a yaml file.
func (c *Config) WriteFile(ymlFile string) error {
	f, err := os.Create(ymlFile)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", ymlFile, err)
	}
	defer func() { err = f.Close() }()
	return c.Write(f)
}
