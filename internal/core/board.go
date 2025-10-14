// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/internal/animation"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Cell is a square on the board.
type Cell struct {
	X, Y int
}

// Board is the game board, which is a grid of cells.
type Board struct {
	Rows, Columns     int
	backgroundSprites *animation.SpriteSheet
}

const (
	CellSize = 24
)

var (
	BlackSquare = animation.FrameLoc{Row: 1, Col: 2}
	WhiteSquare = animation.FrameLoc{Row: 1, Col: 8}
)

// newBoard creates a new, empty board.
func newBoard() Board {
	return Board{
		Rows:              10,
		Columns:           10,
		backgroundSprites: animation.LoadSpriteSheet("dungeon_tiles.png", 4, 9),
	}
}

// Render draws the board to the screen at the given location.
func (b *Board) Render(loc rl.Vector2) error {
	// First draw the board itself
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Columns; j++ {
			var frame = WhiteSquare
			if (i+j)%2 == 0 {
				frame = BlackSquare
			}
			if err := b.backgroundSprites.DrawFrame(
				frame,
				rl.Vector2{X: loc.X + float32(j*CellSize), Y: loc.Y + float32(i*CellSize)},
				rl.Vector2{X: 1, Y: 0}); err != nil {
				return fmt.Errorf("failed to draw cell: %w", err)
			}
		}
	}
	return nil
}
