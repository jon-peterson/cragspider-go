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

// CardinalDirections contains unit vectors for the four cardinal directions: up, right, down, left.
var CardinalDirections = [4]rl.Vector2{
	{X: 0, Y: -1}, // up
	{X: 1, Y: 0},  // right
	{X: 0, Y: 1},  // down
	{X: -1, Y: 0}, // left
}

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
			var f animation.FrameLoc
			if (i+j)%2 == 0 {
				f = animation.FrameLoc{
					Row: IntInRange(0, 1),
					Col: IntInRange(0, 2),
				}
			} else {
				f = animation.FrameLoc{
					Row: IntInRange(0, 1),
					Col: IntInRange(6, 8),
				}
			}
			b.cells[i][j] = Cell{
				frame:    f,
				rotation: Choice(CardinalDirections[:]),
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
