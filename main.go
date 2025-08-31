package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cheersmas/jou/database"
	"github.com/cheersmas/jou/repositories"
	"github.com/cheersmas/jou/services"
	"github.com/cheersmas/jou/tea"
)

func main() {
	ctx := context.Background()
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		panic(err)
	}
	fmt.Print("wtf is going on")
	defer db.Close()

	journalRepo, err := repositories.NewJournalRepository(ctx, db)
	if err != nil {
		log.Fatalf("Failed to initialize journal: %v", err)
	}
	journalService := services.NewJournalService(journalRepo)

	tea.Root(ctx, journalService)
}
