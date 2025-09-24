package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Currency struct {
	Code		string
	Name		string
	Symbol		string
	Decimals	int
}

type Money struct {
	Amount		decimal.Decimal
	Currency	Currency
}

type Asset struct {
	ID			uuid.UUID
	UserID		uuid.UUID
	Name		string
	Value		Money
	Category	string
	CreatedAt	time.Time
}

type User struct {
	ID			uuid.UUID
	Name		string
	CreatedAt	time.Time
}

type Transaction struct {
	ID			uuid.UUID
	UserID		uuid.UUID
	Description	string
	Amount		Money
	Frequency	string
	StartDate	time.Time
	EndDate		time.Time
	CreatedAt	time.Time
}

type ExchangeRate struct {
	UserID	uuid.UUID
	From	string
	To		string
	Rate	decimal.Decimal
}
