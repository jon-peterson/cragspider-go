// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import "fmt"

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
	Name   string
	Color  Color
	Config PieceConfig
}

// String returns a nicely formatted string representation of the piece.
func (p Piece) String() string {
	return fmt.Sprintf("%s %s", p.Color, p.Name)
}

// ValidMoves returns a list of valid positions that the piece can move to from the given starting position.
func (p *Piece) ValidMoves(start Position, b *Board) []Position {
	moves := make([]Position, 0, len(p.Config.Moves))

	for _, move := range p.Config.Moves {
		pos := Position{
			start[0] + move[0],
			start[1] + move[1],
		}
		if b.IsValid(pos) {
			moves = append(moves, pos)
		}
	}

	return moves
}
