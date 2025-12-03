// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	game, err := NewGame()
	require.NoError(t, err)

	// Check that the game was created with the correct initial state
	assert.NotNil(t, game.Board, "Board should be initialized")
	assert.Equal(t, White, game.ActiveColor, "White should be the current player")
	assert.NotNil(t, game.players, "Players map should be initialized")
}

func TestGamePlayers(t *testing.T) {
	game, err := NewGame()
	require.NoError(t, err)

	t.Run("white player is human", func(t *testing.T) {
		whitePlayer := game.GetPlayer(White)
		require.NotNil(t, whitePlayer, "White player should exist")
		assert.True(t, whitePlayer.IsHuman(), "White player should be human")
	})

	t.Run("black player is human", func(t *testing.T) {
		blackPlayer := game.GetPlayer(Black)
		require.NotNil(t, blackPlayer, "Black player should exist")
		assert.True(t, blackPlayer.IsHuman(), "Black player should be human")
	})
}

func TestAdvanceTurn(t *testing.T) {
	game, err := NewGame()
	require.NoError(t, err)

	assert.Equal(t, White, game.ActiveColor, "White should be the first player")

	game.AdvanceTurn()
	assert.Equal(t, Black, game.ActiveColor, "Black should be the current player after advancing")

	game.AdvanceTurn()
	assert.Equal(t, White, game.ActiveColor, "White should be the current player after advancing twice")
}
