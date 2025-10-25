// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game represents a single instance of a game.
type Game struct {
	board  *Board
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
		board:  b,
		config: cfg,
	}, nil
}

// Over returns true if the game is over.
func (g *Game) Over() bool {
	return false
}

// Render draws the current game state to the screen.
func (g *Game) Render() error {
	if err := g.board.Render(rl.Vector2{X: 96, Y: 48}); err != nil {
		return fmt.Errorf("failed to render board: %w", err)
	}
	return nil
}
