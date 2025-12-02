// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"cragspider-go/pkg/graphics"
	"cragspider-go/pkg/random"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/samber/lo"
)

// Square is a physical square on the board.
type Square struct {
	Frame    graphics.FrameCoords
	Rotation rl.Vector2
}

// SquareGrid holds the immutable grid of squares on a board.
// This is shared across all copies of a board since squares never change.
type SquareGrid struct {
	data [][]Square
}

// newSquareGrid creates a new SquareGrid with the given dimensions.
func newSquareGrid(rows, cols int) *SquareGrid {
	data := make([][]Square, rows)
	for i := range data {
		data[i] = make([]Square, cols)
	}
	return &SquareGrid{data: data}
}

// Move represents a [deltaRow, deltaCol] move on the board.
type Move [2]int

// String returns a nicely formatted string representation of a move.
func (m Move) String() string {
	return fmt.Sprintf("[%d,%d]", m[0], m[1])
}

type Position [2]int

// String returns a nicely formatted string representation of the position.
func (pos Position) String() string {
	return fmt.Sprintf("[%d,%d]", pos[0], pos[1])
}

// Add returns the position after moving from the given position by the given move.
func (pos Position) Add(move Move) Position {
	return Position{
		pos[0] + move[0],
		pos[1] + move[1],
	}
}

// Board is the game board, which is a grid of squares upon which there are pieces.
type Board struct {
	Rows, Columns int
	squares       *SquareGrid
	pieces        [][]*Piece
	captured      map[Color][]*Piece // Pieces captured by each color

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
	rows := config.Board.Rows
	columns := config.Board.Columns

	b := &Board{
		Rows:     rows,
		Columns:  columns,
		squares:  newSquareGrid(rows, columns),
		pieces:   make([][]*Piece, rows),
		captured: make(map[Color][]*Piece),
		config:   config,
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
		b.pieces[i] = make([]*Piece, b.Columns)
		for j := 0; j < b.Columns; j++ {
			var f graphics.FrameCoords
			if (i+j)%2 == 0 {
				f = graphics.FrameCoords{
					random.IntInRange(0, 1),
					random.IntInRange(0, 2),
				}
			} else {
				f = graphics.FrameCoords{
					random.IntInRange(0, 1),
					random.IntInRange(6, 8),
				}
			}
			facing := random.Choice(CardinalDirections[:])
			b.squares.data[i][j] = Square{
				Frame:    f,
				Rotation: rl.Vector2{X: float32(facing[0]), Y: float32(facing[1])},
			}
		}
	}
}

// placeStartingPieces places all pieces in their starting positions according to the game config.
func (b *Board) placeStartingPieces() error {
	currentBoard := b
	for _, color := range []Color{White, Black} {
		positions, err := currentBoard.config.Board.GetStartingPositions(color)
		if err != nil {
			return fmt.Errorf("failed to get starting positions for %s: %w", color, err)
		}
		for _, pos := range positions {
			pieceConfig, err := currentBoard.config.GetPieceConfig(pos.Name)
			if err != nil {
				return fmt.Errorf("failed to get config for piece %s: %w", pos.Name, err)
			}

			piece := &Piece{
				Name:   pos.Name,
				Color:  color,
				Config: *pieceConfig,
			}

			newBoard, err := currentBoard.PlacePiece(piece, pos.Position)
			if err != nil {
				return fmt.Errorf("failed to place %s %s at %v: %w", color, pos.Name, pos.Position, err)
			}
			currentBoard = newBoard
		}
	}
	// Copy the accumulated state back to the original board
	b.pieces = currentBoard.pieces
	b.captured = currentBoard.captured
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

// Copy creates a new Board with the same state as the receiver. The pieces and captured maps are
// deep copied, but just the pointers to the squares and config are copied.
func (b *Board) Copy() *Board {
	// Deep copy the pieces grid
	newPieces := make([][]*Piece, b.Rows)
	for i := range b.pieces {
		newPieces[i] = make([]*Piece, b.Columns)
		copy(newPieces[i], b.pieces[i])
	}

	// Deep copy the captured map
	newCaptured := make(map[Color][]*Piece)
	for color, pieces := range b.captured {
		newCaptured[color] = make([]*Piece, len(pieces))
		copy(newCaptured[color], pieces)
	}

	return &Board{
		Rows:     b.Rows,
		Columns:  b.Columns,
		squares:  b.squares,
		pieces:   newPieces,
		captured: newCaptured,
		config:   b.config,
	}
}

// PlacePiece puts the specified piece in the specified location, returning a new board with the change applied.
// Returns an error if the position is occupied or out of bounds.
func (b *Board) PlacePiece(piece *Piece, pos Position) (*Board, error) {
	if piece == nil {
		return nil, fmt.Errorf("piece is nil")
	}
	if !b.IsValid(pos) {
		return nil, fmt.Errorf("%s is out of bounds", pos)
	}
	if b.IsOccupied(pos) {
		return nil, fmt.Errorf("%s is occupied", pos)
	}
	newBoard := b.Copy()
	newBoard.pieces[pos[0]][pos[1]] = piece
	return newBoard, nil
}

// MovePiece moves the existing piece from the specified position, returning a new board with the move applied.
// An error is returned if the move isn't valid for the piece. If the destination is occupied by an opponent's piece,
// that piece is captured in the returned board.
func (b *Board) MovePiece(piece *Piece, start Position, move Move) (*Board, error) {
	// Make sure the piece is actually at that starting position.
	if b.pieces[start[0]][start[1]] != piece {
		return nil, fmt.Errorf("%s is not at %s", piece, start)
	}
	// Make sure that the move being passed in is valid for this piece from the starting position.
	validMoves := piece.ValidMoves(start, b)
	end := start.Add(move)
	if !lo.Contains(validMoves, end) {
		return nil, fmt.Errorf("move %v is not valid for piece %s", move, piece)
	}

	// Copy the board to make modifications
	newBoard := b.Copy()

	// Check if there's an opponent's piece at the destination to capture
	if occupant := newBoard.GetPieceAt(end); occupant != nil {
		newBoard.captured[piece.Color] = append(newBoard.captured[piece.Color], occupant)
		newBoard.pieces[end[0]][end[1]] = nil
	}

	// Put the piece in its new location and clear the old one
	newBoard.pieces[end[0]][end[1]] = piece
	newBoard.pieces[start[0]][start[1]] = nil

	return newBoard, nil
}

// GetSquareAt returns the square at the given position.
func (b *Board) GetSquareAt(pos Position) *Square {
	return &b.squares.data[pos[0]][pos[1]]
}

// GetPieceAt returns the piece at the given position, or nil if empty.
func (b *Board) GetPieceAt(pos Position) *Piece {
	return b.pieces[pos[0]][pos[1]]
}

// GetCapturedPieces returns all pieces captured by the specified color.
func (b *Board) GetCapturedPieces(color Color) []*Piece {
	return b.captured[color]
}

// GetPiecesByColor returns all pieces of the specified color on the board.
// Returns an empty slice if no pieces of that color are found.
func (b *Board) GetPiecesByColor(color Color) []*Piece {
	var pieces []*Piece
	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Columns; col++ {
			if piece := b.pieces[row][col]; piece != nil && piece.Color == color {
				pieces = append(pieces, piece)
			}
		}
	}
	return pieces
}
