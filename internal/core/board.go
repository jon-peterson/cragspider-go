// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/internal/animation"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Square is a physical square on the board.
type Square struct {
	frame    animation.FrameCoords
	rotation rl.Vector2
}

// Position is a [row,col] on the board.
type Position [2]int

// Move represents a [deltaRow, deltaCol] move on the board.
type Move [2]int

// Board is the game board, which is a grid of squares upon which there are pieces.
type Board struct {
	Rows, Columns int
	squares       [][]Square
	pieces        [][]*Piece

	backgroundSprites *animation.SpriteSheet
	whiteSprites      *animation.SpriteSheet
	blackSprites      *animation.SpriteSheet

	config *GameConfig // Reference to the game configuration
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

// newBoard creates a new board with the given game configuration.
// If config is nil, it will use the default configuration.
func newBoard(config *GameConfig) (*Board, error) {
	b := &Board{
		Rows:              10,
		Columns:           10,
		squares:           make([][]Square, 10),
		pieces:            make([][]*Piece, 10),
		backgroundSprites: animation.LoadSpriteSheet("dungeon_tiles.png", 4, 9),
		whiteSprites:      animation.LoadSpriteSheet("adventurer_pieces.png", 6, 18),
		blackSprites:      animation.LoadSpriteSheet("monster_pieces.png", 11, 18),
		config:            config,
	}
	b.initializeSquares()
	err := b.placeStartingPieces()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// initializeSquares initializes the board's squares with surfaces. The surfaces are colored in pairs but are of
// random frame and orientation. This gives the board variety over plays.
func (b *Board) initializeSquares() {
	for i := 0; i < b.Rows; i++ {
		b.squares[i] = make([]Square, b.Columns)
		b.pieces[i] = make([]*Piece, b.Columns)
		for j := 0; j < b.Columns; j++ {
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
			b.squares[i][j] = Square{
				frame:    f,
				rotation: rl.Vector2{X: float32(facing[0]), Y: float32(facing[1])},
			}
		}
	}
}

// placeStartingPieces places all pieces in their starting positions according to the game config.
func (b *Board) placeStartingPieces() error {
	for _, color := range []Color{White, Black} {
		positions, err := b.config.Board.GetStartingPositions(color)
		if err != nil {
			return fmt.Errorf("failed to get starting positions for %s: %w", color, err)
		}
		for _, pos := range positions {
			pieceConfig, err := b.config.GetPieceConfig(pos.Name)
			if err != nil {
				return fmt.Errorf("failed to get config for piece %s: %w", pos.Name, err)
			}

			piece := &Piece{
				name:   pos.Name,
				color:  color,
				config: *pieceConfig,
			}

			if err := b.PlacePiece(*piece, pos.Position); err != nil {
				return fmt.Errorf("failed to place %s %s at %v: %w", color, pos.Name, pos.Position, err)
			}
		}
	}
	return nil
}

// IsOccupied returns true if there's a piece at the specified location.
func (b *Board) IsOccupied(pos Position) bool {
	return b.pieces[pos[0]][pos[1]] != nil
}

// IsValid returns true if the specified position is on the board.
func (b *Board) IsValid(pos Position) bool {
	return pos[0] >= 0 && pos[0] < b.Rows && pos[1] >= 0 && pos[1] < b.Columns
}

// PlacePiece puts the specified piece in the specified location, returning an error if the position is occupied.
// The piece is copied when placed on the board to prevent accidental modifications.
func (b *Board) PlacePiece(piece Piece, pos Position) error {
	if !b.IsValid(pos) {
		return fmt.Errorf("position %v is out of bounds", pos)
	}
	if b.IsOccupied(pos) {
		return fmt.Errorf("position %v is occupied", pos)
	}
	b.pieces[pos[0]][pos[1]] = &piece
	return nil
}

// Render draws the board to the screen at the given location.
func (b *Board) Render(loc rl.Vector2) error {
	// First draw the board itself
	for i := range b.Rows {
		for j := range b.Columns {
			err := b.backgroundSprites.DrawFrameRotated(
				b.squares[i][j].frame,
				rl.Vector2{X: loc.X + float32(j*CellSize*Scale), Y: loc.Y + float32(i*CellSize*Scale)},
				Scale,
				b.squares[i][j].rotation)
			if err != nil {
				return fmt.Errorf("failed to draw cell: %w", err)
			}
		}
	}
	// Now draw each of the pieces on the board
	for i := range b.Rows {
		for j := range b.Columns {
			piece := b.pieces[i][j]
			if piece != nil {
				sheet := b.whiteSprites
				if piece.color == Black {
					sheet = b.blackSprites
				}
				err := sheet.DrawFrame(
					piece.config.Sprites[piece.color][0],
					rl.Vector2{X: loc.X + float32(j*CellSize*Scale), Y: loc.Y + float32(i*CellSize*Scale)},
					Scale)
				if err != nil {
					return fmt.Errorf("failed to draw piece: %w", err)
				}
			}
		}
	}

	return nil
}
