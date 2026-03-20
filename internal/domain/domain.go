package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type CardNumber string

func NewCardNumber(cardNumber int) (CardNumber, error) {
	cnStr := strconv.Itoa(cardNumber)
	if len(cnStr) < 14 || len(cnStr) > 19 {
		return "", errors.New("invalid card number")
	}
	return CardNumber(cnStr), nil
}

func (cn *CardNumber) GetLastFour() int {
	cnStr := string(*cn)
	if len(cnStr) < 4 {
		return 0
	}

	lastFourStr := cnStr[len(cnStr)-4:]
	lastFour, err := strconv.Atoi(lastFourStr)
	if err != nil {
		return 0
	}

	return lastFour
}

type Amount int

func NewAmount(amount int) (Amount, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}
	return Amount(amount), nil
}

func (a *Amount) GetValue() int {
	return int(*a)
}

type CVV string

func NewCVV(cvv int) (CVV, error) {
	cvvStr := strconv.Itoa(cvv)
	isValid, err := regexp.Match(`^[0-9]{3,4}$`, []byte(cvvStr))
	if err != nil || !isValid {
		return "", errors.New("CVV must be 3-4 numbers long")
	}

	return CVV(cvvStr), nil
}

type Expiry struct {
	Month int
	Year  int
}

func NewExpiry(currentDate time.Time, expMonth, expYear int) (Expiry, error) {
	currentMonth := int(currentDate.Month())
	currentYear := int(currentDate.Year())

	if expMonth <= 0 || expMonth > 12 {
		return Expiry{}, errors.New("invalid expiry month")
	}

	if currentYear > expYear || (currentYear == expYear && currentMonth > expMonth) {
		return Expiry{}, errors.New("card has expired")
	}

	return Expiry{Month: expMonth, Year: expYear}, nil
}

func (e Expiry) String() string {
	return fmt.Sprintf("%d/%d", e.Month, e.Year)
}

type Currency string

const (
	GBP Currency = "GBP"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

func ParseCurrency(v string) (Currency, error) {
	c := Currency(v)
	switch c {
	case GBP, USD, EUR:
		return c, nil
	default:
		return "", errors.New("invalid currency")
	}
}

func (c *Currency) GetISO() string {
	return string(*c)
}

type Payment struct {
	CardNumber CardNumber
	Amount     Amount
	Cvv        CVV
	ExpiryDate Expiry
	Currency   Currency
}

func NewPayment(
	currentDate time.Time,
	cardNumber int,
	amount int,
	cvv int,
	expMonth int,
	expYear int,
	currency string,
) (Payment, error) {
	pmtCardNumber, err := NewCardNumber(cardNumber)
	if err != nil {
		return Payment{}, err
	}

	pmtAmount, err := NewAmount(amount)
	if err != nil {
		return Payment{}, err
	}

	pmtCVV, err := NewCVV(cvv)
	if err != nil {
		return Payment{}, err
	}

	pmtExpiry, err := NewExpiry(currentDate, expMonth, expYear)
	if err != nil {
		return Payment{}, err
	}

	pmtCurrency, err := ParseCurrency(currency)
	if err != nil {
		return Payment{}, err
	}

	return Payment{
		CardNumber: pmtCardNumber,
		Amount:     pmtAmount,
		Cvv:        pmtCVV,
		ExpiryDate: pmtExpiry,
		Currency:   pmtCurrency,
	}, nil
}
