// Copyright 2025 Ideograph LLC. All rights reserved.

package animation

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	rl.InitWindow(800, 600, "Test")
	defer rl.CloseWindow()
	m.Run()
}

func TestSpriteSheet_frame(t *testing.T) {
	sheet := Load("test_sheet.png", 4, 7)
	defer sheet.Unload()

	tests := []struct {
		name     string
		row      int
		col      int
		expected rl.Rectangle
		expErr   bool
	}{
		{
			name:     "top-left corner",
			row:      0,
			col:      0,
			expected: rl.Rectangle{X: 0, Y: 0, Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)},
			expErr:   false,
		},
		{
			name:     "top-right corner",
			row:      0,
			col:      6,
			expected: rl.Rectangle{X: float32(sheet.frameWidth), Y: 0, Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)},
			expErr:   false,
		},
		{
			name:     "bottom-left corner",
			row:      3,
			col:      0,
			expected: rl.Rectangle{X: 0, Y: float32(sheet.frameHeight), Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)},
			expErr:   false,
		},
		{
			name:     "bottom-right corner",
			row:      3,
			col:      6,
			expected: rl.Rectangle{X: float32(sheet.frameWidth), Y: float32(sheet.frameHeight), Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)},
			expErr:   false,
		},
		{
			name:   "negative row",
			row:    -1,
			col:    0,
			expErr: true,
		},
		{
			name:   "negative column",
			row:    0,
			col:    -1,
			expErr: true,
		},
		{
			name:   "row out of bounds",
			row:    4,
			col:    0,
			expErr: true,
		},
		{
			name:   "column out of bounds",
			row:    0,
			col:    7,
			expErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := sheet.frame(tt.row, tt.col)
			if tt.expErr {
				assert.Error(t, err, "Expected error for %s", tt.name)
				return
			}
			require.NoError(t, err, "Unexpected error for %s", tt.name)
			assert.Equal(t, tt.expected, frame, "Unexpected frame for %s", tt.name)
		})
	}
}

func TestSpriteSheet_GetRectangle(t *testing.T) {
	sheet := Load("test_sheet.png", 4, 7)
	defer sheet.Unload()

	center := rl.NewVector2(50, 50)
	expected := rl.Rectangle{
		X:      center.X - float32(sheet.frameWidth)/2,
		Y:      center.Y - float32(sheet.frameHeight)/2,
		Width:  float32(sheet.frameWidth),
		Height: float32(sheet.frameHeight),
	}

	rect := sheet.GetRectangle(center)
	assert.Equal(t, expected, rect, "Rectangle should be centered on the given point")
}

func TestSpriteSheet_GetSize(t *testing.T) {
	sheet := Load("test_sheet.png", 4, 7)
	defer sheet.Unload()

	expected := rl.Vector2{X: float32(sheet.frameWidth), Y: float32(sheet.frameHeight)}
	assert.Equal(t, expected, sheet.GetSize(), "Size should match frame dimensions")
}

func TestLoadSpriteSheet_Cache(t *testing.T) {
	// Load the sprite sheet for the first time
	sheet1 := Load("test_sheet.png", 2, 2)
	sheet2 := Load("test_sheet.png", 2, 2)
	defer sheet1.Unload()
	defer sheet2.Unload()

	// Verify that the same pointer is returned
	samePointer := assert.Same(t, sheet1, sheet2, "Should return the same sprite sheet instance from cache")
	if !samePointer {
		t.Log("Cache test failed - different instances returned for same sprite sheet")
	}
}
