package ai

import (
	"fmt"
	"testing"

	"cragspider-go/internal/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomBotNextMove_ReturnsValidMove(t *testing.T) {
	// Create a new game with default config
	game, err := core.NewGame()
	require.NoError(t, err, "should create new game")

	bot := &RandomBot{Color: core.White}

	// Get a move from the bot
	action, err := bot.NextMove(game.Board)

	// Verify we got a valid action
	assert.NoError(t, err, "should return no error")
	assert.NotNil(t, action, "should return an action")
	assert.NotNil(t, action.Piece, "action should have a piece")
	assert.Equal(t, core.White, action.Piece.Color, "piece should be white")

	// Verify the move is valid by checking that the destination is in the valid moves
	startPos, err := game.Board.PieceLocation(action.Piece)
	require.NoError(t, err, "should find piece location")

	validPositions := action.Piece.ValidMoves(startPos, game.Board)
	destination := startPos.Add(action.Move)
	assert.Contains(t, validPositions, destination, "destination should be in valid positions for the piece")
}

func TestRandomBotNextMove_DifferentColorPieces(t *testing.T) {
	game, err := core.NewGame()
	require.NoError(t, err, "should create new game")

	// Test black bot
	blackBot := &RandomBot{Color: core.Black}
	action, err := blackBot.NextMove(game.Board)

	assert.NoError(t, err, "should return no error for black bot")
	assert.NotNil(t, action, "should return an action")
	assert.Equal(t, core.Black, action.Piece.Color, "piece should be black")
}

func TestRandomBotMultipleCalls(t *testing.T) {
	game, err := core.NewGame()
	require.NoError(t, err, "should create new game")

	bot := &RandomBot{Color: core.White}

	// Call NextMove multiple times to check randomness
	moves := make(map[string]bool)
	for i := 0; i < 5; i++ {
		action, err := bot.NextMove(game.Board)
		assert.NoError(t, err)
		assert.NotNil(t, action)

		// Create a unique key for this move
		moveKey := fmt.Sprintf("%s -> %s", action.Piece.Name, action.Move.String())
		moves[moveKey] = true
	}

	// We expect some variation (though with randomness, it's possible all 5 are the same)
	// This is just a sanity check that the code runs without errors
	assert.GreaterOrEqual(t, len(moves), 1, "should have at least one move")
}
