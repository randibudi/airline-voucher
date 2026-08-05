package domain_test

import (
	"math/rand"
	"testing"

	"github.com/randibudi/airline-voucher/backend/internal/domain"
)

func TestSeatGeneratorCreatesThreeUniqueValidSeats(t *testing.T) {
	tests := []struct {
		name     string
		aircraft domain.Aircraft
	}{
		{name: "ATR", aircraft: domain.AircraftATR},
		{name: "Airbus 320", aircraft: domain.AircraftAirbus320},
		{name: "Boeing 737 Max", aircraft: domain.AircraftBoeing737Max},
	}

	for index, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generator := domain.NewSeatGenerator(rand.New(rand.NewSource(int64(index + 1))))
			seats, err := generator.Generate(test.aircraft)
			if err != nil {
				t.Fatalf("generate seats: %v", err)
			}

			unique := make(map[string]struct{}, len(seats))
			for _, seat := range seats {
				if !test.aircraft.IsValidSeat(seat) {
					t.Errorf("seat %q is invalid for %q", seat, test.aircraft)
				}
				unique[seat] = struct{}{}
			}
			if len(unique) != 3 {
				t.Errorf("unique seat count = %d, want 3; seats = %v", len(unique), seats)
			}
		})
	}
}
