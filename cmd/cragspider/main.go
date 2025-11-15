// Copyright 2025 Ideograph LLC. All rights reserved.

package main

import (
	"cragspider-go/internal/scenes"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

// main is the entry point for the Cragspider game. Spawns a window and runs the game inside.
func main() {
	rl.InitWindow(screenWidth, screenHeight, "Cragspider")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)

	if os.Getenv("DEBUG") != "" {
		rl.SetTraceLogLevel(rl.LogDebug)
	}

	// TODO: Start with AttractModeScene when implemented
	sceneCode := scenes.GameplayScene
	for sceneCode != scenes.Quit {
		rl.TraceLog(rl.LogInfo, "Starting scene code %v", sceneCode)
		scene := initScene(sceneCode)
		sceneCode = scene.Loop()
		scene.Close()
	}
}

// initScene initializes and returns the scene corresponding to the given scene code.
func initScene(code scenes.SceneCode) scenes.Scene {
	switch code {
	case scenes.AttractModeScene:
		// TODO
	case scenes.GameplayScene:
		gm := &scenes.Playfield{}
		gm.Init(screenWidth, screenHeight)
		return gm
	case scenes.GameOverScene:
		// TODO
	default:
		rl.TraceLog(rl.LogError, "Unknown or unimplemented scene code %v", code)
	}
	return nil
}
