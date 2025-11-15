// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"cragspider-go/internal/core"
	"cragspider-go/pkg/animation"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/samber/lo"
)

// SelectedPieceAndPosition represents a selected piece on the board and its position.
type SelectedPieceAndPosition struct {
	Piece    *core.Piece
	Position core.Position
}

type Playfield struct {
	game              *core.Game
	boardLoc          rl.Vector2
	selectedPiece     *SelectedPieceAndPosition
	backgroundSprites *animation.SpriteSheet
	whiteSprites      *animation.SpriteSheet
	blackSprites      *animation.SpriteSheet
}

var _ Scene = (*Playfield)(nil)

// Init initializes the playfield scene with the given width and height.
// It uses the default game configuration loaded from the embedded YAML file.
func (p *Playfield) Init(width, height int) {
	config, err := core.GetConfig()
	if err != nil {
		rl.TraceLog(rl.LogFatal, "error loading default configuration: %v", err)
	}
	p.InitWithConfig(width, height, config)
}

// InitWithConfig initializes the playfield scene with the given width, height, and configuration.
// If config is nil, the default configuration is loaded from the embedded YAML file.
func (p *Playfield) InitWithConfig(width, height int, cfg *core.GameConfig) {
	g, err := core.NewGameWithConfig(cfg)
	if err != nil {
		rl.TraceLog(rl.LogFatal, "error creating game: %v", err)
	}
	p.game = g

	// Calculate board dimensions
	boardWidth := p.game.Board.Columns * core.SquareSize
	boardHeight := p.game.Board.Rows * core.SquareSize

	// Center the board in the window by calculating the upper-left corner position
	boardX := float32(width-boardWidth) / 2
	boardY := float32(height-boardHeight) / 2
	p.boardLoc = rl.Vector2{X: boardX, Y: boardY}

	// Initialize sprite sheets for rendering
	p.backgroundSprites = animation.Load("dungeon_tiles.png", 4, 9)
	p.whiteSprites = animation.Load("adventurer_pieces.png", 6, 18)
	p.blackSprites = animation.Load("monster_pieces.png", 11, 18)
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
		selectedPiece := p.SelectedPiece()
		pieceUnderClick := p.PieceUnderClick(rl.GetMousePosition())
		if selectedPiece == nil {
			p.SelectPiece(pieceUnderClick)
		} else {
			// User is trying to move selected piece to a new location
			dest, err := p.PositionUnderClick(rl.GetMousePosition())
			if err != nil {
				// User clicked outside the board, so deselect.
				p.SelectPiece(pieceUnderClick)
			} else {
				// User is trying to move into a new square.
				move := core.Move{
					dest[0] - selectedPiece.Position[0],
					dest[1] - selectedPiece.Position[1],
				}
				err = p.movePiece(selectedPiece, move)
				if err != nil {
					rl.TraceLog(rl.LogWarning, "failed to move piece %s: %s", selectedPiece.Piece, err)
				}
				p.SelectPiece(nil)
			}
		}
	}
}

// movePiece takes the selected piece and tries to make the specified move. This fails if the location isn't
// a valid one.
func (p *Playfield) movePiece(spp *SelectedPieceAndPosition, move core.Move) error {
	err := p.game.Board.MovePiece(spp.Piece, spp.Position, move)
	return err
}

// update updates the game state since the last time through the gameplay loop.
func (p *Playfield) update() {

}

// render draws the current game state to the screen.
func (p *Playfield) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	if err := p.renderBoard(); err != nil {
		rl.TraceLog(rl.LogError, "error rendering game: %v", err)
	}

	rl.EndDrawing()
}

// Close closes the game and cleans up resources.
func (p *Playfield) Close() {
	if p.backgroundSprites != nil {
		p.backgroundSprites.Unload()
	}
	if p.whiteSprites != nil {
		p.whiteSprites.Unload()
	}
	if p.blackSprites != nil {
		p.blackSprites.Unload()
	}
}

// SelectPiece selects the specified piece, unselecting any previously selected piece.
func (p *Playfield) SelectPiece(piece *core.Piece) {
	// If clicking didn't select a piece, unselect any selected piece
	if piece == nil {
		p.selectedPiece = nil
		return
	}
	// Selecting a selected piece unselects it (toggle)
	if p.selectedPiece != nil && p.selectedPiece.Piece == piece {
		p.selectedPiece = nil
		return
	}
	// Find the position of the piece and store it in the struct
	pos, err := p.game.Board.PieceLocation(piece)
	if err != nil {
		p.selectedPiece = nil
		return
	}
	p.selectedPiece = &SelectedPieceAndPosition{
		Piece:    piece,
		Position: pos,
	}
}

