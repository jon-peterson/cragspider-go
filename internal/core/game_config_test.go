package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		// Get the config and run the test function
		cfg, err := GetConfig()
		require.NoError(t, err)

		// Test piece loading
		assert.Len(t, cfg.Pieces, 2, "should load 2 pieces")

		// Test warrior piece
		warrior, err := cfg.GetPieceConfig("warrior")
		require.NoError(t, err)
		assert.Len(t, warrior.Moves, 4, "warrior should have 4 movement paths")
		// Each path should have 2 steps (can move 1 or 2 spaces orthogonally)
		for _, path := range warrior.Moves {
			assert.Len(t, path, 2, "each warrior path should have 2 steps")
		}
		assert.Contains(t, warrior.Sprites, White, "warrior should have white sprites")
		assert.Contains(t, warrior.Sprites, Black, "warrior should have black sprites")

		// Test board dimensions
		assert.Equal(t, 10, cfg.Board.Rows, "board should have 10 rows")
		assert.Equal(t, 10, cfg.Board.Columns, "board should have 10 columns")

		// Test board configuration
		whitePieces, err := cfg.Board.GetStartingPositions(White)
		require.NoError(t, err)
		assert.Len(t, whitePieces, 4, "should have 4 white pieces")

		blackPieces, err := cfg.Board.GetStartingPositions(Black)
		require.NoError(t, err)
		assert.Len(t, blackPieces, 4, "should have 4 black pieces")
	})
}
