// Copyright 2025 Ideograph LLC. All rights reserved.

package ai

import (
	"cragspider-go/internal/core"
	"fmt"
)

// BoardScorer is a way to evaluate who is winning the game just by looking at the current board state
type BoardScorer struct {
	config *AIPlayerConfig
}

// NewBoardScorer creates and returns a new BoardScorer structure, given the name of the AI player from the
// AI configuration file.
func NewBoardScorer(playerName string) (*BoardScorer, error) {
	cfg, err := GetAIConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot get AI config: %w", err)
	}
	playerConfig, err := cfg.GetPlayerConfig(playerName)
	if err != nil {
		return nil, fmt.Errorf("cannot get player config: %w", err)
	}
	return &BoardScorer{config: playerConfig}, nil
}

// Score takes a pointer to a Board with pieces on it and returns a float representing which player
// is in better shape. Positive numbers mean that White is winning; Negative numbers mean that Black
// is winning.
func (bs *BoardScorer) Score(board *core.Board) (float32, error) {
	// For now just add up the various pieces on the board valuing them appropriately
	var score float32

	for _, piece := range board.GetPiecesByColor(core.White) {
		score += bs.config.Scoring[piece.Name]
	}
	for _, piece := range board.GetPiecesByColor(core.Black) {
		score -= bs.config.Scoring[piece.Name]
	}

	return score, nil
}
