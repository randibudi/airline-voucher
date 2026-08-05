package repository_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/randibudi/airline-voucher/backend/internal/domain"
	"github.com/randibudi/airline-voucher/backend/internal/repository"
)

func TestSQLiteRepositoryPersistsVoucherAndEnforcesUniqueness(t *testing.T) {
	ctx := context.Background()
	db, err := repository.OpenSQLite(ctx, filepath.Join(t.TempDir(), "vouchers.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("close database: %v", err)
		}
	}()

	vouchers := repository.NewSQLite(db)
	if err := vouchers.Initialize(ctx); err != nil {
		t.Fatalf("initialize repository: %v", err)
	}

	voucher := domain.Voucher{
		CrewName:     "John Doe",
		CrewID:       "CRW001",
		FlightNumber: "GA123",
		FlightDate:   "2026-08-05",
		Aircraft:     domain.AircraftATR,
		Seats:        [3]string{"1A", "2C", "3F"},
	}
	if err := vouchers.Create(ctx, &voucher); err != nil {
		t.Fatalf("create voucher: %v", err)
	}

	found, exists, err := vouchers.FindByFlight(ctx, "GA123", "2026-08-05")
	if err != nil {
		t.Fatalf("find voucher: %v", err)
	}
	if !exists {
		t.Fatal("voucher does not exist after create")
	}
	if found.ID == 0 || found.CreatedAt == "" {
		t.Errorf("generated fields not populated: ID=%d CreatedAt=%q", found.ID, found.CreatedAt)
	}
	if found.Seats != voucher.Seats {
		t.Errorf("seats = %v, want %v", found.Seats, voucher.Seats)
	}

	duplicate := voucher
	duplicate.ID = 0
	duplicate.CreatedAt = ""
	if err := vouchers.Create(ctx, &duplicate); !errors.Is(err, repository.ErrDuplicateVoucher) {
		t.Errorf("duplicate create error = %v, want ErrDuplicateVoucher", err)
	}
}
