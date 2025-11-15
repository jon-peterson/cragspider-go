// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/internal/animation"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/samber/lo"
)

// Square is a physical square on the board.
type Square struct {
	Frame    animation.FrameCoords
	Rotation rl.Vector2
}

// Position is a [row,col] on the board.
type Position [2]int

// String returns a nicely formatted string representation of the position.
func (pos Position) String() string {
	return fmt.Sprintf("[%d,%d]", pos[0], pos[1])
}

// Move represents a [deltaRow, deltaCol] move on the board.
type Move [2]int

// String returns a nicely formatted string representation of a move.
func (m Move) String() string {
	return fmt.Sprintf("[%d,%d]", m[0], m[1])
}

// Board is the game board, which is a grid of squares upon which there are pieces.
type Board struct {
	Rows, Columns int
	squares       [][]Square
	pieces        [][]*Piece

	config *GameConfig
}

const (
	// SpriteSize is the standard size of one of our square sprites
	SpriteSize = 24
	// Scale is the scale factor for the board, to make each square bigger than 24x24.
	Scale float32 = 3.0
	// SquareSize is the standard size of a square in number of pixels.
	SquareSize int = SpriteSize * int(Scale)
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
		Rows:    10,
		Columns: 10,
		squares: make([][]Square, 10),
		pieces:  make([][]*Piece, 10),
		config:  config,
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
				Frame:    f,
				Rotation: rl.Vector2{X: float32(facing[0]), Y: float32(facing[1])},
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
				Name:   pos.Name,
				Color:  color,
				Config: *pieceConfig,
			}

			if err := b.PlacePiece(piece, pos.Position); err != nil {
				return fmt.Errorf("failed to place %s %s at %v: %w", color, pos.Name, pos.Position, err)
			}
		}
	}
	return nil
}

// IsOccupied returns true if there's a piece at the specified location.
func (b *Board) IsOccupied(pos Position) bool {
	return b.GetPieceAt(pos) != nil
}

// IsValid returns true if the specified position is on the board.
func (b *Board) IsValid(pos Position) bool {
	return pos[0] >= 0 && pos[0] < b.Rows && pos[1] >= 0 && pos[1] < b.Columns
}

// PieceLocation returns the position of the specified piece on the board, assuming that it can be
// found. piece should not be nil; if it is, an error is returned.
func (b *Board) PieceLocation(piece *Piece) (Position, error) {
	if piece == nil {
		return Position{}, fmt.Errorf("cannot find location of nil piece")
	}
	for i := range b.Rows {
		for j := range b.Columns {
			if b.pieces[i][j] == piece {
				return Position{i, j}, nil
			}
		}
	}
	return Position{}, fmt.Errorf("%s not found on board", piece)
}

// PlacePiece puts the specified piece in the specified location, returning an error if the position is occupied.
func (b *Board) PlacePiece(piece *Piece, pos Position) error {
	if piece == nil {
		return fmt.Errorf("piece is nil")
	}
	if !b.IsValid(pos) {
		return fmt.Errorf("%s is out of bounds", pos)
	}
	if b.IsOccupied(pos) {
		return fmt.Errorf("%s is occupied", pos)
	}
	b.pieces[pos[0]][pos[1]] = piece
	return nil
}

// MovePiece moves the existing piece from the specified position. An error is returned if it isn't
// one of the valid moves for the piece.
func (b *Board) MovePiece(piece *Piece, start Position, move Move) error {
	// Make sure the piece is actually at that starting position.
	if b.pieces[start[0]][start[1]] != piece {
		return fmt.Errorf("%s is not at %s", piece, start)
	}
	// Make sure that the move being passed in is valid for this piece from the starting position.
	validMoves := piece.ValidMoves(start, b)
	end := Position{
		start[0] + move[0],
		start[1] + move[1],
	}
	if !lo.Contains(validMoves, end) {
		return fmt.Errorf("move %v is not valid for piece %s", move, piece)
	}
	// Actually put the piece in its new location.
	err := b.PlacePiece(piece, end)
	if err != nil {
		return fmt.Errorf("cannot move %s to %s: %w", piece, end, err)
	}
	// Clear the original starting position so a piece can move back.
	b.pieces[start[0]][start[1]] = nil
	return nil
}

// GetSquareAt returns the square at the given position.
func (b *Board) GetSquareAt(pos Position) *Square {
	return &b.squares[pos[0]][pos[1]]
}

// GetPieceAt returns the piece at the given position, or nil if empty.
func (b *Board) GetPieceAt(pos Position) *Piece {
	return b.pieces[pos[0]][pos[1]]
}
