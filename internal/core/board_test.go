// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	// Create a new board
	cfg, err := GetConfig()
	require.NoError(t, err)
	board, err := newBoard(cfg)
	require.NoError(t, err)

	// Check board dimensions
	assert.Equal(t, 10, board.Rows, "Board should have 10 rows")
	assert.Equal(t, 10, board.Columns, "Board should have 10 columns")

	// Check slice dimensions
	require.Len(t, board.squares, 10, "Squares should have 10 rows")
	for _, row := range board.squares {
		assert.Len(t, row, 10, "Each cell row should have 10 columns")
	}
	require.Len(t, board.pieces, 10, "Pieces should have 10 rows")
	for _, row := range board.pieces {
		assert.Len(t, row, 10, "Each piece row should have 10 columns")
	}

	// Check that backgroundSprites is initialized
	assert.NotNil(t, board.backgroundSprites, "Background sprites should be initialized")

	// Check that each cell has a valid rotation vector
	for i, row := range board.squares {
		for j, cell := range row {
			// Check that rotation is one of the cardinal directions
			dir := Move{int(cell.rotation.X), int(cell.rotation.Y)}
			assert.Contains(t, CardinalDirections, dir,
				"Square at [%d][%d] has invalid rotation vector: %v", i, j, dir)
		}
	}

	// Verify that each piece in the game config is on that position in the board
	whitePieces, err := cfg.Board.GetStartingPositions(White)
	require.NoError(t, err)
	for _, white := range whitePieces {
		pieceOnBoard := board.pieces[white.Position[0]][white.Position[1]]
		assert.NotNil(t, pieceOnBoard, "white piece %s should be at position %v", white.Name, white.Position)
	}
	// Do the same for the black pieces
	blackPieces, err := cfg.Board.GetStartingPositions(Black)
	require.NoError(t, err)
	for _, black := range blackPieces {
		pieceOnBoard := board.pieces[black.Position[0]][black.Position[1]]
		assert.NotNil(t, pieceOnBoard, "black piece %s should be at position %v", black.Name, black.Position)
	}
}

func TestBoard_PositionUnderClick(t *testing.T) {
	// Create a test board with 10x10 grid
	board := &Board{
		Rows:    10,
		Columns: 10,
	}

	// Board's top-left corner at (100, 100) in screen coordinates
	boardLoc := rl.Vector2{X: 100, Y: 100}

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
			clickX:      float32(100 + SquareSize/2),
			clickY:      float32(100 + SquareSize/2),
			expectedRow: 0,
			expectedCol: 0,
			expectErr:   false,
		},
		{
			name:        "on the border between squares",
			clickX:      float32(100 + SquareSize), // Exactly on the border between first and second column
			clickY:      float32(100 + SquareSize), // Exactly on the border between first and second row
			expectedRow: 1,                         // Should go to the next row/column
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
			clickX:    float32(100 + 10*SquareSize),
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
			clickY:    float32(100 + 10*SquareSize),
			expectErr: true,
		},
		{
			name:      "bottom right corner of last square",
			clickX:    float32(100 + 10*SquareSize),
			clickY:    float32(100 + 10*SquareSize),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clickLoc := rl.Vector2{X: tt.clickX, Y: tt.clickY}
			pos, err := board.PositionUnderClick(boardLoc, clickLoc)

			if tt.expectErr {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Did not expect an error")
				assert.Equal(t, tt.expectedRow, pos[0], "Unexpected row")
				assert.Equal(t, tt.expectedCol, pos[1], "Unexpected column")
			}
		})
	}

	t.Run("non-square board", func(t *testing.T) {
		// Test with a non-square board (5 rows, 10 columns)
		nonSquareBoard := &Board{
			Rows:    5,
			Columns: 10,
		}

		tests := []struct {
			name        string
			clickX      float32
			clickY      float32
			expectedRow int
			expectedCol int
			expectErr   bool
		}{
			{
				name:        "top left of non-square board",
				clickX:      50,
				clickY:      50,
				expectedRow: 0,
				expectedCol: 0,
				expectErr:   false,
			},
			{
				name:        "bottom right of non-square board",
				clickX:      float32(50 + 10*SquareSize - 1), // Last column
				clickY:      float32(50 + 5*SquareSize - 1),  // Last row
				expectedRow: 4,
				expectedCol: 9,
				expectErr:   false,
			},
			{
				name:      "outside right edge of non-square board",
				clickX:    float32(50 + 10*SquareSize + 1), // On the right edge (invalid)
				clickY:    float32(50 + 2*SquareSize),      // Somewhere in the middle vertically
				expectErr: true,
			},
			{
				name:      "outside bottom edge of non-square board",
				clickX:    float32(50 + 5*SquareSize),     // Somewhere in the middle horizontally
				clickY:    float32(50 + 5*SquareSize + 1), // On the bottom edge (invalid)
				expectErr: true,
			},
		}

		boardLoc := rl.Vector2{X: 50, Y: 50}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				clickLoc := rl.Vector2{X: tt.clickX, Y: tt.clickY}
				pos, err := nonSquareBoard.PositionUnderClick(boardLoc, clickLoc)

				if tt.expectErr {
					assert.Error(t, err, "Expected an error")
				} else {
					assert.NoError(t, err, "Did not expect an error")
					assert.Equal(t, tt.expectedRow, pos[0], "Unexpected row")
					assert.Equal(t, tt.expectedCol, pos[1], "Unexpected column")
				}
			})
		}
	})
}

