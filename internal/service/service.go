package service

import (
	"errors"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/domain"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
)

var ErrUpstreamUnavailable = errors.New("upstream unavailable")

type BankResponse struct {
	models.PostPaymentResponse
}

type BankClient interface {
	Charge(payment domain.Payment) (BankResponse, error)
}

type PaymentRepository interface {
	GetPayment(id string) *BankResponse
	AddPayment(payment BankResponse)
}

type PaymentService struct {
	bank BankClient
}

func NewPaymentService(
	bank BankClient,
) *PaymentService {
	ps := PaymentService{bank: bank}
	return &ps
}

type PaymentParams struct {
	CardNumber int
	Amount     int
	Cvv        int
	ExpMonth   int
	ExpYear    int
	Currency   string
}

func (s *PaymentService) Process(p PaymentParams) (BankResponse, error) {
	payment, err := domain.NewPayment(
		time.Now(),
		p.CardNumber,
		p.Amount,
		p.Cvv,
		p.ExpMonth,
		p.ExpYear,
		p.Currency,
	)
	if err != nil {
		return buildRejectedResponse(p), err
	}

	bankResp, err := s.bank.Charge(payment)
	if err != nil {
		return BankResponse{}, err
	}

	return bankResp, nil
}

func buildRejectedResponse(p PaymentParams) BankResponse {
	lastFour := p.CardNumber % 10000

	bankResp := BankResponse{
		PostPaymentResponse: models.PostPaymentResponse{
			Id:                 "",
			PaymentStatus:      "Rejected",
			CardNumberLastFour: lastFour,
			ExpiryMonth:        p.ExpMonth,
			ExpiryYear:         p.ExpYear,
			Currency:           p.Currency,
			Amount:             p.Amount,
		},
	}

	return bankResp
}
