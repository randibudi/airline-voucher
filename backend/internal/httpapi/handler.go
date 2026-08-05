package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/randibudi/airline-voucher/backend/internal/domain"
	"github.com/randibudi/airline-voucher/backend/internal/service"
)

// Handler exposes voucher use cases over HTTP.
type Handler struct {
	service *service.Voucher
}

// NewHandler creates voucher HTTP handlers.
func NewHandler(voucherService *service.Voucher) *Handler {
	return &Handler{service: voucherService}
}

type checkRequest struct {
	FlightNumber string `json:"flightNumber"`
	FlightDate   string `json:"flightDate"`
}

type generateRequest struct {
	CrewName     string `json:"crewName"`
	CrewID       string `json:"crewId"`
	FlightNumber string `json:"flightNumber"`
	FlightDate   string `json:"flightDate"`
	Aircraft     string `json:"aircraft"`
}

type voucherResponse struct {
	CrewName     string   `json:"crewName"`
	CrewID       string   `json:"crewId"`
	FlightNumber string   `json:"flightNumber"`
	FlightDate   string   `json:"flightDate"`
	Aircraft     string   `json:"aircraft"`
	Seats        []string `json:"seats"`
}

// Check handles POST /api/check.
func (h *Handler) Check(c fiber.Ctx) error {
	var request checkRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return errorJSON(c, http.StatusBadRequest, "request body must be valid JSON")
	}

	voucher, exists, err := h.service.Check(c.Context(), service.CheckInput{
		FlightNumber: request.FlightNumber,
		FlightDate:   request.FlightDate,
	})
	if err != nil {
		return handleServiceError(c, err)
	}

	var responseVoucher *voucherResponse
	if exists {
		mapped := mapVoucher(voucher)
		responseVoucher = &mapped
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"exists":  exists,
		"voucher": responseVoucher,
	})
}

// Generate handles POST /api/generate.
func (h *Handler) Generate(c fiber.Ctx) error {
	var request generateRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return errorJSON(c, http.StatusBadRequest, "request body must be valid JSON")
	}

	voucher, err := h.service.Generate(c.Context(), service.GenerateInput{
		CrewName:     request.CrewName,
		CrewID:       request.CrewID,
		FlightNumber: request.FlightNumber,
		FlightDate:   request.FlightDate,
		Aircraft:     request.Aircraft,
	})
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.Status(http.StatusOK).JSON(mapVoucher(voucher))
}

func mapVoucher(voucher domain.Voucher) voucherResponse {
	return voucherResponse{
		CrewName:     voucher.CrewName,
		CrewID:       voucher.CrewID,
		FlightNumber: voucher.FlightNumber,
		FlightDate:   voucher.FlightDate,
		Aircraft:     string(voucher.Aircraft),
		Seats:        []string{voucher.Seats[0], voucher.Seats[1], voucher.Seats[2]},
	}
}

func handleServiceError(c fiber.Ctx, err error) error {
	if errors.Is(err, service.ErrValidation) {
		return errorJSON(c, http.StatusBadRequest, err.Error())
	}
	return errorJSON(c, http.StatusInternalServerError, "internal server error")
}

func errorJSON(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{"message": message},
	})
}
