package ai

import "cragspider-go/internal/core"

// Action represents a complete move action, containing the piece and its destination position.
type Action struct {
	Piece       *core.Piece
	Destination core.Position
}

// Agent is an interface for bot implementations that generate moves.
type Agent interface {
	// NextMove returns the next move for the agent given the current board state.
	// It returns an Action (containing the piece and move delta) or an error if no valid moves are available.
	NextMove(board *core.Board) (*Action, error)
}
