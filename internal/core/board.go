// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/internal/animation"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Cell is a physical square on the board.
type Cell struct {
	frame    animation.FrameCoords
	rotation rl.Vector2
}

// Position is a [row,col] on the board.
type Position [2]int

// Move represents a [deltaRow, deltaCol] move on the board.
type Move [2]int

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
var CardinalDirections = [4]Move{
	{0, -1}, // up
	{1, 0},  // right
	{0, 1},  // down
	{-1, 0}, // left
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

	// Initialize the board's cells with surfaces. The surfaces are colored in pairs but are of
	// random frame and orientation. This gives the board variety over plays.
	for i := range b.Rows {
		for j := range b.Columns {
			var f animation.FrameCoords
			if (i+j)%2 == 0 {
				f = animation.FrameCoords{
					IntInRange(0, 1),
					IntInRange(0, 2),
				}
			} else {
				f = animation.FrameCoords{
					IntInRange(0, 1),
					IntInRange(6, 8),
				}
			}
			facing := Choice(CardinalDirections[:])
			b.cells[i][j] = Cell{
				frame:    f,
				rotation: rl.Vector2{X: float32(facing[0]), Y: float32(facing[1])},
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
