// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import rl "github.com/gen2brain/raylib-go/raylib"

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
func (g *Game) Render() {
	g.board.Draw(rl.Vector2{X: 24, Y: 24})
}
