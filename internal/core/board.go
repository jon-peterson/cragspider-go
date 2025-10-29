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
	frame    animation.FrameCoords
	rotation rl.Vector2
}

// Position is a [row,col] on the board.
type Position [2]int

// Move represents a [deltaRow, deltaCol] move on the board.
type Move [2]int

// selectedPieceAndPosition represents a selected piece on the board and its position.
type selectedPieceAndPosition struct {
	Piece    *Piece
	Position Position
}

// Board is the game board, which is a grid of squares upon which there are pieces.
type Board struct {
	Rows, Columns int
	squares       [][]Square
	pieces        [][]*Piece

	backgroundSprites *animation.SpriteSheet
	whiteSprites      *animation.SpriteSheet
	blackSprites      *animation.SpriteSheet

	selectedPiece *selectedPieceAndPosition

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

// PositionUnderClick returns the board position under a mouse click, given the board location in the window
// and where the user clicked. If the user clicked outside the board, then an error is returned.
func (b *Board) PositionUnderClick(boardLoc, clickLoc rl.Vector2) (Position, error) {
	// Shift the position relative to the board upper corner so the click loc is in board space
	adjClickLoc := rl.Vector2{X: clickLoc.X - boardLoc.X, Y: clickLoc.Y - boardLoc.Y}

	// Check if the click is outside the board bounds
	if adjClickLoc.X < 0 || adjClickLoc.X >= float32(SquareSize*b.Columns) ||
		adjClickLoc.Y < 0 || adjClickLoc.Y >= float32(SquareSize*b.Rows) {
		return Position{}, fmt.Errorf("click is outside the board bounds")
	}

	// Just scale the click based on the square size
	return Position{int(adjClickLoc.Y / float32(SquareSize)), int(adjClickLoc.X / float32(SquareSize))}, nil
}

// PieceUnderClick returns the piece under a mouse click, given the board location in the window
// and where the user clicked. If there's no piece there, then nil is returned.
func (b *Board) PieceUnderClick(boardLoc, clickLoc rl.Vector2) *Piece {
	pos, err := b.PositionUnderClick(boardLoc, clickLoc)
	if err != nil {
		return nil
	}
	return b.pieces[pos[0]][pos[1]]
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

// SelectPiece selects the specified piece, unselecting any previously selected piece.
func (b *Board) SelectPiece(p *Piece) {
	if p == nil {
		b.selectedPiece = nil
		return
	}
	// Selecting a selected piece unselects it (toggle)
	if b.selectedPiece != nil && b.selectedPiece.Piece == p {
		b.selectedPiece = nil
		return
	}
	// Find the position of the piece and store it in the struct
	pos, err := b.PieceLocation(p)
	if err != nil {
		b.selectedPiece = nil
		return
	}
	b.selectedPiece = &selectedPieceAndPosition{
		Piece:    p,
		Position: pos,
	}
}

// Render draws the board to the screen with the given board location (where the upper left corner is).
func (b *Board) Render(boardLoc rl.Vector2) error {
	// First draw the board itself
	for i := range b.Rows {
		for j := range b.Columns {
			err := b.backgroundSprites.DrawFrameRotated(
				b.squares[i][j].frame,
				rl.Vector2{X: boardLoc.X + float32(j*SquareSize), Y: boardLoc.Y + float32(i*SquareSize)},
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
				err2 := b.renderPieceOnBoard(boardLoc, piece, j, i)
				if err2 != nil {
					return err2
				}
			}
		}
	}

	return nil
}

// renderPieceOnBoard renders a single piece on the board at the specified position.
func (b *Board) renderPieceOnBoard(boardLoc rl.Vector2, piece *Piece, j int, i int) error {
	// Different sprite sheets for different players of course
	sheet := lo.Ternary(piece.color == White, b.whiteSprites, b.blackSprites)
	isSelected := b.selectedPiece != nil && b.selectedPiece.Piece == piece
	frame := lo.Ternary(isSelected, 0, 1)
	err := sheet.DrawFrame(
		piece.config.Sprites[piece.color][frame],
		rl.Vector2{X: boardLoc.X + float32(j*SquareSize), Y: boardLoc.Y + float32(i*SquareSize)},
		Scale)
	if err != nil {
		return fmt.Errorf("failed to draw piece: %w", err)
	}
	return nil
}