func TestBoard_PieceLocation(t *testing.T) {
	// Create a test board
	board := &Board{
		Rows:    3,
		Columns: 3,
		pieces:  make([][]*Piece, 3),
	}
	for i := range board.pieces {
		board.pieces[i] = make([]*Piece, 3)
	}

	// Create a test piece and place it at (1,1)
	piece := &Piece{name: "test_piece", color: White}
	board.pieces[1][1] = piece

	tests := []struct {
		name        string
		piece       *Piece
		expectedPos Position
		expectErr   bool
		errMsg      string
	}{
		{
			name:        "find existing piece",
			piece:       piece,
			expectedPos: Position{1, 1},
			expectErr:   false,
		},
		{
			name:      "piece not on board",
			piece:     &Piece{name: "other_piece", color: Black},
			expectErr: true,
			errMsg:    "not found on board",
		},
		{
			name:      "nil piece",
			piece:     nil,
			expectErr: true,
			errMsg:    "cannot find location of nil piece",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.PieceLocation(tt.piece)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPos, pos)
			}
		})
	}
}

func TestBoard_MovePiece(t *testing.T) {
	// Create a test board
	board := &Board{
		Rows:    5,
		Columns: 5,
		pieces:  make([][]*Piece, 5),
	}
	for i := range board.pieces {
		board.pieces[i] = make([]*Piece, 5)
	}

	// Create test pieces
	mainPiece := &Piece{
		name:  "test_piece",
		color: White,
		config: PieceConfig{
			Name: "test_piece",
			Moves: []Move{
				{0, 1},  // up
				{1, 0},  // right
				{0, -1}, // down
				{-1, 0}, // left
			},
		},
	}
	blockerPiece := &Piece{name: "blocker", color: Black}
	wrongPiece := &Piece{name: "wrong_piece", color: White}
	middlePos := Position{2, 2}

	tests := []struct {
		name       string
		piece      *Piece
		startPos   Position
		move       Move
		setup      func() // Optional setup function
		wantErr    bool
		wantErrMsg string
		verify     func(t *testing.T, b *Board)
	}{
		{
			name:     "valid move right",
			piece:    mainPiece,
			startPos: middlePos,
			move:     Move{1, 0},
			wantErr:  false,
			verify: func(t *testing.T, b *Board) {
				// Verify the piece was moved
				assert.Nil(t, b.pieces[middlePos[0]][middlePos[1]], "Original position should be empty")
				expectedPos := Position{middlePos[0] + 1, middlePos[1]}
				assert.Equal(t, mainPiece, b.pieces[expectedPos[0]][expectedPos[1]], "Piece should be at new position")
			},
		},
		{
			name:       "invalid diagonal move",
			piece:      mainPiece,
			startPos:   middlePos,
			move:       Move{1, 1},
			wantErr:    true,
			wantErrMsg: "is not valid",
		},
		{
			name:       "move off board",
			piece:      mainPiece,
			startPos:   Position{0, 0},
			move:       Move{0, -1},
			wantErr:    true,
			wantErrMsg: "is not valid",
		},
		{
			name:     "move to occupied position",
			piece:    mainPiece,
			startPos: middlePos,
			move:     Move{1, 0},
			setup: func() {
				// Place a blocker to the right
				blockerPos := Position{middlePos[0] + 1, middlePos[1]}
				board.PlacePiece(blockerPiece, blockerPos)
			},
			wantErr:    true,
			wantErrMsg: "cannot move",
		},
		{
			name:     "move from empty position",
			piece:    mainPiece,
			startPos: Position{4, 4},
			move:     Move{0, 1},
			setup: func() {
				board.pieces[4][4] = nil
			},
			wantErr:    true,
			wantErrMsg: "is not at [4,4]",
		},
		{
			name:       "move wrong piece",
			piece:      wrongPiece,
			startPos:   middlePos,
			move:       Move{0, 1},
			wantErr:    true,
			wantErrMsg: "is not at [2,2]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the board state
			for i := range board.pieces {
				board.pieces[i] = make([]*Piece, 5)
			}
			// Put the piece in the right place
			err := board.PlacePiece(mainPiece, tt.startPos)
			require.NoError(t, err, "Failed to place piece")

			// Run any test-specific setup
			if tt.setup != nil {
				tt.setup()
			}

			// Execute the move
			err = board.MovePiece(tt.piece, tt.startPos, tt.move)

			// Verify results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, board)
				}
			}
		})
	}
}

