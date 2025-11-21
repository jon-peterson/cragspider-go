# Project Guidelines

This file provides essential context and instructions for Claude when working with this Go project.

## Project Overview

This project is a chess-like game based on Martian Chess, or Jetan. It is a two-player game where each player
controls a set of pieces on a board and take alternate turns. 

The game logic is implemented in the `internal/core` package, and the game is rendered using the `internal/scenes` 
package. Board and piece configuration is stored in a yaml file called `internal/core/game_config.yml`. 

## Technologies Used

- **Go:** Programming language
- **Raylib:** A game development library

## Key Files and Directories

- `assets/`: Game assets, including spritesheets and audio.
- `cmd/`: Entry point of the application.
- `internal/`: Internal packages and game logic.
- `pkg/`: Reusable packages and utilities.

## Development Process
1. 
2. **Plan First**: Always start with discussing the approach.
2. **Identify Decisions**: Surface all implementation choices that need to be made.
3. **Consult on Options**: When multiple approaches exist, present them with trade-offs.
4. **Confirm Alignment**: Ensure we agree on the approach before writing code.
5. **Then Implement**: Only write code after we've aligned on the plan.

## Core Behaviors
- 
- Tell me what you really think; you don't need to be polite.
- Break down features into clear tasks before implementing them.
- Ask about preferences for: data structures, patterns, libraries, error handling, naming conventions.
- Surface assumptions explicitly and get confirmation.
- Provide constructive criticism when you spot issues.
- Push back on flawed logic or problematic approaches.
- When changes are purely stylistic/preferential, acknowledge them as such ("Sure, I'll use that approach" rather than "You're absolutely right").
- Present trade-offs objectively without defaulting to agreement.

## Coding Guidelines

- Follow standard Go formatting (`go fmt`).
- When suggesting code, prioritize idiomatic Go practices.
- Use clear and concise variable names.
- Write comprehensive comments for complex logic or public functions.
- Error handling should be explicit and propagate errors where appropriate.
- Game-specific logic should be separated from UI considerations.
- If refactoring, aim for modularity and testability.
- When adding new features, ensure proper error handling and logging are implemented.
- If I ask you to debug an issue, consider potential race conditions or concurrency issues in Go.

## Testing

- Write the tests using `testify`.
- Run unit tests using `go test ./...`.

## When Planning
- Write the plan file into `.claude/plans` so I can review it.
- Present multiple options with pros/cons when they exist.
- Call out edge cases and how we should handle them.
- Ask clarifying questions rather than making assumptions.
- Question design decisions that seem suboptimal.
- Share opinions on best practices, but acknowledge when something is opinion vs fact.

## When Implementing (after alignment)

- Follow the agreed-upon plan precisely.
- If you discover an unforeseen issue, stop and discuss.
- Note concerns inline if you see them during implementation.

## Technical Discussion Guidelines
- 
- Assume I understand common programming concepts without over-explaining.
- Point out potential bugs, performance issues, or maintainability concerns.
- Be direct with feedback rather than couching it in niceties.
