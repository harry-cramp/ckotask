package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/service"
	"github.com/go-chi/chi/v5"
)

type PaymentsHandler struct {
	storage    *repository.PaymentsRepository
	payService *service.PaymentService
}

func NewPaymentsHandler(
	storage *repository.PaymentsRepository,
	payService *service.PaymentService,
) *PaymentsHandler {
	return &PaymentsHandler{
		storage:    storage,
		payService: payService,
	}
}

// GetHandler returns an http.HandlerFunc that handles HTTP GET requests.
// It retrieves a payment record by its ID from the storage.
// The ID is expected to be part of the URL.
func (h *PaymentsHandler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		payment := h.storage.GetPayment(id)

		if payment != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(payment); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// PostHandler returns an http.HandlerFunc that handles HTTP POST requests.
// It makes a payment request to the bank and stores the record by ID in storage.
// The card info and payment amount should be included in the request body.
func (h *PaymentsHandler) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var paymentRequest models.PostPaymentRequest
		err := json.NewDecoder(r.Body).Decode(&paymentRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bankResp, err := h.payService.Process(service.PaymentParams{
			CardNumber: paymentRequest.CardNumber,
			Amount:     paymentRequest.Amount,
			Cvv:        paymentRequest.Cvv,
			ExpMonth:   paymentRequest.ExpiryMonth,
			ExpYear:    paymentRequest.ExpiryYear,
			Currency:   paymentRequest.Currency,
		})
		if err != nil {
			if err == service.ErrUpstreamUnavailable {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else if bankResp.PaymentStatus == "Rejected" {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			_ = json.NewEncoder(w).Encode(bankResp)
			return
		}

		h.storage.AddPayment(bankResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(bankResp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
