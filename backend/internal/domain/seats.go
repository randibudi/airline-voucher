package domain

import (
	"fmt"
	"math/rand"
)

// RandomSource provides the controlled randomness needed by seat generation.
type RandomSource interface {
	Intn(n int) int
}

// SeatGenerator creates unique seats for supported aircraft.
type SeatGenerator struct {
	random RandomSource
}

// NewSeatGenerator creates a seat generator with an injectable random source.
func NewSeatGenerator(random RandomSource) *SeatGenerator {
	return &SeatGenerator{random: random}
}

// NewDefaultSeatGenerator creates a concurrency-safe generator for production use.
func NewDefaultSeatGenerator() *SeatGenerator {
	return NewSeatGenerator(globalRandom{})
}

type globalRandom struct{}

func (globalRandom) Intn(n int) int {
	return rand.Intn(n)
}

// Generate returns exactly three unique seats valid for the aircraft.
func (g *SeatGenerator) Generate(aircraft Aircraft) ([3]string, error) {
	var seats [3]string
	layout, ok := supportedAircraft[aircraft]
	if !ok {
		return seats, fmt.Errorf("unsupported aircraft %q", aircraft)
	}

	available := make([]string, 0, layout.maxRow*len(layout.letters))
	for row := 1; row <= layout.maxRow; row++ {
		for _, letter := range layout.letters {
			available = append(available, fmt.Sprintf("%d%c", row, letter))
		}
	}

	for i := range seats {
		index := g.random.Intn(len(available))
		seats[i] = available[index]
		available[index] = available[len(available)-1]
		available = available[:len(available)-1]
	}
	return seats, nil
}
