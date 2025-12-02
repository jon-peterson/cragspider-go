// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	// Create a new board
	cfg, err := GetConfig()
	require.NoError(t, err)
	board, err := newBoard(cfg)
	require.NoError(t, err)

	// Check board dimensions match config
	assert.Equal(t, cfg.Board.Rows, board.Rows, "Board rows should match config")
	assert.Equal(t, cfg.Board.Columns, board.Columns, "Board columns should match config")

	// Check slice dimensions
	require.Len(t, board.squares, cfg.Board.Rows, "Squares should have correct number of rows")
	for _, row := range board.squares {
		assert.Len(t, row, cfg.Board.Columns, "Each cell row should have correct number of columns")
	}
	require.Len(t, board.pieces, cfg.Board.Rows, "Pieces should have correct number of rows")
	for _, row := range board.pieces {
		assert.Len(t, row, cfg.Board.Columns, "Each piece row should have correct number of columns")
	}

	// Check that each cell has a valid rotation vector
	for i, row := range board.squares {
		for j, cell := range row {
			// Check that rotation is one of the cardinal directions
			dir := Move{int(cell.Rotation.X), int(cell.Rotation.Y)}
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
	piece := &Piece{Name: "test_piece", Color: White}
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
			piece:     &Piece{Name: "other_piece", Color: Black},
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
		Rows:     5,
		Columns:  5,
		pieces:   make([][]*Piece, 5),
		captured: make(map[Color][]*Piece),
	}
	for i := range board.pieces {
		board.pieces[i] = make([]*Piece, 5)
	}

	// Create test pieces
	mainPiece := &Piece{
		Name:  "test_piece",
		Color: White,
		Config: PieceConfig{
			Name: "test_piece",
			Moves: [][]Move{
				{{0, 1}},  // right path
				{{0, -1}}, // left path
				{{1, 0}},  // down path
				{{-1, 0}}, // up path
			},
		},
	}
	blockerPiece := &Piece{Name: "blocker", Color: White} // Same color as mainPiece
	wrongPiece := &Piece{Name: "wrong_piece", Color: White}
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
			name:     "move to occupied position (same color)",
			piece:    mainPiece,
			startPos: middlePos,
			move:     Move{1, 0},
			setup: func() {
				// Place a blocker (same color) to the right
				blockerPos := Position{middlePos[0] + 1, middlePos[1]}
				newBoard, err := board.PlacePiece(blockerPiece, blockerPos)
				require.NoError(t, err)
				board = newBoard
			},
			wantErr:    true,
			wantErrMsg: "is not valid",
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
			newBoard, err := board.PlacePiece(mainPiece, tt.startPos)
			require.NoError(t, err, "Failed to place piece")
			board = newBoard

			// Run any test-specific setup
			if tt.setup != nil {
				tt.setup()
			}

			// Execute the move
			resultBoard, err := board.MovePiece(tt.piece, tt.startPos, tt.move)

			// Verify results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, resultBoard)
				}
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
		Name:  "test",
		Color: White,
	}

	t.Run("place piece in empty position", func(t *testing.T) {
		pos := Position{1, 1}
		newBoard, err := board.PlacePiece(piece, pos)
		require.NoError(t, err, "Should be able to place piece in empty position")
		assert.NotNil(t, newBoard.pieces[pos[0]][pos[1]], "Piece should be placed on the board")
		assert.Equal(t, piece.Name, newBoard.pieces[pos[0]][pos[1]].Name, "Placed piece should have the correct name")
		assert.Equal(t, piece.Color, newBoard.pieces[pos[0]][pos[1]].Color, "Placed piece should have the correct color")
	})

	t.Run("cannot place in occupied position", func(t *testing.T) {
		pos := Position{2, 2}
		// First placement should succeed
		newBoard, err := board.PlacePiece(piece, pos)
		require.NoError(t, err, "First placement should succeed")

		// Second placement should fail
		_, err = newBoard.PlacePiece(piece, pos)
		assert.Error(t, err, "Should not be able to place piece in occupied position")
	})

	t.Run("cannot place out of bounds", func(t *testing.T) {
		outOfBoundsPositions := []Position{
			{-1, 0},            // row too small
			{board.Rows, 0},    // row too large
			{0, -1},            // col too small
			{0, board.Columns}, // col too large
		}

		for _, pos := range outOfBoundsPositions {
			_, err := board.PlacePiece(piece, pos)
			assert.Error(t, err, "Should not be able to place piece at position %v", pos)
		}
	})
}

