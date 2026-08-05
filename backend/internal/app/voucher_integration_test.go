package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/randibudi/airline-voucher/backend/internal/app"
	"github.com/randibudi/airline-voucher/backend/internal/domain"
	"github.com/randibudi/airline-voucher/backend/internal/httpapi"
	"github.com/randibudi/airline-voucher/backend/internal/repository"
	"github.com/randibudi/airline-voucher/backend/internal/service"
)

type voucherJSON struct {
	CrewName     string   `json:"crewName"`
	CrewID       string   `json:"crewId"`
	FlightNumber string   `json:"flightNumber"`
	FlightDate   string   `json:"flightDate"`
	Aircraft     string   `json:"aircraft"`
	Seats        []string `json:"seats"`
}

func TestVoucherGenerateCheckAndIdempotentGenerate(t *testing.T) {
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

	voucherRepository := repository.NewSQLite(db)
	if err := voucherRepository.Initialize(ctx); err != nil {
		t.Fatalf("initialize repository: %v", err)
	}
	generator := domain.NewSeatGenerator(rand.New(rand.NewSource(1)))
	voucherService := service.NewVoucher(voucherRepository, generator)
	application := app.New(httpapi.NewHandler(voucherService))

	generateBody := []byte(`{
		"crewName":"John Doe",
		"crewId":"CRW001",
		"flightNumber":" ga123 ",
		"flightDate":"2026-08-05",
		"aircraft":"ATR"
	}`)
	first := postJSON[voucherJSON](t, application, "/api/generate", generateBody)
	if first.FlightNumber != "GA123" {
		t.Errorf("flightNumber = %q, want GA123", first.FlightNumber)
	}
	if len(first.Seats) != 3 {
		t.Fatalf("seat count = %d, want 3", len(first.Seats))
	}

	checkBody := []byte(`{"flightNumber":"GA123","flightDate":"2026-08-05"}`)
	checked := postJSON[struct {
		Exists  bool         `json:"exists"`
		Voucher *voucherJSON `json:"voucher"`
	}](t, application, "/api/check", checkBody)
	if !checked.Exists || checked.Voucher == nil {
		t.Fatal("check response does not contain the generated voucher")
	}
	if !reflect.DeepEqual(*checked.Voucher, first) {
		t.Errorf("checked voucher = %+v, want %+v", *checked.Voucher, first)
	}

	second := postJSON[voucherJSON](t, application, "/api/generate", generateBody)
	if !reflect.DeepEqual(second, first) {
		t.Errorf("second voucher = %+v, want %+v", second, first)
	}
}

func postJSON[T any](t *testing.T, application *fiber.App, path string, body []byte) T {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := application.Test(request)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Errorf("close response body: %v", err)
		}
	}()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("POST %s status = %d, want %d", path, response.StatusCode, http.StatusOK)
	}

	var result T
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode POST %s response: %v", path, err)
	}
	return result
}
