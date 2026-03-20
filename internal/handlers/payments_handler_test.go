package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentHandler(t *testing.T) {
	payment := models.PostPaymentResponse{
		Id:                 "test-id",
		PaymentStatus:      "test-successful-status",
		CardNumberLastFour: 1234,
		ExpiryMonth:        10,
		ExpiryYear:         2035,
		Currency:           "GBP",
		Amount:             100,
	}
	pr := repository.NewPaymentsRepository()
	pr.AddPayment(service.BankResponse{PostPaymentResponse: payment})

	payments := NewPaymentsHandler(pr, &service.PaymentService{})

	r := chi.NewRouter()
	r.Get("/api/payments/{id}", payments.GetHandler())

	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() error {
		return httpServer.ListenAndServe()
	}()

	t.Run("PaymentFound", func(t *testing.T) {
		// Create a new HTTP request for testing
		req, _ := http.NewRequest("GET", "/api/payments/test-id", nil)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
	t.Run("PaymentNotFound", func(t *testing.T) {
		// Create a new HTTP request for testing with a non-existing payment ID
		req, _ := http.NewRequest("GET", "/api/payments/NonExistingID", nil)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the HTTP status code in the response
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

type MockBankClient struct{}

func (m *MockBankClient) Charge(p domain.Payment) (service.BankResponse, error) {
	lastFour := p.CardNumber.GetLastFour()

	status := "Declined"
	if lastFour%2 == 1 {
		status = "Authorized"
	} else if lastFour%10 == 0 {
		return service.BankResponse{}, service.ErrUpstreamUnavailable
	}

	return service.BankResponse{
		PostPaymentResponse: models.PostPaymentResponse{
			Id:            "test123",
			PaymentStatus: status,
		},
	}, nil
}

func TestPostPaymentHandler(t *testing.T) {
	pr := repository.NewPaymentsRepository()
	bc := &MockBankClient{}
	ps := service.NewPaymentService(bc)

	payments := NewPaymentsHandler(pr, ps)

	r := chi.NewRouter()
	r.Post("/api/payments", payments.PostHandler())

	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() error {
		return httpServer.ListenAndServe()
	}()

	tests := []struct {
		name          string
		cardNumber    int
		paymentStatus string
		httpStatus    int
	}{
		{"make payment - authorised", 12341234123411, "Authorized", http.StatusOK},
		{"make payment - declined", 12341234123412, "Declined", http.StatusOK},
		{"make payment - rejected", 12341234, "Rejected", http.StatusBadRequest},
		{"make payment - service unavailable", 123412341234120, "", http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentReq := models.PostPaymentRequest{
				CardNumber:  tt.cardNumber,
				ExpiryMonth: 4,
				ExpiryYear:  2029,
				Currency:    "GBP",
				Amount:      100,
				Cvv:         123,
			}

			jsonData, _ := json.Marshal(paymentReq)

			// Create a new HTTP request for testing
			req, _ := http.NewRequest("POST", "/api/payments", bytes.NewBuffer(jsonData))

			// Create a new HTTP request recorder for recording the response
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			// Check the body is not nil
			assert.NotNil(t, w.Body)

			// Check the HTTP status code in the response
			if status := w.Code; status != tt.httpStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			paymentResp := models.PostPaymentResponse{}
			_ = json.NewDecoder(w.Body).Decode(&paymentResp)
			assert.Equal(t, tt.paymentStatus, paymentResp.PaymentStatus)

			if tt.paymentStatus == "Authorized" {
				// Check result is stored
				storedResult := pr.GetPayment(paymentResp.Id)
				assert.NotNil(t, storedResult)
			}
		})
	}
}
