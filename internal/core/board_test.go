package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	// Create a new board
	board := newBoard()

	// Check board dimensions
	assert.Equal(t, 10, board.Rows, "Board should have 10 rows")
	assert.Equal(t, 10, board.Columns, "Board should have 10 columns")

	// Check cells slice dimensions
	require.Len(t, board.cells, 10, "Cells should have 10 rows")
	for _, row := range board.cells {
		assert.Len(t, row, 10, "Each row should have 10 columns")
	}

	// Check that backgroundSprites is initialized
	assert.NotNil(t, board.backgroundSprites, "Background sprites should be initialized")

	// Check that each cell has a valid rotation vector
	for i, row := range board.cells {
		for j, cell := range row {
			// Check that rotation is one of the cardinal directions
			assert.Contains(t, CardinalDirections, cell.rotation,
				"Cell at [%d][%d] has invalid rotation vector: %v", i, j, cell.rotation)
		}
	}
}
