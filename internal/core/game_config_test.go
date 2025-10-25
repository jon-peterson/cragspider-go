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
		assert.Len(t, warrior.Moves, 8, "warrior should have 8 possible moves")
		assert.Contains(t, warrior.Sprites, White, "warrior should have white sprites")
		assert.Contains(t, warrior.Sprites, Black, "warrior should have black sprites")

		// Test board configuration
		whitePieces, err := cfg.Board.GetStartingPositions(White)
		require.NoError(t, err)
		assert.Len(t, whitePieces, 4, "should have 4 white pieces")

		blackPieces, err := cfg.Board.GetStartingPositions(Black)
		require.NoError(t, err)
		assert.Len(t, blackPieces, 4, "should have 4 black pieces")
	})
}
