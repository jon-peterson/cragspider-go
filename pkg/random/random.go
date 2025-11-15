package random

import "math/rand"

// IntInRange returns a random int between min and max, inclusive. Panics if min >= max.
func IntInRange(min, max int) int {
	if min >= max {
		panic("min must be less than max")
	}
	return min + rand.Intn(max-min+1)
}

// Chance returns true if a random number between 0 and 1 is less than chance.
// Panics if chance is outside the range 0-1.
func Chance(chance float32) bool {
	if chance < 0 || chance > 1 {
		panic("chance must be between 0 and 1")
	}
	return rand.Float32() < chance
}

// Choice returns a random element from the given slice. Panics if the slice is empty.
func Choice[T any](items []T) T {
	if len(items) == 0 {
		panic("cannot choose from empty slice")
	}
	return items[rand.Intn(len(items))]
}
