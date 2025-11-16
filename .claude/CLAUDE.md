# Project Guidelines

This file provides essential context and instructions for Claude when working with this Go project.

## Project Overview

This project is a chess-like game based on Martian Chess, or Jetan. It is a two-player game where each player
controls a set of pieces on a board and take alternate turns. 

The game logic is implemented in the `internal/core` package, and the game is rendered using the `internal/scenes` 
package. Board and piece configuration is stored in a yaml file called `internal/core/game_config.yml`. 

## Technologies Used

* **Go:** Programming language
* **Raylib:** A game development library

## Code Style Guidelines

* Follow standard Go formatting (`go fmt`).
* Use clear and concise variable names.
* Write comprehensive comments for complex logic or public functions.
* Error handling should be explicit and propagate errors where appropriate.

## Testing Instructions

* Run unit tests using `go test ./...`.

## Key Files and Directories

* `assets/`: Game assets, including spritesheets and audio.
* `cmd/`: Entry point of the application.
* `internal/`: Internal packages and game logic.
* `pkg/`: Reusable packages and utilities.

## Specific Instructions for Claude

* Tell me what you really think; you don't need to be polite.
* Game-specific logic should be separated from UI considerations.
* When suggesting code, prioritize idiomatic Go practices.
* If refactoring, aim for modularity and testability.
* When adding new features, ensure proper error handling and logging are implemented.
* If I ask you to debug an issue, consider potential race conditions or concurrency issues in Go.
