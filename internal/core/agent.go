// Copyright 2025 Ideograph LLC. All rights reserved.

package core

// Action represents a complete move action, containing the piece and its move delta.
type Action struct {
	Piece *Piece
	Move  Move
}

// AgentStrategy is an interface for bot implementations that generate moves.
type AgentStrategy interface {
	// NextMove returns the next move for the strategy given the current board state.
	// It returns an Action (containing the piece and move delta) or an error if no valid moves are available.
	NextMove(board *Board) (*Action, error)
}
