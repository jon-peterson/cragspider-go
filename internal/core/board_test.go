package core

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	// Create a new board
	board := newBoard()

	// Check board dimensions
	if board.Rows != 10 {
		t.Errorf("Expected board.Rows to be 10, got %d", board.Rows)
	}
	if board.Columns != 10 {
		t.Errorf("Expected board.Columns to be 10, got %d", board.Columns)
	}

	// Check cells slice dimensions
	if len(board.cells) != 10 {
		t.Fatalf("Expected 10 rows in cells, got %d", len(board.cells))
	}
	for i, row := range board.cells {
		if len(row) != 10 {
			t.Errorf("Expected 10 columns in row %d, got %d", i, len(row))
		}
	}

	// Check that backgroundSprites is initialized
	if board.backgroundSprites == nil {
		t.Error("Expected backgroundSprites to be initialized, got nil")
	}

	// Check that each cell has a valid rotation vector
	for i, row := range board.cells {
		for j, cell := range row {
			// Check that rotation is one of the cardinal directions
			valid := false
			for _, dir := range CardinalDirections {
				if cell.rotation == dir {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("Cell at [%d][%d] has invalid rotation vector: %v", i, j, cell.rotation)
			}
		}
	}
}
