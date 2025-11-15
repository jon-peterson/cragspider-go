# Project Guidelines

This file provides essential context and instructions for Claude when working with this Go project.

## Project Overview

This project is a simple REST API built with Go, using the `gorilla/mux` router for routing and `gorm.io/gorm` for ORM.
It manages user data with basic CRUD operations.

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
* When suggesting code, prioritize idiomatic Go practices.
* If refactoring, aim for modularity and testability.
* When adding new features, ensure proper error handling and logging are implemented.
* If I ask you to debug an issue, consider potential race conditions or concurrency issues in Go.
