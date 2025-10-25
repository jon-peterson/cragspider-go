// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"fmt"
)

// Game represents a single instance of a game.
type Game struct {
	Board  *Board
	config *GameConfig
}

// NewGame returns a new game with a board ready to be played.
func NewGame() (*Game, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
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
