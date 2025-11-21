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
// For each path, the piece can move to any position along that path until it encounters a blocking piece:
// - If a same-color piece blocks, the path ends and that position cannot be moved to
// - If an opposite-color piece blocks, the piece can capture it but cannot continue past
func (p *Piece) ValidMoves(start Position, b *Board) []Position {
	moves := make([]Position, 0)

	// Process each path independently
	for _, path := range p.Config.Moves {
		// Skip empty paths
		if len(path) == 0 {
			continue
		}

		currentPos := start

		// Walk along the path, one delta at a time
		for _, delta := range path {
			nextPos := Position{
				currentPos[0] + delta[0],
				currentPos[1] + delta[1],
			}

			// If the next position is off the board, this path ends
			if !b.IsValid(nextPos) {
				break
			}

			// Check what's at the next position
			occupant := b.GetPieceAt(nextPos)

			// If a same-color piece is there, we cannot move to this position or continue
			if occupant != nil && occupant.Color == p.Color {
				break
			}

			// We can move to this position
			moves = append(moves, nextPos)

			// If an opposite-color piece is there, we can capture but cannot continue
			if occupant != nil && occupant.Color != p.Color {
				break
			}

			// No piece is blocking, continue along the path
			currentPos = nextPos
		}
	}

	return moves
}
