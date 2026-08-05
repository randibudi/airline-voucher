package domain

import (
	"fmt"
	"strconv"
)

// Aircraft identifies a supported aircraft seat layout.
type Aircraft string

const (
	AircraftATR          Aircraft = "ATR"
	AircraftAirbus320    Aircraft = "Airbus 320"
	AircraftBoeing737Max Aircraft = "Boeing 737 Max"
)

var supportedAircraft = map[Aircraft]seatLayout{
	AircraftATR:          {maxRow: 18, letters: "ACDF"},
	AircraftAirbus320:    {maxRow: 32, letters: "ABCDEF"},
	AircraftBoeing737Max: {maxRow: 32, letters: "ABCDEF"},
}

type seatLayout struct {
	maxRow  int
	letters string
}

// ParseAircraft accepts only the canonical aircraft names exposed by the API.
func ParseAircraft(value string) (Aircraft, error) {
	aircraft := Aircraft(value)
	if _, ok := supportedAircraft[aircraft]; !ok {
		return "", fmt.Errorf("unsupported aircraft %q", value)
	}
	return aircraft, nil
}

// IsValidSeat reports whether a seat belongs to the aircraft layout.
func (a Aircraft) IsValidSeat(seat string) bool {
	layout, ok := supportedAircraft[a]
	if !ok || len(seat) < 2 {
		return false
	}

	row, err := strconv.Atoi(seat[:len(seat)-1])
	if err != nil || row < 1 || row > layout.maxRow {
		return false
	}

	letter := seat[len(seat)-1]
	for i := range layout.letters {
		if layout.letters[i] == letter {
			return true
		}
	}
	return false
}

// Voucher contains one persisted seat assignment.
type Voucher struct {
	ID           int64
	CrewName     string
	CrewID       string
	FlightNumber string
	FlightDate   string
	Aircraft     Aircraft
	Seats        [3]string
	CreatedAt    string
}
