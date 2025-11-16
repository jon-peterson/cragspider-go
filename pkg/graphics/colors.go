// Copyright 2025 Ideograph LLC. All rights reserved.

package graphics

import rl "github.com/gen2brain/raylib-go/raylib"

// LightenColor mixes the given color with white by the specified amount (0.0 to 1.0).
func LightenColor(original rl.Color, amount float32) rl.Color {
	amount = rl.Clamp(amount, 0.0, 1.0)
	newR := float32(original.R)*(1.0-amount) + float32(255)*amount
	newG := float32(original.G)*(1.0-amount) + float32(255)*amount
	newB := float32(original.B)*(1.0-amount) + float32(255)*amount
	return rl.NewColor(
		byte(newR),
		byte(newG),
		byte(newB),
		original.A,
	)
}
