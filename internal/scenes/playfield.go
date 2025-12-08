// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"cragspider-go/internal/ai"
	"cragspider-go/internal/core"
	"cragspider-go/pkg/graphics"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
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
	backgroundSprites *graphics.SpriteSheet
	whiteSprites      *graphics.SpriteSheet
	blackSprites      *graphics.SpriteSheet
	moveExecutionChan chan struct {
		piece    *core.Piece
		position core.Position
		move     core.Move
	}
	planningMove bool
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
	whitePlayer := core.NewHumanPlayer()
	blackPlayer := core.NewAIPlayer("Random AI", ai.NewRandomBot(core.Black))

	g, err := core.NewGameWithConfigAndPlayers(cfg, whitePlayer, blackPlayer)
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
	p.backgroundSprites = graphics.Load("dungeon_tiles.png", 4, 9)
	p.whiteSprites = graphics.Load("adventurer_pieces.png", 6, 18)
	p.blackSprites = graphics.Load("monster_pieces.png", 11, 18)

	// Initialize the channel for AI move execution
	p.moveExecutionChan = make(chan struct {
		piece    *core.Piece
		position core.Position
		move     core.Move
	}, 1)
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
		pieceUnderClick := p.PieceUnderMouse(rl.GetMousePosition())
		if p.selectedPiece == nil {
			p.SelectPiece(pieceUnderClick)
		} else {
			// User is trying to move selected piece to a new location
			dest, err := p.PositionUnderMouse(rl.GetMousePosition())
			if err != nil {
				// User clicked outside the board, so deselect.
				p.SelectPiece(pieceUnderClick)
			} else {
				// User is trying to move into a new square.
				move := core.Move{
					dest[0] - p.selectedPiece.Position[0],
					dest[1] - p.selectedPiece.Position[1],
				}
				err = p.movePiece(p.selectedPiece, move)
				if err != nil {
					rl.TraceLog(rl.LogWarning, "failed to move piece %s: %s", p.selectedPiece.Piece, err)
				}
				p.SelectPiece(nil)
			}
		}
	}
}

// movePiece takes the selected piece and tries to make the specified move. This fails if the location isn't
// a valid one. If the move succeeds, the turn is advanced to the next player.
func (p *Playfield) movePiece(spp *SelectedPieceAndPosition, move core.Move) error {
	newBoard, err := p.game.Board.MovePiece(spp.Piece, spp.Position, move)
	if err == nil {
		p.game.Board = newBoard
		p.game.AdvanceTurn()
	}
	return err
}

// update updates the game state since the last time through the gameplay loop.
// If the current player is AI controlled, executes their move automatically.
func (p *Playfield) update() {
	currentPlayer := p.game.GetPlayer(p.game.ActiveColor)
	if !currentPlayer.IsAI() {
		return
	}

	// First, check if there's a pending move that should be executed
	select {
	case execution := <-p.moveExecutionChan:
		// Execute the pending move
		p.movePiece(&SelectedPieceAndPosition{Piece: execution.piece, Position: execution.position}, execution.move)
		p.SelectPiece(nil)
		p.planningMove = false
		return
	default:
		// No pending move; continue with planning the next AI move
	}

	// If we're not already planning a move, we should kick that off now
	if !p.planningMove {
		p.planningMove = true
		go p.planAIMove(currentPlayer)
	}
}

// planAIMove plans and executes an AI move with visualization.
// It runs in a goroutine to avoid blocking the main loop during AI planning.
func (p *Playfield) planAIMove(player *core.Player) {
	// Get the AI's next move
	action, err := player.Strategy.NextMove(p.game.Board)
	if err != nil {
		rl.TraceLog(rl.LogWarning, "AI player failed to generate move: %v", err)
		// Skip turn if AI has no valid moves
		p.game.AdvanceTurn()
		return
	}

	// Find the current position of the piece
	currentPos, err := p.game.Board.PieceLocation(action.Piece)
	if err != nil {
		rl.TraceLog(rl.LogWarning, "failed to find AI piece location: %v", err)
		p.game.AdvanceTurn()
		return
	}

	// Convert Action to Move
	move, err := p.game.ActionToMove(action)
	if err != nil {
		rl.TraceLog(rl.LogWarning, "failed to convert AI action to move: %v", err)
		p.game.AdvanceTurn()
		return
	}

	// Select the piece to visualize the valid moves
	p.SelectPiece(action.Piece)

	// Wait 1 second before executing the move
	time.Sleep(1 * time.Second)

	// Signal the main loop to execute this move
	p.moveExecutionChan <- struct {
		piece    *core.Piece
		position core.Position
		move     core.Move
	}{
		piece:    action.Piece,
		position: currentPos,
		move:     move,
	}
}

// render draws the current game state to the screen.
func (p *Playfield) render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	if err := p.renderBoard(); err != nil {
		rl.TraceLog(rl.LogError, "error rendering game: %v", err)
	}

	if err := p.renderCapturedPieces(); err != nil {
		rl.TraceLog(rl.LogError, "error rendering captured pieces: %v", err)
	}

	p.renderStatus()

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
// Only allows selecting pieces that belong to the current player.
func (p *Playfield) SelectPiece(piece *core.Piece) {
	// If clicking didn't select a piece, unselect any selected piece
	if piece == nil {
		p.selectedPiece = nil
		return
	}
	// Prevent selecting pieces that don't belong to the current player
	if piece.Color != p.game.ActiveColor {
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

// MouseIsOverBoard returns true if and only if the mouse location is somewhere on the playable board.
func (p *Playfield) MouseIsOverBoard(mouseLoc rl.Vector2) bool {
	adjClickLoc := rl.Vector2{X: mouseLoc.X - p.boardLoc.X, Y: mouseLoc.Y - p.boardLoc.Y}
	return adjClickLoc.X >= 0 && adjClickLoc.X < float32(core.SquareSize*p.game.Board.Columns) &&
		adjClickLoc.Y >= 0 && adjClickLoc.Y < float32(core.SquareSize*p.game.Board.Rows)
}

// PositionUnderMouse returns the board position under a mouse position. If it's outside the board, then an error
// is returned.
func (p *Playfield) PositionUnderMouse(mouseLoc rl.Vector2) (core.Position, error) {
	// Shift the position relative to the board upper corner so the click loc is in board space
	adjClickLoc := rl.Vector2{X: mouseLoc.X - p.boardLoc.X, Y: mouseLoc.Y - p.boardLoc.Y}
	if !p.MouseIsOverBoard(mouseLoc) {
		return core.Position{}, fmt.Errorf("click is outside the board bounds")
	}
	return core.Position{int(adjClickLoc.Y / float32(core.SquareSize)), int(adjClickLoc.X / float32(core.SquareSize))}, nil
}

// PieceUnderMouse returns the piece under a mouse position. If there's no piece there, then nil is returned.
func (p *Playfield) PieceUnderMouse(clickLoc rl.Vector2) *core.Piece {
	pos, err := p.PositionUnderMouse(clickLoc)
	if err != nil {
		return nil
	}
	return p.game.Board.GetPieceAt(pos)
}
