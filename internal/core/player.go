// Copyright 2025 Ideograph LLC. All rights reserved.

package core

// Player represents a player in the game.
type Player struct {
	Name     string
	Strategy AgentStrategy
}

// NewHumanPlayer returns a new player set up as a Human.
func NewHumanPlayer() *Player {
	return &Player{Name: "Human"}
}

// NewAIPlayer returns a new player object with the specified AI strategy.
func NewAIPlayer(name string, agentStrategy AgentStrategy) *Player {
	return &Player{Name: name, Strategy: agentStrategy}
}

// String returns the name of the player, human or AI.
func (p Player) String() string {
	if p.Strategy == nil {
		return "Human"
	}
	return p.Name
}

// IsHuman returns true if the player is human controlled.
func (p Player) IsHuman() bool {
	return p.Strategy == nil
}

// IsAI returns true if the player is AI controlled.
func (p Player) IsAI() bool {
	return p.Strategy != nil
}
