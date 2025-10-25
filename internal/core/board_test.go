// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	// Create a new board
	cfg, err := GetConfig()
	require.NoError(t, err)
	board, err := newBoard(cfg)
	require.NoError(t, err)

	// Check board dimensions
	assert.Equal(t, 10, board.Rows, "Board should have 10 rows")
	assert.Equal(t, 10, board.Columns, "Board should have 10 columns")

	// Check slice dimensions
	require.Len(t, board.squares, 10, "Squares should have 10 rows")
	for _, row := range board.squares {
		assert.Len(t, row, 10, "Each cell row should have 10 columns")
	}
	require.Len(t, board.pieces, 10, "Pieces should have 10 rows")
	for _, row := range board.pieces {
		assert.Len(t, row, 10, "Each piece row should have 10 columns")
	}

	// Check that backgroundSprites is initialized
	assert.NotNil(t, board.backgroundSprites, "Background sprites should be initialized")

	// Check that each cell has a valid rotation vector
	for i, row := range board.squares {
		for j, cell := range row {
			// Check that rotation is one of the cardinal directions
			dir := Move{int(cell.rotation.X), int(cell.rotation.Y)}
			assert.Contains(t, CardinalDirections, dir,
				"Square at [%d][%d] has invalid rotation vector: %v", i, j, dir)
		}
	}

	// Verify that each piece in the game config is on that position in the board
	whitePieces, err := cfg.Board.GetStartingPositions(White)
	require.NoError(t, err)
	for _, white := range whitePieces {
		pieceOnBoard := board.pieces[white.Position[0]][white.Position[1]]
		assert.NotNil(t, pieceOnBoard, "white piece %s should be at position %v", white.Name, white.Position)
	}
	// Do the same for the black pieces
	blackPieces, err := cfg.Board.GetStartingPositions(Black)
	require.NoError(t, err)
	for _, black := range blackPieces {
		pieceOnBoard := board.pieces[black.Position[0]][black.Position[1]]
		assert.NotNil(t, pieceOnBoard, "black piece %s should be at position %v", black.Name, black.Position)
	}
}

func TestBoard_PlacePiece(t *testing.T) {
	cfg, err := GetConfig()
	require.NoError(t, err)
	board, err := newBoard(cfg)
	require.NoError(t, err)

	// Create a test piece
	piece := Piece{
		name:  "test",
		color: White,
	}

	t.Run("place piece in empty position", func(t *testing.T) {
		pos := Position{1, 1}
		err := board.PlacePiece(piece, pos)
		require.NoError(t, err, "Should be able to place piece in empty position")
		assert.NotNil(t, board.pieces[pos[0]][pos[1]], "Piece should be placed on the board")
		assert.Equal(t, piece.name, board.pieces[pos[0]][pos[1]].name, "Placed piece should have the correct name")
		assert.Equal(t, piece.color, board.pieces[pos[0]][pos[1]].color, "Placed piece should have the correct color")
	})

	t.Run("cannot place in occupied position", func(t *testing.T) {
		pos := Position{2, 2}
		// First placement should succeed
		err := board.PlacePiece(piece, pos)
		require.NoError(t, err, "First placement should succeed")

		// Second placement should fail
		err = board.PlacePiece(piece, pos)
		assert.Error(t, err, "Should not be able to place piece in occupied position")
	})

	t.Run("cannot place out of bounds", func(t *testing.T) {
		outOfBoundsPositions := []Position{
			{-1, 0}, // row too small
			{10, 0}, // row too large
			{0, -1}, // col too small
			{0, 10}, // col too large
		}

		for _, pos := range outOfBoundsPositions {
			err := board.PlacePiece(piece, pos)
			assert.Error(t, err, "Should not be able to place piece at position %v", pos)
		}
	})
}
