// Copyright 2025 Ideograph LLC. All rights reserved.

package core

// Color represents the color of a piece or player.
type Color string

const (
	// White represents the white player/pieces.
	White Color = "white"
	// Black represents the black player/pieces.
	Black Color = "black"
)

// Piece is a white or black piece on the board with its associated data
type Piece struct {
	name     string
	color    Color
	config   PieceConfig
	selected bool
}

// ToggleSelected toggles the selected state of the piece.
func (p *Piece) ToggleSelected() {
	p.selected = !p.selected
}

// ValidMoves returns a list of valid positions that the piece can move to given the starting position.
func (p *Piece) ValidMoves(start Position, b *Board) []Position {
	moves := make([]Position, 0, len(p.config.Moves))

	for _, move := range p.config.Moves {
		pos := Position{
			start[0] + move[0],
			start[1] + move[1],
		}
		moves = append(moves, pos)
	}

	return moves
}
