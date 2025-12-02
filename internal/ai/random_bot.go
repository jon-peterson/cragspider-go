package ai

import (
	"cragspider-go/internal/core"
	"fmt"
	"math/rand"
)

// RandomBot is an AI agent that makes random valid moves.
type RandomBot struct {
	Color core.Color
}

// NextMove returns a random valid move for a random piece of the bot's color.
// Returns nil and an error if no valid moves are available for any piece.
func (rb *RandomBot) NextMove(board *core.Board) (*Action, error) {
	// Get all pieces of the bot's color
	botPieces := board.GetPiecesByColor(rb.Color)

	if len(botPieces) == 0 {
		return nil, fmt.Errorf("no pieces found for color %s", rb.Color)
	}

	// Shuffle the pieces to randomize which one we try first
	rand.Shuffle(len(botPieces), func(i, j int) {
		botPieces[i], botPieces[j] = botPieces[j], botPieces[i]
	})

	// For each piece, try to find a valid move
	for _, piece := range botPieces {
		pos, err := board.PieceLocation(piece)
		if err != nil {
			continue // Skip if we can't find the piece location
		}

		validPositions := piece.ValidMoves(pos, board)
		if len(validPositions) == 0 {
			continue
		}

		// Randomly select one of the valid positions
		selectedPos := validPositions[rand.Intn(len(validPositions))]
		move := calculateDelta(pos, selectedPos)

		return &Action{
			Piece: piece,
			Move:  move,
		}, nil
	}

	return nil, fmt.Errorf("no valid moves available for color %s", rb.Color)
}

// calculateDelta computes the move delta from a starting position to an ending position.
func calculateDelta(from, to core.Position) core.Move {
	return core.Move{
		to[0] - from[0],
		to[1] - from[1],
	}
}
