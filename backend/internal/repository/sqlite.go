package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/randibudi/airline-voucher/backend/internal/domain"
	"modernc.org/sqlite"
)

// ErrDuplicateVoucher indicates a unique flight number and date conflict.
var ErrDuplicateVoucher = errors.New("voucher already exists")

// SQLiteRepository stores vouchers in SQLite.
type SQLiteRepository struct {
	db *sql.DB
}

// OpenSQLite opens a SQLite database and verifies the connection.
func OpenSQLite(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return db, nil
}

// NewSQLite creates a voucher repository over an open database.
func NewSQLite(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// Initialize creates the application schema when it does not exist.
func (r *SQLiteRepository) Initialize(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS vouchers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	crew_name TEXT NOT NULL,
	crew_id TEXT NOT NULL,
	flight_number TEXT NOT NULL,
	flight_date TEXT NOT NULL,
	aircraft TEXT NOT NULL,
	seat_1 TEXT NOT NULL,
	seat_2 TEXT NOT NULL,
	seat_3 TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(flight_number, flight_date)
);`
	if _, err := r.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}
	return nil
}

// FindByFlight returns a voucher and whether it exists.
func (r *SQLiteRepository) FindByFlight(ctx context.Context, flightNumber, flightDate string) (domain.Voucher, bool, error) {
	const query = `
SELECT id, crew_name, crew_id, flight_number, flight_date, aircraft,
       seat_1, seat_2, seat_3, created_at
FROM vouchers
WHERE flight_number = ? AND flight_date = ?`

	var voucher domain.Voucher
	var aircraft string
	err := r.db.QueryRowContext(ctx, query, flightNumber, flightDate).Scan(
		&voucher.ID,
		&voucher.CrewName,
		&voucher.CrewID,
		&voucher.FlightNumber,
		&voucher.FlightDate,
		&aircraft,
		&voucher.Seats[0],
		&voucher.Seats[1],
		&voucher.Seats[2],
		&voucher.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Voucher{}, false, nil
	}
	if err != nil {
		return domain.Voucher{}, false, fmt.Errorf("find voucher: %w", err)
	}

	voucher.Aircraft = domain.Aircraft(aircraft)
	return voucher, true, nil
}

// Create persists a voucher and populates its generated fields.
func (r *SQLiteRepository) Create(ctx context.Context, voucher *domain.Voucher) error {
	const query = `
INSERT INTO vouchers (
	crew_name, crew_id, flight_number, flight_date, aircraft,
	seat_1, seat_2, seat_3
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		voucher.CrewName,
		voucher.CrewID,
		voucher.FlightNumber,
		voucher.FlightDate,
		voucher.Aircraft,
		voucher.Seats[0],
		voucher.Seats[1],
		voucher.Seats[2],
	).Scan(&voucher.ID, &voucher.CreatedAt)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 {
			return ErrDuplicateVoucher
		}
		return fmt.Errorf("create voucher: %w", err)
	}
	return nil
}
