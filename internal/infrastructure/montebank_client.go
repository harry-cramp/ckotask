package infrastructure

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/service"

	"github.com/google/uuid"
)

type MontebankClient struct{}

func NewMontebankClient() *MontebankClient {
	mc := MontebankClient{}
	return &mc
}

type PaymentRequest struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Currency   string `json:"currency"`
	Amount     int    `json:"amount"`
	Cvv        string `json:"cvv"`
}

type PaymentResponse struct {
	Authorized        bool   `json:"authorized"`
	AuthorizationCode string `json:"authorization_code"`
}

func (m *MontebankClient) Charge(payment domain.Payment) (service.BankResponse, error) {
	url := "http://localhost:8080/payments"

	paymentReq := PaymentRequest{
		CardNumber: string(payment.CardNumber),
		ExpiryDate: payment.ExpiryDate.String(),
		Currency:   payment.Currency.GetISO(),
		Amount:     payment.Amount.GetValue(),
		Cvv:        string(payment.Cvv),
	}

	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return service.BankResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return service.BankResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return service.BankResponse{}, service.ErrUpstreamUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return service.BankResponse{}, service.ErrUpstreamUnavailable
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return service.BankResponse{}, err
	}

	status := "Declined"
	if resp.StatusCode == http.StatusOK && paymentResp.Authorized {
		status = "Authorized"
	}

	return service.BankResponse{
		PostPaymentResponse: models.PostPaymentResponse{
			Id:                 uuid.New().String(),
			PaymentStatus:      status,
			CardNumberLastFour: payment.CardNumber.GetLastFour(),
			ExpiryMonth:        payment.ExpiryDate.Month,
			ExpiryYear:         payment.ExpiryDate.Year,
			Currency:           payment.Currency.GetISO(),
			Amount:             payment.Amount.GetValue(),
		},
	}, nil
}
