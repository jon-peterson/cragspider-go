// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	board Board
}

// NewGame returns a new game with a board ready to be played.
func NewGame() Game {
	return Game{
		board: newBoard(),
	}
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
