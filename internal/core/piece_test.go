// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTestBoard creates a test board with the specified dimensions
func createTestBoard(rows, cols int) *Board {
	return &Board{
		Rows:    rows,
		Columns: cols,
		squares: make([][]Square, rows),
		pieces:  make([][]*Piece, rows),
	}
}

func TestPiece_validMoves(t *testing.T) {
	// Create a test piece that can move two spaces in each cardinal direction
	piece := getTestPiece()

	// Create an empty 5x5 board (0,0 is top-left)
	board := createTestBoard(5, 5)

	// Get valid moves from the center of the board (row 2, col 2)
	validMoves := piece.ValidMoves(Position{2, 2}, board)

	// Expected moves (one and two steps in each cardinal direction from position [2,2])
	expectedMoves := []Position{
		{2, 3}, // One step down
		{2, 4}, // Two steps down
		{2, 1}, // One step up
		{2, 0}, // Two steps up
		{3, 2}, // One step right
		{4, 2}, // Two steps right
		{1, 2}, // One step left
		{0, 2}, // Two steps left
	}

	// Check that we got the expected number of moves
	assert.Len(t, validMoves, len(expectedMoves), "Unexpected number of valid moves")

	// Check each expected move is in the valid moves
	for _, expected := range expectedMoves {
		found := false
		for _, move := range validMoves {
			if move[0] == expected[0] && move[1] == expected[1] {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected move %v not found in valid moves", expected)
	}
}

// getTestPiece returns a standard piece that can move two steps in each cardinal direction
func getTestPiece() *Piece {
	piece := &Piece{
		name:  "test_piece",
		color: White,
		config: PieceConfig{
			Name: "test_piece",
			Moves: []Move{
				{1, 0},
				{2, 0},
				{-1, 0},
				{-2, 0},
				{0, 1},
				{0, 2},
				{0, -1},
				{0, -2},
			},
		},
	}
	return piece
}
