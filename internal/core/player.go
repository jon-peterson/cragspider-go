// Copyright 2025 Ideograph LLC. All rights reserved.

package core

import "fmt"

// Player represents a player in the game.
type Player struct {
	Name     string
	strategy AgentStrategy
}

// newHumanPlayer returns a new player set up as a Human.
func newHumanPlayer() *Player {
	return &Player{Name: "Human"}
}

// newAgentPlayer returns a new player object with the specified AI strategy.
func newAIPlayer(name string, agentStrategy AgentStrategy) *Player {
	return &Player{Name: name, strategy: agentStrategy}
}

// String returns the name of the player, human or AI.
func (p Player) String() string {
	if p.strategy == nil {
		return "Human"
	}
	return fmt.Sprintf("%s", p.Name)
}

// IsHuman returns true if the player is human controlled.
func (p Player) IsHuman() bool {
	return p.strategy == nil
}
