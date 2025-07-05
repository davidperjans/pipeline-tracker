package main

import (
	"fmt"
	"log"

	"github.com/davidperjans/pipeline-tracker/internal/storage"
)

func main() {
	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
}
