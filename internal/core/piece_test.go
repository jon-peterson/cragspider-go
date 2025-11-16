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
		expectedLength int
		blockedMoves   []Position           // Moves that should be blocked (for documentation)
		setup          func(*Board, *Piece) // Optional setup to place pieces on board
	}{
		{
			name:           "center position",
			startPos:       Position{2, 2},
			expectedLength: 8,
			blockedMoves:   []Position{},
		},
		{
			name:           "left edge position",
			startPos:       Position{2, 0},
			expectedLength: 6, // Moves that go off board (y < 0) are filtered
			blockedMoves:   []Position{{2, -1}, {2, -2}},
		},
		{
			name:           "blocked by same color piece",
			startPos:       Position{2, 2},
			expectedLength: 4, // Two-step moves are blocked by same-color pieces
			blockedMoves:   []Position{{2, 4}, {2, 0}, {4, 2}, {0, 2}},
			setup: func(board *Board, piece *Piece) {
				// Place same-color pieces at two-step positions
				samePiece1 := &Piece{Name: "blocker1", Color: White}
				samePiece2 := &Piece{Name: "blocker2", Color: White}
				samePiece3 := &Piece{Name: "blocker3", Color: White}
				samePiece4 := &Piece{Name: "blocker4", Color: White}

				board.pieces[2][4] = samePiece1 // Two steps down
				board.pieces[2][0] = samePiece2 // Two steps up
				board.pieces[4][2] = samePiece3 // Two steps right
				board.pieces[0][2] = samePiece4 // Two steps left
			},
		},
		{
			name:           "opposite color pieces do not block",
			startPos:       Position{2, 2},
			expectedLength: 8, // All moves including enemy-occupied ones are valid
			blockedMoves:   []Position{},
			setup: func(board *Board, piece *Piece) {
				// Place opposite-color pieces at two-step positions
				enemyPiece1 := &Piece{Name: "enemy1", Color: Black}
				enemyPiece2 := &Piece{Name: "enemy2", Color: Black}
				enemyPiece3 := &Piece{Name: "enemy3", Color: Black}
				enemyPiece4 := &Piece{Name: "enemy4", Color: Black}

				board.pieces[2][4] = enemyPiece1 // Two steps down
				board.pieces[2][0] = enemyPiece2 // Two steps up
				board.pieces[4][2] = enemyPiece3 // Two steps right
				board.pieces[0][2] = enemyPiece4 // Two steps left
			},
		},
	}

	// Create a test piece that can move two spaces in each cardinal direction
	piece := getTestPiece()

	// Create an empty 5x5 board (0,0 is top-left)
	board := createTestBoard(5, 5)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset board pieces for this test
			for i := range board.pieces {
				board.pieces[i] = make([]*Piece, 5)
			}

			// Run any test-specific setup
			if tt.setup != nil {
				tt.setup(board, piece)
			}

			validMoves := piece.ValidMoves(tt.startPos, board)

			// Check that we got the expected number of moves
			assert.Len(t, validMoves, tt.expectedLength, "Unexpected number of valid moves")

			// Check that all returned moves are within board bounds
			for _, move := range validMoves {
				assert.True(t, move[0] >= 0 && move[0] < 5, "Move %v has invalid x coordinate (should be 0-4)", move)
				assert.True(t, move[1] >= 0 && move[1] < 5, "Move %v has invalid y coordinate (should be 0-4)", move)
			}

			// Verify blocked moves are not in the valid moves
			for _, blocked := range tt.blockedMoves {
				if board.IsValid(blocked) {
					assert.False(t, lo.Contains(validMoves, blocked), "Blocked move %v should not be in valid moves", blocked)
				}
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
