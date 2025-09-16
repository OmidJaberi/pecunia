package main

import (
	"fmt"
	"os"
	
	//"github.com/jmoiron/sqlx"
	"github.com/OmidJaberi/pecunia/internal/db"
)

func main() {
	database := db.Connect("pecunia.db")

	migration, err := os.ReadFile("migrations/0001_init.sql")
	if err != nil {
		panic(err)
	}
	db.Migrate(database, string(migration))

	fmt.Println("Pecunia initialized!")
}
