// ABOUTME: Configuration file loader for .bingoctl.yaml metadata
// ABOUTME: Reads service mappings from bingo project template
package template

import (
	"os"

	"gopkg.in/yaml.v3"
)

// BingoctlConfig represents .bingoctl.yaml configuration file structure
type BingoctlConfig struct {
	Version  int                    `yaml:"version"`
	Services map[string]ServiceInfo `yaml:"services"`
}

// ServiceInfo describes a service's directory structure
type ServiceInfo struct {
	Cmd         string `yaml:"cmd"`
	Internal    string `yaml:"internal"`
	Description string `yaml:"description"`
}

// loadBingoctlConfig loads and parses .bingoctl.yaml
func loadBingoctlConfig(path string) (*BingoctlConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config BingoctlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
