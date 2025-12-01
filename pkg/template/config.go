// ABOUTME: Configuration file loader for .bingo.yaml metadata
// ABOUTME: Reads service mappings from bingo project template
package template

import (
	"os"

	"gopkg.in/yaml.v3"
)

// BingoConfig represents .bingo.yaml configuration file structure
type BingoConfig struct {
	Version  int                    `yaml:"version"`
	Services map[string]ServiceInfo `yaml:"services"`
}

// ServiceInfo describes a service's directory structure
type ServiceInfo struct {
	Cmd         string `yaml:"cmd"`
	Internal    string `yaml:"internal"`
	Description string `yaml:"description"`
}

// LoadBingoConfig loads and parses .bingo.yaml
func LoadBingoConfig(path string) (*BingoConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config BingoConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
