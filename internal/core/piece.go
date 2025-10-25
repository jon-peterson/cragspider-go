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
