// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/internal/animation"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Cell is a physical square on the board.
type Cell struct {
	frame    animation.FrameLoc
	rotation rl.Vector2
}

// Board is the game board, which is a grid of cells.
type Board struct {
	Rows, Columns     int
	backgroundSprites *animation.SpriteSheet
	cells             [][]Cell
}

const (
	CellSize = 24
	Scale    = 3.0
)

var (
	BlackSquare = animation.FrameLoc{Row: 1, Col: 2}
	WhiteSquare = animation.FrameLoc{Row: 1, Col: 8}
)

// newBoard creates a new, empty board.
func newBoard() Board {
	b := Board{
		Rows:              10,
		Columns:           10,
		cells:             make([][]Cell, 10),
		backgroundSprites: animation.LoadSpriteSheet("dungeon_tiles.png", 4, 9),
	}
	for i := range b.cells {
		b.cells[i] = make([]Cell, 10)
	}

	// Initialize the board's cells with frames
	for i := range b.Rows {
		for j := range b.Columns {
			var f = WhiteSquare
			if (i+j)%2 == 0 {
				f = BlackSquare
			}
			b.cells[i][j] = Cell{
				frame:    f,
				rotation: rl.Vector2{X: 1, Y: 0},
			}
		}
	}

	return b
}

// Render draws the board to the screen at the given location.
func (b *Board) Render(loc rl.Vector2) error {
	// First draw the board itself
	for i := range b.Rows {
		for j := range b.Columns {
			if err := b.backgroundSprites.DrawFrame(
				b.cells[i][j].frame,
				rl.Vector2{X: loc.X + float32(j*CellSize*Scale), Y: loc.Y + float32(i*CellSize*Scale)},
				Scale,
				b.cells[i][j].rotation); err != nil {
				return fmt.Errorf("failed to draw cell: %w", err)
			}
		}
	}
	return nil
}
