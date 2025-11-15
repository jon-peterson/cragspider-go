// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	"github.com/samber/lo"
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
	tests := []struct {
		name           string
		startPos       Position
		expectedMoves  []Position
		expectedLength int
	}{
		{
			name:     "center position",
			startPos: Position{2, 2},
			expectedMoves: []Position{
				{2, 3}, // One step down
				{2, 4}, // Two steps down
				{2, 1}, // One step up
				{2, 0}, // Two steps up
				{3, 2}, // One step right
				{4, 2}, // Two steps right
				{1, 2}, // One step left
				{0, 2}, // Two steps left
			},
			expectedLength: 8,
		},
		{
			name:     "left edge position",
			startPos: Position{2, 0},
			expectedMoves: []Position{
				{2, 1}, // One step down
				{2, 2}, // Two steps down
				{3, 0}, // One step right
				{4, 0}, // Two steps right
				{1, 0}, // One step left
				{0, 0}, // Two steps left
				// Note: Moves that would go off the top (y < 0) are filtered out
			},
			expectedLength: 6, // 2 moves filtered out (y = -1 and y = -2)
		},
	}

	// Create a test piece that can move two spaces in each cardinal direction
	piece := getTestPiece()

	// Create an empty 5x5 board (0,0 is top-left)
	board := createTestBoard(5, 5)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validMoves := piece.ValidMoves(tt.startPos, board)

			// Check that we got the expected number of moves
			assert.Len(t, validMoves, tt.expectedLength, "Unexpected number of valid moves")

			// Check that all returned moves are within board bounds
			for _, move := range validMoves {
				assert.True(t, move[0] >= 0 && move[0] < 5, "Move %v has invalid x coordinate (should be 0-4)", move)
				assert.True(t, move[1] >= 0 && move[1] < 5, "Move %v has invalid y coordinate (should be 0-4)", move)
			}

			// Check each expected move is in the valid moves
			for _, expected := range tt.expectedMoves {
				assert.True(t, lo.Contains(validMoves, expected), "Expected move %v not found in valid moves", expected)
			}
		})
	}
}

// getTestPiece returns a standard piece that can move two steps in each cardinal direction
func getTestPiece() *Piece {
	piece := &Piece{
		Name:  "test_piece",
		Color: White,
		Config: PieceConfig{
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
