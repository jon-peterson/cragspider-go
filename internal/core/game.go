// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game represents a single instance of a game.
type Game struct {
	Board  *Board
	config *GameConfig
}

// NewGame returns a new game with the standard configuration.
func NewGame() (*Game, error) {
	cfg, err := GetConfig()
	if err != nil {
		rl.TraceLog(rl.LogFatal, "error loading default configuration: %v", err)
	}
	return NewGameWithConfig(cfg)
}

// NewGameWithConfig returns a new game with a board ready to be played. It accepts a GameConfig parameter to allow
// custom configurations.
func NewGameWithConfig(cfg *GameConfig) (*Game, error) {
	if cfg == nil {
		return nil, fmt.Errorf("game configuration cannot be nil")
	}
	b, err := newBoard(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	return &Game{
		Board:  b,
		config: cfg,
	}, nil
}

// Over returns true if the game is over.
func (g *Game) Over() bool {
	return false
}
