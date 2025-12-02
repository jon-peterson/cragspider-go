// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// createTestBoard creates a test board with the specified dimensions
func createTestBoard(rows, cols int) *Board {
	board := &Board{
		Rows:    rows,
		Columns: cols,
		squares: newSquareGrid(rows, cols),
		pieces:  make([][]*Piece, rows),
	}
	for i := range board.pieces {
		board.pieces[i] = make([]*Piece, cols)
	}
	return board
}

func TestPiece_validMoves(t *testing.T) {
	tests := []struct {
		name           string
		startPos       Position
		expectedLength int
		expectedMoves  []Position           // Exact positions expected in valid moves
		setup          func(*Board, *Piece) // Optional setup to place pieces on board
	}{
		{
			name:           "center position - all paths open",
			startPos:       Position{2, 2},
			expectedLength: 8,
			expectedMoves: []Position{
				{3, 2}, {4, 2}, // Path down
				{1, 2}, {0, 2}, // Path up
				{2, 3}, {2, 4}, // Path right
				{2, 1}, {2, 0}, // Path left
			},
		},
		{
			name:           "left edge position - left path blocked",
			startPos:       Position{2, 0},
			expectedLength: 6,
			expectedMoves: []Position{
				{3, 0}, {4, 0}, // Path down
				{1, 0}, {0, 0}, // Path up
				{2, 1}, {2, 2}, // Path right
				// Path left blocked immediately by board edge
			},
		},
		{
			name:           "same-color piece blocks path at step 2",
			startPos:       Position{2, 2},
			expectedLength: 4, // Can only do first step on each path
			expectedMoves: []Position{
				{3, 2}, // Path down, step 1 only (blocked by piece at {4,2})
				{1, 2}, // Path up, step 1 only (blocked by piece at {0,2})
				{2, 3}, // Path right, step 1 only (blocked by piece at {2,4})
				{2, 1}, // Path left, step 1 only (blocked by piece at {2,0})
			},
			setup: func(board *Board, piece *Piece) {
				// Block second step of each path with same-color pieces
				blocker1 := &Piece{Name: "blocker1", Color: White}
				blocker2 := &Piece{Name: "blocker2", Color: White}
				blocker3 := &Piece{Name: "blocker3", Color: White}
				blocker4 := &Piece{Name: "blocker4", Color: White}

				board.pieces[4][2] = blocker1 // Blocks down path at step 2
				board.pieces[0][2] = blocker2 // Blocks up path at step 2
				board.pieces[2][4] = blocker3 // Blocks right path at step 2
				board.pieces[2][0] = blocker4 // Blocks left path at step 2
			},
		},
		{
			name:           "same-color piece blocks path at step 1",
			startPos:       Position{2, 2},
			expectedLength: 2, // Only right path is open, with 2 moves
			expectedMoves: []Position{
				{2, 3}, {2, 4}, // Path right: no blockers
				// Paths down, up, left blocked at step 1
			},
			setup: func(board *Board, piece *Piece) {
				// Block first step of down, up, and left paths
				blocker1 := &Piece{Name: "blocker1", Color: White}
				blocker2 := &Piece{Name: "blocker2", Color: White}
				blocker3 := &Piece{Name: "blocker3", Color: White}

				board.pieces[3][2] = blocker1 // Blocks down path at step 1
				board.pieces[1][2] = blocker2 // Blocks up path at step 1
				board.pieces[2][1] = blocker3 // Blocks left path at step 1
			},
		},
		{
			name:           "opposite-color piece allows capture and ends path",
			startPos:       Position{2, 2},
			expectedLength: 6, // Can capture on step 2 for two paths, all steps for one path
			expectedMoves: []Position{
				{3, 2},         // Path down, step 1
				{4, 2},         // Path down, step 2 (opposite piece - can capture)
				{1, 2},         // Path up, step 1
				{0, 2},         // Path up, step 2 (opposite piece - can capture)
				{2, 3}, {2, 4}, // Path right: no blockers
				// Path left is blocked at step 1 by same-color piece
			},
			setup: func(board *Board, piece *Piece) {
				// Place opposite-color pieces at step 2 of down and up paths
				enemy1 := &Piece{Name: "enemy1", Color: Black}
				enemy2 := &Piece{Name: "enemy2", Color: Black}
				blocker := &Piece{Name: "blocker", Color: White}

				board.pieces[4][2] = enemy1  // Down path step 2: opposite color
				board.pieces[0][2] = enemy2  // Up path step 2: opposite color
				board.pieces[2][1] = blocker // Left path step 1: same color
			},
		},
	}

	// Create a test piece with 4 cardinal paths
	piece := getTestPiece()

	// Create an empty 5x5 board
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
				assert.True(t, move[0] >= 0 && move[0] < 5, "Move %v has invalid row coordinate (should be 0-4)", move)
				assert.True(t, move[1] >= 0 && move[1] < 5, "Move %v has invalid col coordinate (should be 0-4)", move)
			}

			// Verify expected moves are all present
			for _, expected := range tt.expectedMoves {
				assert.True(t, lo.Contains(validMoves, expected), "Expected move %v not found in valid moves", expected)
			}
		})
	}
}

// getTestPiece returns a standard piece with 4 cardinal paths: up to 2 steps in each direction
func getTestPiece() *Piece {
	piece := &Piece{
		Name:  "test_piece",
		Color: White,
		Config: PieceConfig{
			Name: "test_piece",
			Moves: [][]Move{
				// Path down: can move 1 or 2 steps down
				{{1, 0}, {1, 0}},
				// Path up: can move 1 or 2 steps up
				{{-1, 0}, {-1, 0}},
				// Path right: can move 1 or 2 steps right
				{{0, 1}, {0, 1}},
				// Path left: can move 1 or 2 steps left
				{{0, -1}, {0, -1}},
			},
		},
	}
	return piece
}
