package repository

import (
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/service"
)

type PaymentsRepository struct {
	payments []service.BankResponse
}

func NewPaymentsRepository() *PaymentsRepository {
	return &PaymentsRepository{
		payments: []service.BankResponse{},
	}
}

func (ps *PaymentsRepository) GetPayment(id string) *service.BankResponse {
	for _, element := range ps.payments {
		if element.Id == id {
			return &element
		}
	}
	return nil
}

func (ps *PaymentsRepository) AddPayment(payment service.BankResponse) {
	ps.payments = append(ps.payments, payment)
}
