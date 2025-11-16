// Copyright 2025 Ideograph LLC. All rights reserved.

package scenes

import (
	"cragspider-go/internal/core"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/samber/lo"
)

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
