package ai

import "cragspider-go/internal/core"

// Action represents a complete move action, containing both the piece and the move delta.
type Action struct {
	Piece *core.Piece
	Move  core.Move
}

// Agent is an interface for bot implementations that generate moves.
type Agent interface {
	// NextMove returns the next move for the agent given the current board state.
	// It returns an Action (containing the piece and move delta) or an error if no valid moves are available.
	NextMove(board *core.Board) (*Action, error)
}