// SelectedPiece returns the currently selected piece and position, or null if there is none.
func (p *Playfield) SelectedPiece() *SelectedPieceAndPosition {
	return lo.Ternary(p.selectedPiece != nil, p.selectedPiece, nil)
}

// PositionUnderClick returns the board position under a mouse click.
// If the user clicked outside the board, then an error is returned.
func (p *Playfield) PositionUnderClick(clickLoc rl.Vector2) (core.Position, error) {
	// Shift the position relative to the board upper corner so the click loc is in board space
	adjClickLoc := rl.Vector2{X: clickLoc.X - p.boardLoc.X, Y: clickLoc.Y - p.boardLoc.Y}

	// Check if the click is outside the board bounds
	if adjClickLoc.X < 0 || adjClickLoc.X >= float32(core.SquareSize*p.game.Board.Columns) ||
		adjClickLoc.Y < 0 || adjClickLoc.Y >= float32(core.SquareSize*p.game.Board.Rows) {
		return core.Position{}, fmt.Errorf("click is outside the board bounds")
	}

	// Just scale the click based on the square size
	return core.Position{int(adjClickLoc.Y / float32(core.SquareSize)), int(adjClickLoc.X / float32(core.SquareSize))}, nil
}

// PieceUnderClick returns the piece under a mouse click.
// If there's no piece there, then nil is returned.
func (p *Playfield) PieceUnderClick(clickLoc rl.Vector2) *core.Piece {
	pos, err := p.PositionUnderClick(clickLoc)
	if err != nil {
		return nil
	}
	return p.game.Board.GetPieceAt(pos)
}

// renderBoard draws the board to the screen with the given board location (where the upper left corner is).
func (p *Playfield) renderBoard() error {
	// If there's a selected piece, figure out its valid moves so we can tint the squares
	var tintedPositions []core.Position
	if p.selectedPiece != nil {
		tintedPositions = p.selectedPiece.Piece.ValidMoves(p.selectedPiece.Position, p.game.Board)
		tintedPositions = append(tintedPositions, p.selectedPiece.Position)
	}
	// First draw the board itself
	for i := range p.game.Board.Rows {
		for j := range p.game.Board.Columns {
			// If there's a valid move on this square, or if it's the currently selected piece, tint it
			tint := lo.Ternary(lo.Contains(tintedPositions, core.Position{i, j}), rl.Green, rl.White)
			err := p.backgroundSprites.DrawFrame(
				p.game.Board.GetSquareAt(core.Position{i, j}).Frame,
				rl.Vector2{X: p.boardLoc.X + float32(j*core.SquareSize), Y: p.boardLoc.Y + float32(i*core.SquareSize)},
				core.Scale,
				p.game.Board.GetSquareAt(core.Position{i, j}).Rotation,
				tint)
			if err != nil {
				return fmt.Errorf("failed to draw cell: %w", err)
			}
		}
	}
	// Now draw each of the pieces on the board
	for i := range p.game.Board.Rows {
		for j := range p.game.Board.Columns {
			piece := p.game.Board.GetPieceAt(core.Position{i, j})
			if piece != nil {
				err2 := p.renderPieceOnBoard(piece, j, i)
				if err2 != nil {
					return err2
				}
			}
		}
	}

	return nil
}

// renderPieceOnBoard renders a single piece on the board at the specified position.
func (p *Playfield) renderPieceOnBoard(piece *core.Piece, j int, i int) error {
	// Different sprite sheets for different players of course
	sheet := lo.Ternary(piece.Color == core.White, p.whiteSprites, p.blackSprites)
	isSelected := p.selectedPiece != nil && p.selectedPiece.Piece == piece
	frame := lo.Ternary(isSelected, 0, 1)
	err := sheet.DrawFrame(
		piece.Config.Sprites[piece.Color][frame],
		rl.Vector2{X: p.boardLoc.X + float32(j*core.SquareSize), Y: p.boardLoc.Y + float32(i*core.SquareSize)},
		core.Scale,
		rl.Vector2{X: 1.0, Y: 0.0},
		rl.White)
	if err != nil {
		return fmt.Errorf("failed to draw piece: %w", err)
	}
	return nil
}