func TestBoard_Capture(t *testing.T) {
	// Create a test board
	board := &Board{
		Rows:     5,
		Columns:  5,
		pieces:   make([][]*Piece, 5),
		captured: make(map[Color][]*Piece),
	}
	for i := range board.pieces {
		board.pieces[i] = make([]*Piece, 5)
	}

	// Create test pieces
	whitePiece := &Piece{
		Name:  "white_warrior",
		Color: White,
		Config: PieceConfig{
			Name: "white_warrior",
			Moves: [][]Move{
				{{1, 0}, {1, 0}},   // down path
				{{-1, 0}, {-1, 0}}, // up path
				{{0, 1}, {0, 1}},   // right path
				{{0, -1}, {0, -1}}, // left path
			},
		},
	}

	blackPiece := &Piece{
		Name:  "black_warrior",
		Color: Black,
	}

	t.Run("capture opponent piece", func(t *testing.T) {
		// Reset board
		for i := range board.pieces {
			board.pieces[i] = make([]*Piece, 5)
		}
		board.captured = make(map[Color][]*Piece)

		// Place white piece at (2, 2) and black piece at (2, 4)
		whitePiece.Config.Moves = [][]Move{{{0, 1}, {0, 1}}, {{0, -1}, {0, -1}}} // right and left paths with 2 steps each
		board.pieces[2][2] = whitePiece
		board.pieces[2][4] = blackPiece

		// Move white piece to capture black piece
		resultBoard, err := board.MovePiece(whitePiece, Position{2, 2}, Move{0, 2})
		require.NoError(t, err, "Should be able to move and capture")

		// Verify white piece moved
		assert.Equal(t, whitePiece, resultBoard.pieces[2][4], "White piece should be at capture position")
		assert.Nil(t, resultBoard.pieces[2][2], "Original position should be empty")

		// Verify black piece was captured
		capturedByWhite := resultBoard.GetCapturedPieces(White)
		assert.Len(t, capturedByWhite, 1, "White should have captured 1 piece")
		assert.Equal(t, blackPiece, capturedByWhite[0], "Captured piece should be the black piece")
	})

	t.Run("multiple captures tracked separately by color", func(t *testing.T) {
		// Reset board
		for i := range board.pieces {
			board.pieces[i] = make([]*Piece, 5)
		}
		board.captured = make(map[Color][]*Piece)

		whitePiece1 := &Piece{Name: "white1", Color: White, Config: PieceConfig{Moves: [][]Move{{{0, 1}, {0, 1}}}}}
		whitePiece2 := &Piece{Name: "white2", Color: White, Config: PieceConfig{Moves: [][]Move{{{0, 1}, {0, 1}}}}}
		blackPiece1 := &Piece{Name: "black1", Color: Black}
		blackPiece2 := &Piece{Name: "black2", Color: Black}

		// Set up positions
		board.pieces[0][0] = whitePiece1
		board.pieces[0][1] = blackPiece1
		board.pieces[1][0] = whitePiece2
		board.pieces[1][1] = blackPiece2

		// White captures black pieces
		resultBoard, err := board.MovePiece(whitePiece1, Position{0, 0}, Move{0, 1})
		require.NoError(t, err)
		board = resultBoard

		resultBoard, err = board.MovePiece(whitePiece2, Position{1, 0}, Move{0, 1})
		require.NoError(t, err)
		board = resultBoard

		// Verify white has 2 captured pieces
		capturedByWhite := board.GetCapturedPieces(White)
		assert.Len(t, capturedByWhite, 2, "White should have 2 captured pieces")
		assert.Equal(t, blackPiece1, capturedByWhite[0])
		assert.Equal(t, blackPiece2, capturedByWhite[1])

		// Verify black has 0 captured pieces
		capturedByBlack := board.GetCapturedPieces(Black)
		assert.Len(t, capturedByBlack, 0, "Black should have 0 captured pieces")
	})

	t.Run("can move to empty square without capturing", func(t *testing.T) {
		// Reset board
		for i := range board.pieces {
			board.pieces[i] = make([]*Piece, 5)
		}
		board.captured = make(map[Color][]*Piece)

		whitePiece := &Piece{Name: "white", Color: White, Config: PieceConfig{Moves: [][]Move{{{1, 0}}, {{0, 1}}}}}
		board.pieces[2][2] = whitePiece

		// Move to empty square
		resultBoard, err := board.MovePiece(whitePiece, Position{2, 2}, Move{1, 0})
		require.NoError(t, err)

		// Verify no captures
		assert.Len(t, resultBoard.GetCapturedPieces(White), 0, "No pieces should be captured")
		assert.Equal(t, whitePiece, resultBoard.pieces[3][2], "Piece should be at new position")
	})
}
