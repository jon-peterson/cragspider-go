// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/pkg/animation"
	_ "embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

// SpriteCoords is a location [row,col] in a spritesheet.
type SpriteCoords []animation.FrameCoords

// PieceConfig represents a type of piece in the game, like bishop or pawn.
type PieceConfig struct {
	Name    string                 `yaml:"name"`
	Sprites map[Color]SpriteCoords `yaml:"sprites"`
	Moves   []Move                 `yaml:"moves"`
}

// BoardPosition represents a starting position on the board: what piece and where.
type BoardPosition struct {
	Name     string   `yaml:"name"`
	Position Position `yaml:"position"`
}

// BoardConfig is how the board looks at the very start of the game.
type BoardConfig struct {
	Rows    int             `yaml:"rows"`
	Columns int             `yaml:"columns"`
	White   []BoardPosition `yaml:"white"`
	Black   []BoardPosition `yaml:"black"`
}

// GetStartingPositions returns the starting positions for the specified color.
func (b *BoardConfig) GetStartingPositions(color Color) ([]BoardPosition, error) {
	switch color {
	case White:
		return b.White, nil
	case Black:
		return b.Black, nil
	default:
		return nil, fmt.Errorf("invalid color: %s, must be 'white' or 'black'", color)
	}
}

// GameConfig holds all the parameters for how the game is played.
type GameConfig struct {
	Pieces []PieceConfig `yaml:"pieces"`
	Board  BoardConfig   `yaml:"board"`
}

var (
	config     *GameConfig
	configOnce sync.Once
)

//go:embed game_config.yml
var configData []byte

// GetConfig loads and returns the game configuration from the embedded YAML file.
// It uses sync.Once to ensure the configuration is only loaded once.
func GetConfig() (*GameConfig, error) {
	var loadErr error

	configOnce.Do(func() {
		var cfg GameConfig
		if err := yaml.Unmarshal(configData, &cfg); err != nil {
			loadErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}

		config = &cfg
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return config, nil
}

// GetPieceConfig returns the configuration for the specified piece name.
// It returns an error if the piece is not found.
func (g *GameConfig) GetPieceConfig(name string) (*PieceConfig, error) {
	for _, piece := range g.Pieces {
		if piece.Name == name {
			return &piece, nil
		}
	}
	return nil, fmt.Errorf("piece '%s' not found in configuration", name)
}
