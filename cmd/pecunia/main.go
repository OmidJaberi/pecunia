package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/OmidJaberi/pecunia/internal/db"
	"github.com/OmidJaberi/pecunia/internal/domain"
	"github.com/OmidJaberi/pecunia/internal/exchangegraph"
)

func main() {
	database := db.Connect("pecunia.db")
	migration, err := os.ReadFile("migrations/0001_init.sql")
	if err != nil {
		panic(err)
	}
	db.Migrate(database, string(migration))

	// Setup repos
	currencies := db.NewCurrencyRepo(database)
	assets := db.NewAssetRepo(database)
	rates := db.NewExchangeRateRepo(database)

	// Seed currencies
	seedCurrencies(currencies)

	// Test user
	userID := uuid.New()

	// Seed exchange rates
	rates.Upsert(domain.ExchangeRate{UserID: userID, From: "USD", To: "IRR", Rate: decimal.NewFromInt(990000)})
	rates.Upsert(domain.ExchangeRate{UserID: userID, From: "BTC", To: "USD", Rate: decimal.NewFromInt(60000)})
	rates.Upsert(domain.ExchangeRate{UserID: userID, From: "CNY", To: "USD", Rate: decimal.NewFromInt(7)})

	assets.Insert(domain.Asset{
		ID:			uuid.New(),
		UserID:		userID,
		Name:		"Bit-CoinSack",
		Value:		domain.Money{
			Amount:		decimal.NewFromFloat(1.33),
			Currency:	domain.Currency{Code: "BTC"},
		},
		Category:	"investment",
		CreatedAt:	time.Now(),
	})

	assets.Insert(domain.Asset{
		ID:			uuid.New(),
		UserID:		userID,
		Name:		"Cash",
		Value:		domain.Money{
			Amount:		decimal.NewFromFloat(753),
			Currency:	domain.Currency{Code: "USD"},
		},
		Category:	"investment",
		CreatedAt:	time.Now(),
	})

	// Load data and compute totals
	exRates, _ := rates.ListByUser(userID)
	g := exchangegraph.NewCurrencyGraph(exRates)

	allAssets, _ := assets.ListByUserID(userID)
	totalUSD := decimal.Zero
	for _, a := range allAssets {
		val, err := g.Convert(a.Value.Amount, a.Value.Currency.Code, "USD")
		if err != nil {
			fmt.Println("conversion error:", err)
			continue
		}
		fmt.Printf("%s: %s %s\n", a.Name, a.Value.Amount.String(), a.Value.Currency.Code)
		totalUSD = totalUSD.Add(val)
	}
	fmt.Println("Total assets in USD:", totalUSD.String())
}

func seedCurrencies(repo *db.CurrencyRepo) {
	defaults := []domain.Currency{
		{Code: "USD", Name: "US Dollar", Symbol: "$", Decimals: 2},
		{Code: "IRR", Name: "Iranian Rial", Symbol: "﷼", Decimals: 0},
		{Code: "CNY", Name: "US Dollar", Symbol: "¥", Decimals: 2},
		{Code: "BTC", Name: "Bitcoin", Symbol: "B", Decimals: 8},
		{Code: "grGOLD18", Name: "18k Gold Gram", Symbol: "g", Decimals: 3},
	}
	for _, c := range defaults {
		_ = repo.Insert(c)
	}
}
