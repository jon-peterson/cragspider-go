// Copyright 2025 Ideograph LLC. All rights reserved.

package ai

import (
	"testing"

	"cragspider-go/internal/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoardScorer_ValidPlayerName(t *testing.T) {
	scorer, err := NewBoardScorer("doofus")

	require.NoError(t, err, "should create scorer with valid player name")
	assert.NotNil(t, scorer, "scorer should not be nil")
	assert.NotNil(t, scorer.config, "scorer config should not be nil")
	assert.Equal(t, "doofus", scorer.config.Name, "scorer should have correct player name")
}

func TestNewBoardScorer_InvalidPlayerName(t *testing.T) {
	scorer, err := NewBoardScorer("nonexistent")

	assert.Error(t, err, "should return error for invalid player name")
	assert.Nil(t, scorer, "scorer should be nil on error")
}

func TestScore_CustomPieces(t *testing.T) {
	// Create a game to get a valid board
	game, err := core.NewGame()
	require.NoError(t, err, "should create new game")

	// Count existing pieces on the board to calculate baseline
	whitePieces := game.Board.GetPiecesByColor(core.White)
	blackPieces := game.Board.GetPiecesByColor(core.Black)

	// Create scorer to get piece values
	scorer, err := NewBoardScorer("doofus")
	require.NoError(t, err, "should create scorer")

	// Calculate expected score from initial board
	var expectedScore float32
	for _, piece := range whitePieces {
		expectedScore += scorer.config.Scoring[piece.Name]
	}
	for _, piece := range blackPieces {
		expectedScore -= scorer.config.Scoring[piece.Name]
	}

	// Score the board
	score, err := scorer.Score(game.Board)
	require.NoError(t, err, "should score board without error")

	// Verify the calculation matches
	assert.Equal(t, expectedScore, score, "score should match manual calculation")

	// Now add a piece and verify the score changes correctly
	// Find an empty position
	emptyPos := core.Position{3, 3}
	if game.Board.IsOccupied(emptyPos) {
		// Try another position if this one is occupied
		emptyPos = core.Position{4, 4}
	}
	require.False(t, game.Board.IsOccupied(emptyPos), "test position should be empty")

	// Add a white warrior (value = 1)
	newWarrior := &core.Piece{Name: "warrior", Color: core.White}
	newBoard, err := game.Board.PlacePiece(newWarrior, emptyPos)
	require.NoError(t, err, "should place piece")

	// Score should increase by 1
	newScore, err := scorer.Score(newBoard)
	require.NoError(t, err, "should score new board")
	assert.Equal(t, expectedScore+1, newScore, "score should increase by warrior value (1)")

	// Add a black padwar (value = 2)
	emptyPos2 := core.Position{3, 4}
	require.False(t, newBoard.IsOccupied(emptyPos2), "second test position should be empty")

	blackPadwar := &core.Piece{Name: "padwar", Color: core.Black}
	newBoard2, err := newBoard.PlacePiece(blackPadwar, emptyPos2)
	require.NoError(t, err, "should place second piece")

	// Score should decrease by 2 from the previous score (because it's a black piece)
	finalScore, err := scorer.Score(newBoard2)
	require.NoError(t, err, "should score final board")
	assert.Equal(t, expectedScore+1-2, finalScore, "score should be adjusted by both pieces: +1 (white warrior) -2 (black padwar)")
}