func TestBoard_SelectPiece(t *testing.T) {
	// Create a test board
	board := &Board{
		Rows:    3,
		Columns: 3,
		pieces:  make([][]*Piece, 3),
	}
	for i := range board.pieces {
		board.pieces[i] = make([]*Piece, 3)
	}

	// Create test pieces
	piece0 := &Piece{name: "piece0", color: White}
	piece1 := &Piece{name: "piece1", color: White}
	piece2 := &Piece{name: "piece2", color: Black}
	board.pieces[0][0] = piece0
	board.pieces[1][1] = piece1
	board.pieces[2][2] = piece2

	tests := []struct {
		name            string
		pieceToSelect   *Piece
		expectSelected  bool
		expectedPos     Position
		expectNilSelect bool
	}{
		{
			name:           "select first piece",
			pieceToSelect:  piece1,
			expectSelected: true,
			expectedPos:    Position{1, 1},
		},
		{
			name:           "select second piece",
			pieceToSelect:  piece2,
			expectSelected: true,
			expectedPos:    Position{2, 2},
		},
		{
			name:           "select already selected piece unselects it",
			pieceToSelect:  piece0,
			expectSelected: false, // First selection will be true, then we'll select again
		},
		{
			name:           "select nil (unselect)",
			pieceToSelect:  nil,
			expectSelected: false,
		},
		{
			name:            "select non-existent piece",
			pieceToSelect:   &Piece{name: "non-existent", color: White},
			expectSelected:  false,
			expectNilSelect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start tests with piece 0 already selected
			board.SelectPiece(piece0)

			// Select the new piece
			board.SelectPiece(tt.pieceToSelect)

			if tt.expectNilSelect {
				assert.Nil(t, board.selectedPiece)
				return
			}

			if tt.expectSelected {
				require.NotNil(t, board.selectedPiece, "Expected a piece to be selected")
				assert.Equal(t, tt.pieceToSelect, board.selectedPiece.Piece)
				assert.Equal(t, tt.expectedPos, board.selectedPiece.Position)
			} else {
				assert.Nil(t, board.selectedPiece, "Expected no piece to be selected")
			}
		})
	}
}

func TestBoard_PlacePiece(t *testing.T) {
	cfg, err := GetConfig()
	require.NoError(t, err)
	board, err := newBoard(cfg)
	require.NoError(t, err)

	// Create a test piece
	piece := &Piece{
		name:  "test",
		color: White,
	}

	t.Run("place piece in empty position", func(t *testing.T) {
		pos := Position{1, 1}
		err := board.PlacePiece(piece, pos)
		require.NoError(t, err, "Should be able to place piece in empty position")
		assert.NotNil(t, board.pieces[pos[0]][pos[1]], "Piece should be placed on the board")
		assert.Equal(t, piece.name, board.pieces[pos[0]][pos[1]].name, "Placed piece should have the correct name")
		assert.Equal(t, piece.color, board.pieces[pos[0]][pos[1]].color, "Placed piece should have the correct color")
	})

	t.Run("cannot place in occupied position", func(t *testing.T) {
		pos := Position{2, 2}
		// First placement should succeed
		err := board.PlacePiece(piece, pos)
		require.NoError(t, err, "First placement should succeed")

		// Second placement should fail
		err = board.PlacePiece(piece, pos)
		assert.Error(t, err, "Should not be able to place piece in occupied position")
	})

	t.Run("cannot place out of bounds", func(t *testing.T) {
		outOfBoundsPositions := []Position{
			{-1, 0}, // row too small
			{10, 0}, // row too large
			{0, -1}, // col too small
			{0, 10}, // col too large
		}

		for _, pos := range outOfBoundsPositions {
			err := board.PlacePiece(piece, pos)
			assert.Error(t, err, "Should not be able to place piece at position %v", pos)
		}
	})
}
