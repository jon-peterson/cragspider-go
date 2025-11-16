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
		game:     game,
		boardLoc: rl.Vector2{X: 0, Y: 0}, // Simple location for test
	}

	tests := []struct {
		name          string
		selectedPiece *core.Piece
		selectedPos   core.Position
		mousePos      rl.Vector2
		expectedColor rl.Color
		expectedMoves []core.Position
		checkColor    bool // Whether to check color when positions don't match exactly
	}{
		{
			name:          "no selected piece, mouse not over board",
			selectedPiece: nil,
			selectedPos:   core.Position{},
			mousePos:      rl.Vector2{X: -100, Y: -100}, // Off the board
			expectedColor: rl.Color{},
			expectedMoves: []core.Position{},
			checkColor:    false,
		},
		{
			name: "with selected white piece, mouse not over board",
			selectedPiece: &core.Piece{
				Name:  "test",
				Color: core.White,
				Config: core.PieceConfig{
					Moves: []core.Move{{-1, 0}, {1, 0}},
				},
			},
			selectedPos:   core.Position{3, 3},
			mousePos:      rl.Vector2{X: -100, Y: -100}, // Off the board
			expectedColor: rl.Green,
			expectedMoves: []core.Position{{3, 3}, {2, 3}, {4, 3}},
			checkColor:    true,
		},
		{
			name: "with selected black piece, mouse not over board",
			selectedPiece: &core.Piece{
				Name:  "test",
				Color: core.Black,
				Config: core.PieceConfig{
					Moves: []core.Move{{-1, 0}, {1, 0}},
				},
			},
			selectedPos:   core.Position{3, 3},
			mousePos:      rl.Vector2{X: -100, Y: -100}, // Off the board
			expectedColor: rl.Red,
			expectedMoves: []core.Position{{3, 3}, {2, 3}, {4, 3}},
			checkColor:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.selectedPiece == nil {
				pf.selectedPiece = nil
			} else {
				pf.selectedPiece = &SelectedPieceAndPosition{
					Piece:    tt.selectedPiece,
					Position: tt.selectedPos,
				}
			}

			tints := pf.getTintedPositions(tt.mousePos)

			if len(tt.expectedMoves) == 0 {
				// No selected piece and mouse not over board
				assert.Empty(t, tints, "Expected no tints")
			} else {
				// We have explicit expected moves (for selected piece cases)
				assert.Len(t, tints, len(tt.expectedMoves), "Expected tints for all moves")
				for _, move := range tt.expectedMoves {
					tint, exists := tints[move]
					assert.True(t, exists, "Expected tint for position %v", move)
					if tt.checkColor {
						assert.Equal(t, tt.expectedColor, tint, "Expected %v tint for position %v", tt.expectedColor, move)
					}
				}
			}
		})
	}

	// Test mouse hovering over actual board pieces
	t.Run("mouse hover over white piece", func(t *testing.T) {
		pf.selectedPiece = nil // Clear selection
		// Find a white piece on the board
		var whitePiece *core.Piece
		var whitePiecePos core.Position
		for i := range pf.game.Board.Rows {
			for j := range pf.game.Board.Columns {
				p := pf.game.Board.GetPieceAt(core.Position{i, j})
				if p != nil && p.Color == core.White {
					whitePiece = p
					whitePiecePos = core.Position{i, j}
					break
				}
			}
			if whitePiece != nil {
				break
			}
		}
		assert.NotNil(t, whitePiece, "Should find at least one white piece on board")

		// Calculate mouse position at center of that square
		mouseX := float32(whitePiecePos[1]*core.SquareSize + core.SquareSize/2)
		mouseY := float32(whitePiecePos[0]*core.SquareSize + core.SquareSize/2)

		tints := pf.getTintedPositions(rl.Vector2{X: mouseX, Y: mouseY})

		// Just verify we have some tints, don't check specific colors
		assert.NotEmpty(t, tints, "Should have tints when mouse is over piece")
	})
}
