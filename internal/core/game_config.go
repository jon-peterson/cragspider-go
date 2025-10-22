package core

import (
	_ "embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

// Color represents the color of a piece or player
type Color string

const (
	// White represents the white player/pieces
	White Color = "white"
	// Black represents the black player/pieces
	Black Color = "black"
)

type Position [2]int // [row, col]
type Move [2]int     // [deltaRow, deltaCol]
type SpriteCoords []Position

type PieceConfig struct {
	Name    string                  `yaml:"name"`
	Sprites map[string]SpriteCoords `yaml:"sprites"`
	Moves   []Move                  `yaml:"moves"`
}

// BoardConfig is the starting status of the board
type BoardConfig struct {
	White [][]interface{} `yaml:"white"`
	Black [][]interface{} `yaml:"black"`
}

// GetStartingPositions returns the starting positions for the specified color.
func (b *BoardConfig) GetStartingPositions(color Color) ([][]interface{}, error) {
	switch color {
	case White:
		return b.White, nil
	case Black:
		return b.Black, nil
	default:
		return nil, fmt.Errorf("invalid color: %s, must be 'white' or 'black'", color)
	}
}

// GameConfig is the full game configuration
type GameConfig struct {
	Pieces []PieceConfig `yaml:"pieces"`
	Board  BoardConfig   `yaml:"board"`
}

var (
	config     *GameConfig
	configOnce sync.Once
)

//go:embed game_config.yaml
var configData []byte

// LoadConfig loads the game configuration from the embedded YAML file.
// It uses sync.Once to ensure the configuration is only loaded once.
func LoadConfig() (*GameConfig, error) {
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

// GetConfig returns the loaded configuration.
// It returns an error if the configuration has not been loaded yet.
func GetConfig() (*GameConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration not loaded, call LoadConfig first")
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
