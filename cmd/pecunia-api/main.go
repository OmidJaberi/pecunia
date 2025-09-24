package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OmidJaberi/pecunia/internal/db"
	"github.com/OmidJaberi/pecunia/internal/domain"
	"github.com/OmidJaberi/pecunia/internal/server"
)

func main() {
	dbPath := "pecunia.db"
	database := db.Connect(dbPath)
	defer database.Close()

	migration, err := os.ReadFile("migrations/0001_init.sql")
	if err != nil {
		log.Fatalf("read migration: %v", err)
	}
	db.Migrate(database, string(migration))

	// seed a few currencies (idempotent)
	cRepo := db.NewCurrencyRepo(database)
	_ = cRepo.Insert(domain.Currency{Code: "USD", Name: "US Dollar", Symbol: "$", Decimals: 2})
	_ = cRepo.Insert(domain.Currency{Code: "IRR", Name: "Iranian Rial", Symbol: "﷼", Decimals: 0})
	_ = cRepo.Insert(domain.Currency{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", Decimals: 2})
	_ = cRepo.Insert(domain.Currency{Code: "BTC", Name: "Bitcoin", Symbol: "₿", Decimals: 8})
	_ = cRepo.Insert(domain.Currency{Code: "grGOLD18", Name: "18k Gold Gram", Symbol: "g", Decimals: 3})

	api := server.NewAPI(database)
	httpAddr := ":8080"
	fmt.Println("starting pecunia api on", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, server.Routes(api)))
}
