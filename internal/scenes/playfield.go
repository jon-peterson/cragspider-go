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

	// User click is used to select a piece, unselect a piece, or move a piece depending
	// on the current state of the board.
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		selectedPiece := p.game.Board.SelectedPiece()
		pieceUnderClick := p.game.Board.PieceUnderClick(p.boardLoc, rl.GetMousePosition())
		if selectedPiece == nil {
			// User is trying to select a piece.
			p.game.Board.SelectPiece(pieceUnderClick)
		} else if pieceUnderClick == nil {
			// User is moving the selected piece to an empty square.
			dest, err := p.game.Board.PositionUnderClick(p.boardLoc, rl.GetMousePosition())
			if err != nil {
				// User clicked outside the board, so deselect.
				p.game.Board.SelectPiece(pieceUnderClick)
			} else {
				// User is trying to move into a new square.
				move := core.Move{
					dest[0] - selectedPiece.Position[0],
					dest[1] - selectedPiece.Position[1],
				}
				err := p.game.Board.MovePiece(selectedPiece.Piece, selectedPiece.Position, move)
				if err != nil {
					rl.TraceLog(rl.LogWarning, "failed to move piece %s: %s", selectedPiece.Piece, err)
				}
				// Regardless of whether move was allowed, toggle the selection.
				p.game.Board.SelectPiece(pieceUnderClick)
			}
		}
	}
}

// update updates the game state since the last time through the gameplay loop.
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
