// Copyright 2025 Ideograph LLC. All rights reserved.

package ai

import (
	"cragspider-go/internal/core"
	"cragspider-go/pkg/random"
	"fmt"
	"math/rand"
)

// RandomBot is an AI agent that makes random valid moves.
type RandomBot struct {
	Color core.Color
}

// NewRandomBot returns a new RandomBot structure for the specified color.
func NewRandomBot(color core.Color) *RandomBot {
	return &RandomBot{Color: color}
}

// NextMove returns a random valid move for a random piece of the bot's color.
// Returns nil and an error if no valid moves are available for any piece.
func (rb *RandomBot) NextMove(board *core.Board) (*core.Action, error) {
	// Get all pieces of the bot's color
	botPieces := board.GetPiecesByColor(rb.Color)

	if len(botPieces) == 0 {
		return nil, fmt.Errorf("no pieces found for color %s", rb.Color)
	}

	// Shuffle the pieces to randomize which one we try first
	rand.Shuffle(len(botPieces), func(i, j int) {
		botPieces[i], botPieces[j] = botPieces[j], botPieces[i]
	})

	// As soon as there's a piece with at least one valid move, return one of them randomly
	for _, piece := range botPieces {
		pos, err := board.PieceLocation(piece)
		if err != nil {
			continue // Skip if we can't find the piece location
		}

		validPositions := piece.ValidNextPositions(pos, board)
		if len(validPositions) == 0 {
			continue
		}

		selectedPos := random.Choice(validPositions)
		move := core.Move{
			selectedPos[0] - pos[0],
			selectedPos[1] - pos[1],
		}
		return &core.Action{
			Piece: piece,
			Move:  move,
		}, nil
	}

	return nil, fmt.Errorf("no valid moves available for color %s", rb.Color)
}
