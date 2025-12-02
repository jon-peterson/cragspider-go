// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"testing"

	"cragspider-go/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayfield_PositionUnderClick(t *testing.T) {
	// Create a test playfield with a 10x10 board
	pf := &Playfield{
		boardLoc: rl.Vector2{X: 100, Y: 100},
	}

	// Create a mock game with a board
	game, err := core.NewGame()
	require.NoError(t, err)
	pf.game = game

	tests := []struct {
		name        string
		clickX      float32
		clickY      float32
		expectedRow int
		expectedCol int
		expectErr   bool
	}{
		{
			name:        "top left corner",
			clickX:      100, // Left edge of first column
			clickY:      100, // Top edge of first row
			expectedRow: 0,
			expectedCol: 0,
			expectErr:   false,
		},
		{
			name:        "just inside top left corner",
			clickX:      101,
			clickY:      101,
			expectedRow: 0,
			expectedCol: 0,
			expectErr:   false,
		},
		{
			name:        "middle of first square",
			clickX:      float32(100 + core.SquareSize/2),
			clickY:      float32(100 + core.SquareSize/2),
			expectedRow: 0,
			expectedCol: 0,
			expectErr:   false,
		},
		{
			name:        "on the border between squares",
			clickX:      float32(100 + core.SquareSize), // Exactly on the border between first and second column
			clickY:      float32(100 + core.SquareSize), // Exactly on the border between first and second row
			expectedRow: 1,                              // Should go to the next row/column
			expectedCol: 1,
			expectErr:   false,
		},
		{
			name:      "outside left",
			clickX:    99,  // Left of the board
			clickY:    150, // Within board height
			expectErr: true,
		},
		{
			name:      "outside right",
			clickX:    float32(100 + 10*core.SquareSize),
			clickY:    150,
			expectErr: true,
		},
		{
			name:      "outside top",
			clickX:    150,
			clickY:    99, // Above the board
			expectErr: true,
		},
		{
			name:      "outside bottom",
			clickX:    150,
			clickY:    float32(100 + 10*core.SquareSize),
			expectErr: true,
		},
		{
			name:      "bottom right corner of last square",
			clickX:    float32(100 + 10*core.SquareSize),
			clickY:    float32(100 + 10*core.SquareSize),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clickLoc := rl.Vector2{X: tt.clickX, Y: tt.clickY}
			pos, err := pf.PositionUnderMouse(clickLoc)

			if tt.expectErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Did not expect an error")
				assert.Equal(t, tt.expectedRow, pos[0], "Unexpected row")
				assert.Equal(t, tt.expectedCol, pos[1], "Unexpected column")
			}
		})
	}
}

func TestPlayfield_MouseIsOverBoard(t *testing.T) {
	// Create a test playfield with a 10x10 board
	pf := &Playfield{
		boardLoc: rl.Vector2{X: 100, Y: 100},
	}

	// Create a mock game with a board
	game, err := core.NewGame()
	require.NoError(t, err)
	pf.game = game

	tests := []struct {
		name     string
		mouseX   float32
		mouseY   float32
		expected bool
	}{
		{
			name:     "top left corner",
			mouseX:   100, // Left edge of first column
			mouseY:   100, // Top edge of first row
			expected: true,
		},
		{
			name:     "middle of first square",
			mouseX:   float32(100 + core.SquareSize/2),
			mouseY:   float32(100 + core.SquareSize/2),
			expected: true,
		},
		{
			name:     "just outside left",
			mouseX:   99,  // Left of the board
			mouseY:   150, // Within board height
			expected: false,
		},
		{
			name:     "just outside right",
			mouseX:   float32(100 + 10*core.SquareSize),
			mouseY:   150,
			expected: false,
		},
		{
			name:     "just outside top",
			mouseX:   150,
			mouseY:   99, // Above the board
			expected: false,
		},
		{
			name:     "just outside bottom",
			mouseX:   150,
			mouseY:   float32(100 + 10*core.SquareSize),
			expected: false,
		},
		{
			name:     "bottom right corner of last square",
			mouseX:   float32(100 + 10*core.SquareSize - 1),
			mouseY:   float32(100 + 10*core.SquareSize - 1),
			expected: true,
		},
		{
			name:     "exactly at bottom right corner",
			mouseX:   float32(100 + 10*core.SquareSize),
			mouseY:   float32(100 + 10*core.SquareSize),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mouseLoc := rl.Vector2{X: tt.mouseX, Y: tt.mouseY}
			result := pf.MouseIsOverBoard(mouseLoc)
			assert.Equal(t, tt.expected, result, "Unexpected result for %s", tt.name)
		})
	}
}

func TestPlayfield_SelectPiece(t *testing.T) {
	// Create a test playfield with a custom board
	game, err := core.NewGame()
	require.NoError(t, err)

	pf := &Playfield{
		game: game,
	}

	// Create test pieces
	piece0 := &core.Piece{Name: "piece0", Color: core.White, Config: core.PieceConfig{}}
	piece1 := &core.Piece{Name: "piece1", Color: core.White, Config: core.PieceConfig{}}
	piece2 := &core.Piece{Name: "piece2", Color: core.Black, Config: core.PieceConfig{}}
	newBoard, err := game.Board.PlacePiece(piece0, core.Position{4, 4})
	require.NoError(t, err)
	game.Board = newBoard
	newBoard, err = game.Board.PlacePiece(piece1, core.Position{5, 5})
	require.NoError(t, err)
	game.Board = newBoard
	newBoard, err = game.Board.PlacePiece(piece2, core.Position{6, 6})
	require.NoError(t, err)
	game.Board = newBoard

	tests := []struct {
		name            string
		pieceToSelect   *core.Piece
		expectSelected  bool
		expectedPos     core.Position
		expectNilSelect bool
		selectTwice     bool
	}{
		{
			name:           "select first white piece",
			pieceToSelect:  piece1,
			expectSelected: true,
			expectedPos:    core.Position{5, 5},
		},
		{
			name:           "select another white piece",
			pieceToSelect:  piece0,
			expectSelected: true,
			expectedPos:    core.Position{4, 4},
		},
		{
			name:           "select already selected piece unselects it",
			pieceToSelect:  piece0,
			expectSelected: false,
			selectTwice:    true,
		},
		{
			name:           "select nil (unselect)",
			pieceToSelect:  nil,
			expectSelected: false,
		},
		{
			name:            "cannot select black piece during white turn",
			pieceToSelect:   piece2,
			expectSelected:  false,
			expectNilSelect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset selected piece for each test
			pf.selectedPiece = nil

			// For the toggle test, we need to select the piece first, then select it again
			if tt.selectTwice {
				pf.SelectPiece(tt.pieceToSelect)
			}

			// Select the piece
			pf.SelectPiece(tt.pieceToSelect)

			if tt.expectNilSelect {
				assert.Nil(t, pf.selectedPiece)
				return
			}

			if tt.expectSelected {
				require.NotNil(t, pf.selectedPiece, "Expected a piece to be selected")
				assert.Equal(t, tt.pieceToSelect, pf.selectedPiece.Piece)
				assert.Equal(t, tt.expectedPos, pf.selectedPiece.Position)
			} else {
				assert.Nil(t, pf.selectedPiece, "Expected no piece to be selected")
			}
		})
	}
}
