package main

import (
	"github.com/cheersmas/jou/cmd"
)

func main() {
	// ctx := context.Background()
	// db, err := sql.Open(DB_DRIVER_NAME, DB_DATA_SOURCE_NAME)
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
	// defer db.Close()

	// journalRepository, err := repositories.NewJournalRepository(ctx, db)
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }

	// journalService := services.NewJournalService(journalRepository)

	// newJournal := &domains.Journal{
	// 	Content: "so today was a great day for me",
	// }
	// id, err := journalService.Create(ctx, *newJournal)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// journal, err := journalService.Read(ctx, id)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("content %v", journal.Content)
	cmd.Execute()
}
