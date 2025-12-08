// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"cragspider-go/internal/core"
	"cragspider-go/pkg/graphics"
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/samber/lo"
)

// positionTintMap is a map of board positions to their tint colors
type positionTintMap map[core.Position]color.RGBA

// renderBoard draws the board to the screen with the given board location (where the upper left corner is).
func (p *Playfield) renderBoard() error {
	// First draw the board itself
	tints := p.getTintedPositions(rl.GetMousePosition())
	for i := range p.game.Board.Rows {
		for j := range p.game.Board.Columns {
			pos := core.Position{i, j}
			tint, exists := tints[pos]
			if !exists {
				tint = rl.White
			}
			err := p.backgroundSprites.DrawFrame(
				p.game.Board.GetSquareAt(pos).Frame,
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

// getTintedPositions returns a map of positions on the board that should be tinted to their corresponding colors.
// It checks for pieces under the mouse first, then falls back to the selected piece.
func (p *Playfield) getTintedPositions(mousePos rl.Vector2) positionTintMap {
	tints := make(positionTintMap)
	// If mouse is hovering over a piece, tint it and where it can move to
	pieceUnderMouse := p.PieceUnderMouse(mousePos)
	if pieceUnderMouse != nil {
		// If the user has selected a piece, don't tint any of its other pieces
		if p.selectedPiece == nil || p.selectedPiece.Piece.Color != pieceUnderMouse.Color {
			pos, err := p.game.Board.PieceLocation(pieceUnderMouse)
			if err == nil {
				// Tint the piece under the mouse hover and its valid moves
				tint := graphics.LightenColor(lo.Ternary(pieceUnderMouse.Color == core.White, rl.Green, rl.Red), 0.75)
				tints[pos] = tint
				positions := pieceUnderMouse.ValidNextPositions(pos, p.game.Board)
				for _, movePos := range positions {
					tints[movePos] = tint
				}
			}
		}
	}
	// The selected piece and all its valid moves should be tinted
	if p.selectedPiece != nil {
		tint := lo.Ternary(p.selectedPiece.Piece.Color == core.White, rl.Green, rl.Red)
		tints[p.selectedPiece.Position] = tint
		positions := p.selectedPiece.Piece.ValidNextPositions(p.selectedPiece.Position, p.game.Board)
		for _, pos := range positions {
			tints[pos] = tint
		}
	}
	return tints
}

// renderPieceOnBoard renders a single piece on the board at the specified position.
func (p *Playfield) renderPieceOnBoard(piece *core.Piece, j int, i int) error {
	isSelected := p.selectedPiece != nil && p.selectedPiece.Piece == piece
	frame := lo.Ternary(isSelected, 0, 1)
	location := rl.Vector2{X: p.boardLoc.X + float32(j*core.SquareSize), Y: p.boardLoc.Y + float32(i*core.SquareSize)}
	return p.renderPieceAtLocationWithFrame(piece, location, frame)
}

// renderCapturedPieces renders the captured pieces on the sides of the board.
// White captured pieces are displayed on the left, black on the right.
func (p *Playfield) renderCapturedPieces() error {
	const (
		padding        = 20
		piecePadding   = 10
		capturedPerRow = 3
	)

	// Render white captured pieces on the left
	whiteCaptured := p.game.Board.GetCapturedPieces(core.White)
	for idx, piece := range whiteCaptured {
		row := idx / capturedPerRow
		col := idx % capturedPerRow
		x := p.boardLoc.X - float32((capturedPerRow*core.SquareSize)+padding)
		y := p.boardLoc.Y + float32(row*(core.SquareSize+piecePadding))

		if err := p.renderPieceAtLocationWithFrame(piece, rl.Vector2{X: x + float32(col*core.SquareSize), Y: y}, 1); err != nil {
			return err
		}
	}

	// Render black captured pieces on the right
	blackCaptured := p.game.Board.GetCapturedPieces(core.Black)
	for idx, piece := range blackCaptured {
		row := idx / capturedPerRow
		col := idx % capturedPerRow
		x := p.boardLoc.X + float32(p.game.Board.Columns*core.SquareSize+padding)
		y := p.boardLoc.Y + float32(row*(core.SquareSize+piecePadding))

		if err := p.renderPieceAtLocationWithFrame(piece, rl.Vector2{X: x + float32(col*core.SquareSize), Y: y}, 1); err != nil {
			return err
		}
	}

	return nil
}

// renderPieceAtLocationWithFrame renders a piece at the specified screen location with a specific frame.
func (p *Playfield) renderPieceAtLocationWithFrame(piece *core.Piece, location rl.Vector2, frame int) error {
	sheet := lo.Ternary(piece.Color == core.White, p.whiteSprites, p.blackSprites)
	err := sheet.DrawFrame(
		piece.Config.Sprites[piece.Color][frame],
		location,
		core.Scale,
		rl.Vector2{X: 1.0, Y: 0.0},
		rl.White)
	if err != nil {
		return fmt.Errorf("failed to draw piece: %w", err)
	}
	return nil
}

// renderStatus renders all status needs of the playfield, like whose turn it is.
func (p *Playfield) renderStatus() {
	turnText := lo.Ternary(p.game.ActiveColor == core.White, "White's Turn", "Black's Turn")
	turnColor := lo.Ternary(p.game.ActiveColor == core.White, rl.Black, rl.DarkGray)
	fontSize := int32(24)
	x := int32(20)
	y := int32(20)
	rl.DrawText(turnText, x, y, fontSize, turnColor)
}
