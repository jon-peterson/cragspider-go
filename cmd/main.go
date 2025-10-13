package main

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1024.0
	screenHeight = 768.0
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Cragspider")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	if os.Getenv("DEBUG") != "" {
		rl.SetTraceLogLevel(rl.LogDebug)
	}
}
