// Copyright 2025 Ideograph LLC. All rights reserved.

package ai

import (
	_ "embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

// AIPlayerConfig represents the configuration for an AI player.
type AIPlayerConfig struct {
	Name    string             `yaml:"name"`
	Scoring map[string]float32 `yaml:"scoring"`
}

// AIConfig holds all AI player configurations.
type AIConfig struct {
	Players []AIPlayerConfig `yaml:"players"`
}

var (
	aiConfig     *AIConfig
	aiConfigOnce sync.Once
)

//go:embed ai_config.yml
var aiConfigData []byte

// GetAIConfig loads and returns the AI configuration from the embedded YAML file.
// It uses sync.Once to ensure the configuration is only loaded once.
func GetAIConfig() (*AIConfig, error) {
	var loadErr error

	aiConfigOnce.Do(func() {
		var cfg AIConfig
		if err := yaml.Unmarshal(aiConfigData, &cfg); err != nil {
			loadErr = fmt.Errorf("failed to unmarshal AI config: %w", err)
			return
		}

		aiConfig = &cfg
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return aiConfig, nil
}

// GetPlayerConfig returns the configuration for the specified player name.
// It returns an error if the player is not found.
func (a *AIConfig) GetPlayerConfig(name string) (*AIPlayerConfig, error) {
	for _, player := range a.Players {
		if player.Name == name {
			return &player, nil
		}
	}
	return nil, fmt.Errorf("AI player '%s' not found in configuration", name)
}
