// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"testing"

	"cragspider-go/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
)

func TestPlayfield_getTintedPositions(t *testing.T) {
	// Create a test game with default configuration
	game, err := core.NewGame()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Create a test playfield with the game
	pf := &Playfield{
		game: game,
	}

	tests := []struct {
		name          string
		piece         *core.Piece
		position      core.Position
		expectedColor rl.Color
		expectedMoves []core.Position
	}{
		{
			name:          "no selected piece",
			piece:         nil,
			position:      core.Position{},
			expectedColor: rl.Color{},
			expectedMoves: []core.Position{},
		},
		{
			name: "with selected white piece",
			piece: &core.Piece{
				Name:  "test",
				Color: core.White,
				Config: core.PieceConfig{
					Moves: []core.Move{{-1, 0}, {1, 0}},
				},
			},
			position:      core.Position{3, 3},
			expectedColor: rl.Green,
			expectedMoves: []core.Position{{3, 3}, {2, 3}, {4, 3}},
		},
		{
			name: "with selected black piece",
			piece: &core.Piece{
				Name:  "test",
				Color: core.Black,
				Config: core.PieceConfig{
					Moves: []core.Move{{-1, 0}, {1, 0}},
				},
			},
			position:      core.Position{3, 3},
			expectedColor: rl.Red,
			expectedMoves: []core.Position{{3, 3}, {2, 3}, {4, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.piece == nil {
				pf.selectedPiece = nil
			} else {
				pf.selectedPiece = &SelectedPieceAndPosition{
					Piece:    tt.piece,
					Position: tt.position,
				}
			}

			tints := pf.getTintedPositions()

			if tt.piece == nil {
				assert.Empty(t, tints, "Expected no tints when no piece is selected")
			} else {
				assert.Len(t, tints, len(tt.expectedMoves), "Expected tints for all valid moves")
				for _, move := range tt.expectedMoves {
					tint, exists := tints[move]
					assert.True(t, exists, "Expected tint for position %v", move)
					assert.Equal(t, tt.expectedColor, tint, "Expected %v tint for position %v", tt.expectedColor, move)
				}
			}
		})
	}
}
