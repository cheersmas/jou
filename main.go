package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/cheersmas/jou/domains"
	"github.com/cheersmas/jou/repositories"
	"github.com/cheersmas/jou/services"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_DRIVER_NAME      = "sqlite3"
	DB_DATA_SOURCE_NAME = "lite.db"
)

func main() {
	ctx := context.Background()
	db, err := sql.Open(DB_DRIVER_NAME, DB_DATA_SOURCE_NAME)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer db.Close()

	journalRepository, err := repositories.NewJournalRepository(ctx, db)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	journalService := services.NewJournalService(journalRepository)

	newJournal := &domains.Journal{
		Content: "so today was a great day for me",
	}
	id, err := journalService.Create(ctx, *newJournal)
	if err != nil {
		log.Fatal(err)
	}

	journal, err := journalService.Read(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("content %v", journal.Content)
}
