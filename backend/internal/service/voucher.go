package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/randibudi/airline-voucher/backend/internal/domain"
	"github.com/randibudi/airline-voucher/backend/internal/repository"
)

// ErrValidation identifies invalid client input.
var ErrValidation = errors.New("validation error")

// VoucherRepository defines the persistence operations needed by the service.
type VoucherRepository interface {
	FindByFlight(ctx context.Context, flightNumber, flightDate string) (domain.Voucher, bool, error)
	Create(ctx context.Context, voucher *domain.Voucher) error
}

// SeatGenerator defines seat generation used by the service.
type SeatGenerator interface {
	Generate(aircraft domain.Aircraft) ([3]string, error)
}

// CheckInput contains a voucher lookup request.
type CheckInput struct {
	FlightNumber string
	FlightDate   string
}

// GenerateInput contains a voucher creation request.
type GenerateInput struct {
	CrewName     string
	CrewID       string
	FlightNumber string
	FlightDate   string
	Aircraft     string
}

// Voucher coordinates voucher lookup and generation.
type Voucher struct {
	repository VoucherRepository
	generator  SeatGenerator
}

// NewVoucher creates a voucher service.
func NewVoucher(repository VoucherRepository, generator SeatGenerator) *Voucher {
	return &Voucher{repository: repository, generator: generator}
}

// Check finds a voucher by normalized flight number and date.
func (s *Voucher) Check(ctx context.Context, input CheckInput) (domain.Voucher, bool, error) {
	flightNumber, err := validateFlight(input.FlightNumber, input.FlightDate)
	if err != nil {
		return domain.Voucher{}, false, err
	}
	return s.repository.FindByFlight(ctx, flightNumber, input.FlightDate)
}

// Generate returns an existing voucher or atomically creates a new one.
func (s *Voucher) Generate(ctx context.Context, input GenerateInput) (domain.Voucher, error) {
	input.CrewName = strings.TrimSpace(input.CrewName)
	input.CrewID = strings.TrimSpace(input.CrewID)
	if input.CrewName == "" {
		return domain.Voucher{}, validationError("crewName is required")
	}
	if input.CrewID == "" {
		return domain.Voucher{}, validationError("crewId is required")
	}

	flightNumber, err := validateFlight(input.FlightNumber, input.FlightDate)
	if err != nil {
		return domain.Voucher{}, err
	}
	aircraft, err := domain.ParseAircraft(input.Aircraft)
	if err != nil {
		return domain.Voucher{}, validationError("aircraft must be ATR, Airbus 320, or Boeing 737 Max")
	}

	existing, found, err := s.repository.FindByFlight(ctx, flightNumber, input.FlightDate)
	if err != nil {
		return domain.Voucher{}, err
	}
	if found {
		return existing, nil
	}

	seats, err := s.generator.Generate(aircraft)
	if err != nil {
		return domain.Voucher{}, fmt.Errorf("generate seats: %w", err)
	}
	voucher := domain.Voucher{
		CrewName:     input.CrewName,
		CrewID:       input.CrewID,
		FlightNumber: flightNumber,
		FlightDate:   input.FlightDate,
		Aircraft:     aircraft,
		Seats:        seats,
	}
	if err := s.repository.Create(ctx, &voucher); err == nil {
		return voucher, nil
	} else if !errors.Is(err, repository.ErrDuplicateVoucher) {
		return domain.Voucher{}, err
	}

	existing, found, err = s.repository.FindByFlight(ctx, flightNumber, input.FlightDate)
	if err != nil {
		return domain.Voucher{}, err
	}
	if !found {
		return domain.Voucher{}, errors.New("voucher conflict occurred but voucher was not found")
	}
	return existing, nil
}

func validateFlight(flightNumber, flightDate string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(flightNumber))
	if normalized == "" {
		return "", validationError("flightNumber is required")
	}
	if flightDate == "" {
		return "", validationError("flightDate is required")
	}
	parsed, err := time.Parse("2006-01-02", flightDate)
	if err != nil || parsed.Format("2006-01-02") != flightDate {
		return "", validationError("flightDate must use YYYY-MM-DD")
	}
	return normalized, nil
}

func validationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}
