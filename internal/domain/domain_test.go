package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCardNumber(t *testing.T) {
	tests := []struct {
		name      string
		cardNum   int
		lastFour  int
		shouldErr bool
	}{
		{"no card number", 0, 0, true},
		{"card number too short", 1000200030004, 0, true},
		{"shortest valid card number", 12341234123456, 3456, false},
		{"longest valid card number", 1234123412345656567, 6567, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardNumber, err := NewCardNumber(tt.cardNum)

			if tt.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, cardNumber.GetLastFour(), tt.lastFour)
			}
		})
	}
}

func TestNewAmount(t *testing.T) {
	tests := []struct {
		name      string
		amount    int
		shouldErr bool
	}{
		{"smallest valid amount", 1, false},
		{"large valid amount", 1000000, false},
		{"amount not valid", 0, true},
		{"negative amount", -1, true},
		{"large negative amount", -100000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amt, err := NewAmount(tt.amount)

			if tt.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.amount, amt.GetValue())
			}
		})
	}
}

func TestNewCVV(t *testing.T) {
	tests := []struct {
		name      string
		cvv       int
		shouldErr bool
	}{
		{"empty cvv", 0, true},
		{"cvv too short", 23, true},
		{"valid cvv (length 3)", 123, false},
		{"valid cvv (length 4)", 4567, false},
		{"cvv too long", 89012, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCVV(tt.cvv)

			if tt.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestExpiry(t *testing.T) {
	tests := []struct {
		name      string
		expMonth  int
		expYear   int
		dateStr   string
		shouldErr bool
	}{
		{"invalid month (0)", 0, 2099, "", true},
		{"invalid month (13)", 13, 2099, "", true},
		{"expired card (past year)", 5, 2025, "", true},
		{"expired card (current year)", 2, 2026, "", true},
		{"card not expired (current month and year)", 5, 2026, "05/2026", false},
		{"card not expired (current year)", 8, 2026, "08/2026", false},
		{"card not expired", 5, 2029, "05/2029", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
			_, err := NewExpiry(now, tt.expMonth, tt.expYear)

			if tt.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestParseCurrency(t *testing.T) {
	tests := []struct {
		name      string
		isoCode   string
		shouldErr bool
	}{
		{"invalid currency code", "ABC", true},
		{"unsupported currency code", "JPY", true},
		{"supported currency code (GBP)", "GBP", false},
		{"supported currency code (USD)", "USD", false},
		{"supported currency code (EUR)", "EUR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCurrency(tt.isoCode)

			if tt.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestNewPayment(t *testing.T) {
	tests := []struct {
		name       string
		cardNumber int
		amount     int
		cvv        int
		expMonth   int
		expYear    int
		currency   string
		shouldErr  bool
	}{
		{
			"failed to create payment (invalid card number)",
			0,
			100,
			123,
			10,
			2026,
			"USD",
			true,
		},
		{
			"failed to create payment (invalid amount)",
			1234567812345678,
			-1000,
			123,
			10,
			2026,
			"USD",
			true,
		},
		{
			"failed to create payment (invalid cvv)",
			1234567812345678,
			1000,
			12345,
			10,
			2026,
			"USD",
			true,
		},
		{
			"failed to create payment (expired card)",
			1234567812345678,
			1000,
			123,
			10,
			2023,
			"USD",
			true,
		},
		{
			"failed to create payment (unsupported currency)",
			1234567812345678,
			1000,
			123,
			10,
			2026,
			"JPY",
			true,
		},
		{
			"payment creation success",
			1234567812345678,
			1000,
			123,
			10,
			2026,
			"GBP",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
			_, err := NewPayment(
				now,
				tt.cardNumber,
				tt.amount,
				tt.cvv,
				tt.expMonth,
				tt.expYear,
				tt.currency,
			)

			if tt.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
