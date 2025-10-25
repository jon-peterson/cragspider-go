// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"cragspider-go/internal/core"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Playfield struct {
	game     *core.Game
	boardLoc rl.Vector2
}

var _ Scene = (*Playfield)(nil)

// Init initializes the playfield scene with the given width and height.
func (p *Playfield) Init(width, height int) {
	g, err := core.NewGame()
	if err != nil {
		rl.TraceLog(rl.LogFatal, "error creating game: %v", err)
	}
	p.boardLoc = rl.Vector2{X: 48, Y: 24}
	p.game = g
}

// Loop is the basic gameplay loop. Returns a scene code to indicate the next scene.
func (p *Playfield) Loop() SceneCode {
	for !rl.WindowShouldClose() && !p.game.Over() {
		p.handleInput()
		p.update()
		p.render()
	}
	return Quit
}

// handleInput processes keyboard and mouse input.
func (p *Playfield) handleInput() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		piece := p.game.Board.PieceUnderClick(p.boardLoc, rl.GetMousePosition())
		if piece != nil {
			piece.ToggleSelected()
		}
	}
}

// update updates the game state since the last time throuh the gameplay loop.
func (p *Playfield) update() {

}

// render draws the current game state to the screen.
func (p *Playfield) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	if err := p.game.Board.Render(p.boardLoc); err != nil {
		rl.TraceLog(rl.LogError, "error rendering game: %v", err)
	}

	rl.EndDrawing()
}

// Close closes the game and cleans up resources.
func (p *Playfield) Close() {
	// no op until there are some resources
}
