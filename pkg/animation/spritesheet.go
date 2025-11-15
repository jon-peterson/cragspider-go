// Copyright 2025 Ideograph LLC. All rights reserved.

package animation

import (
	"fmt"
	"image/color"
	"math"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SpriteSheet struct {
	name        string
	texture     rl.Texture2D // The texture with the packed sprites
	frameWidth  int          // Width of each frame in pixels
	frameHeight int          // Height of each frame pixels
	rows        int          // the number of rows in the spritesheet
	cols        int          // the number of columns in the spritesheet
	origin      rl.Vector2   // The middle of the sprite (for rotation)
}

// FrameCoords represents a location in the spritesheet by row and column.
type FrameCoords [2]int

// spriteSheetMap manages a named map of SpriteSheets with locking.
type spriteSheetMap struct {
	spritesMap map[string]*SpriteSheet
	mapLock    sync.RWMutex
}

var sheetMap *spriteSheetMap = newSheetMap()

// newSheetMap creates a new spriteSheetMap with an empty map.
func newSheetMap() *spriteSheetMap {
	return &spriteSheetMap{
		spritesMap: make(map[string]*SpriteSheet),
	}
}

// Load creates a new spritesheet from the given file, with deferred loading of
// the actual texture. These are cached for performance.
func Load(file string, rows, cols int) *SpriteSheet {
	sheetMap.mapLock.RLock()
	if sprite, ok := sheetMap.spritesMap[file]; ok {
		sheetMap.mapLock.RUnlock()
		return sprite
	}
	// First time being asked for this spritesheet, so we need to create cache and return
	sheetMap.mapLock.RUnlock()
	sheetMap.mapLock.Lock()
	defer sheetMap.mapLock.Unlock()
	s := SpriteSheet{
		name: file,
		rows: rows,
		cols: cols,
	}
	sheetMap.spritesMap[file] = &s
	return &s
}

// Unload unloads the sprite sheet from memory.
func (s *SpriteSheet) Unload() {
	sheetMap.mapLock.Lock()
	defer sheetMap.mapLock.Unlock()
	delete(sheetMap.spritesMap, s.name)
	rl.UnloadTexture(s.texture)
}

// populateTexture loads a spritesheet from the saved filename, or returns an error if it can't.
// It initializes it with the specified rows and columns. SpriteSheets are cached.
func (s *SpriteSheet) populateTexture() error {
	sheetMap.mapLock.Lock()
	defer sheetMap.mapLock.Unlock()
	if s.frameWidth != 0 {
		// Already loaded; just return
		return nil
	}
	sheetTexture := rl.LoadTexture("assets/sprites/" + s.name)
	if int(sheetTexture.Width)%s.cols != 0 || int(sheetTexture.Height)%s.rows != 0 {
		return fmt.Errorf("spritesheet of dimensions (%d,%d) can't be broken into %d rows and %d cols",
			sheetTexture.Width, sheetTexture.Height, s.rows, s.cols)
	}
	s.texture = sheetTexture
	s.frameWidth = int(sheetTexture.Width) / s.cols
	s.frameHeight = int(sheetTexture.Height) / s.rows
	s.origin = rl.NewVector2(float32(s.frameWidth)/2, float32(s.frameHeight)/2)
	sheetMap.spritesMap[s.name] = s
	return nil
}

// String returns a string representation of the spritesheet.
func (s *SpriteSheet) String() string {
	return fmt.Sprintf("%s (%dx%d)", s.name, s.frameWidth, s.frameHeight)
}

// DrawFrame draws the sprite at the given frame at the given location, scale, and rotation. The location is the upper
// left hand corner of the sprite, but rotation goes around the center of the sprite.
func (s *SpriteSheet) DrawFrame(frameCoords FrameCoords, loc rl.Vector2, scale float32, rot rl.Vector2, tint color.RGBA) error {
	if s.frameWidth == 0 {
		// Texture hasn't been loaded yet, so load it now
		if err := s.populateTexture(); err != nil {
			return err
		}
	}
	frame, err := s.frame(frameCoords[0], frameCoords[1])
	if err != nil {
		return err
	}
	width := float32(s.frameWidth) * scale
	height := float32(s.frameHeight) * scale
	destination := rl.Rectangle{
		X:      loc.X + (width / 2),
		Y:      loc.Y + (height / 2),
		Width:  width,
		Height: height,
	}
	rotationDegrees := float32(math.Atan2(float64(rot.Y), float64(rot.X)) * 180 / math.Pi)
	rl.DrawTexturePro(s.texture, frame, destination, rl.Vector2Scale(s.origin, scale), rotationDegrees, tint)
	return nil
}

// frame returns the rectangle for the given frame in the spritesheet.
func (s *SpriteSheet) frame(row, col int) (rl.Rectangle, error) {
	if row < 0 || row >= int(s.rows) || col < 0 || col >= int(s.cols) {
		return rl.Rectangle{}, fmt.Errorf("frame (%d,%d) is out of bounds", row, col)
	}
	return rl.Rectangle{
		X:      float32(col) * float32(s.frameWidth),
		Y:      float32(row) * float32(s.frameHeight),
		Width:  float32(s.frameWidth),
		Height: float32(s.frameHeight),
	}, nil
}

// FrameLocation returns the coordinates for the frame assuming row-first ordering.
func (s *SpriteSheet) FrameLocation(f int) (FrameCoords, error) {
	if f < 0 || f >= int(s.rows*s.cols) {
		return FrameCoords{}, fmt.Errorf("frame %d is out of bounds", f)
	}
	row := f / int(s.cols)
	col := f % int(s.cols)
	return [2]int{row, col}, nil
}

// GetSize returns the size of the spritesheet in pixels as a vector.
func (s *SpriteSheet) GetSize() rl.Vector2 {
	return rl.Vector2{
		X: float32(s.frameWidth),
		Y: float32(s.frameHeight),
	}
}

// GetRectangle returns the bounding rectangle where this sprite will be drawn centered at center.
func (s *SpriteSheet) GetRectangle(center rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      center.X - float32(s.frameWidth)/2,
		Y:      center.Y - float32(s.frameHeight)/2,
		Width:  float32(s.frameWidth),
		Height: float32(s.frameHeight),
	}
}
