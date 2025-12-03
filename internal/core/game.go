// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game represents a single instance of a game.
type Game struct {
	Board       *Board
	config      *GameConfig
	ActiveColor Color
	players     map[Color]*Player
}

// NewGame returns a new game with the standard configuration.
func NewGame() (*Game, error) {
	cfg, err := GetConfig()
	if err != nil {
		rl.TraceLog(rl.LogFatal, "error loading default configuration: %v", err)
	}
	return NewGameWithConfig(cfg)
}

// NewGameWithConfig returns a new game with a board ready to be played. It accepts a GameConfig parameter to allow
// custom configurations.
func NewGameWithConfig(cfg *GameConfig) (*Game, error) {
	if cfg == nil {
		return nil, fmt.Errorf("game configuration cannot be nil")
	}
	b, err := newBoard(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	return &Game{
		Board:       b,
		config:      cfg,
		ActiveColor: White,
		players: map[Color]*Player{
			White: newHumanPlayer(),
			Black: newHumanPlayer(),
		},
	}, nil
}

// Over returns true if the game is over.
func (g *Game) Over() bool {
	return false
}

// AdvanceTurn advances the game to the next player's turn.
func (g *Game) AdvanceTurn() {
	if g.ActiveColor == White {
		g.ActiveColor = Black
	} else {
		g.ActiveColor = White
	}
}

// GetPlayer returns the player with the given color.
func (g *Game) GetPlayer(color Color) *Player {
	return g.players[color]
}

// SetPlayer sets the player for the given color.
func (g *Game) SetPlayer(color Color, player *Player) {
	g.players[color] = player
}

// ActionToMove converts an Action (with absolute destination Position) to a Move (with delta).
// This is needed because agents return absolute positions, but Board.MovePiece expects moves with deltas.
func (g *Game) ActionToMove(action *Action) (Move, error) {
	currentPos, err := g.Board.PieceLocation(action.Piece)
	if err != nil {
		return Move{}, fmt.Errorf("cannot find piece location: %w", err)
	}
	return Move{
		action.Destination[0] - currentPos[0],
		action.Destination[1] - currentPos[1],
	}, nil
}
