package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"github.com/OmidJaberi/pecunia/internal/domain"
)

// UserRepo
type UserRepo struct{ db *sqlx.DB }

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(u *domain.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, name, created_at) VALUES (?, ?, ?)`,
		u.ID.String(), u.Name, u.CreatedAt.Unix())
	return err
}

func (r *UserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var u domain.User
	var createdAt int64
	err := r.db.QueryRowx(`SELECT id, name, created_at FROM users WHERE id = ?`, id.String()).Scan(&u.ID, &u.Name, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CurrencyRepo
type CurrencyRepo struct { db *sqlx.DB  }

func NewCurrencyRepo(db *sqlx.DB) *CurrencyRepo { return &CurrencyRepo{db: db} }

func (r *CurrencyRepo) Insert(c domain.Currency) error {
	_, err := r.db.Exec(`
		INSERT INTO currencies (code, name, symbol, decimals)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(code) DO NOTHING`, c.Code, c.Name, c.Symbol, c.Decimals,
	)
	return err
}

func (r *CurrencyRepo) List() ([]domain.Currency, error) {
	var list []domain.Currency
	err := r.db.Select(&list, `SELECT code, name, symbol, decimals FROM currencies`)
	return list, err
}

// AssetRepo
type AssetRepo struct { db *sqlx.DB  }

func NewAssetRepo(db *sqlx.DB) *AssetRepo { return &AssetRepo{db: db} }

func (r *AssetRepo) Insert(a domain.Asset) error {
	_, err := r.db.Exec(`
		INSERT INTO assets (id, user_id, name, currency_code, amount, category, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.UserID, a.Name, a.Value.Currency.Code, a.Value.Amount, a.Category, a.CreatedAt.Unix(),
	)
	return err
}

func (r *AssetRepo) ListByUserID(userID uuid.UUID) ([]domain.Asset, error) {
	rows, err := r.db.Queryx(`
		SELECT id, user_id, name, currency_code, amount, category, created_at
		FROM assets WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Asset
	for rows.Next() {
		var (
			id			uuid.UUID
			uid			uuid.UUID
			name		string
			code		string
			amount		decimal.Decimal
			category	string
			createdAt	int64
		)
		if err := rows.Scan(&id, &uid, &name, &code, &amount, &category, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, domain.Asset{
			ID:			id,
			UserID:		uid,
			Name:		name,
			Value:		domain.Money{
				Amount:		amount,
				Currency:	domain.Currency{Code: code}, // Not filled for now
			},
			Category:	category,
			CreatedAt:	time.Unix(createdAt, 0),
		})
	}
	return result, nil
}

// ExchangeRateRepo
type ExchangeRateRepo struct{ db *sqlx.DB }

func NewExchangeRateRepo(db *sqlx.DB) *ExchangeRateRepo { return &ExchangeRateRepo{db: db} }

func (r *ExchangeRateRepo) Upsert(er domain.ExchangeRate) error {
	_, err := r.db.Exec(`
		INSERT INTO exchange_rates (user_id, from_currency, to_currency, rate)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, from_currency, to_currency)
		DO UPDATE SET rate = excluded.rate`,
		er.UserID, er.From, er.To, er.Rate,
	)
	return err
}

func (r *ExchangeRateRepo) ListByUser(userID uuid.UUID) ([]domain.ExchangeRate, error) {
	var list []domain.ExchangeRate
	err := r.db.Select(&list, `
		SELECT user_id AS "userid", from_currency AS "from", to_currency AS "to", rate
		FROM exchange_rates WHERE user_id = ?`, userID)
	return list, err
}
