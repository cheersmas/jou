package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/cheersmas/jou/database"
	"github.com/cheersmas/jou/repositories"
	"github.com/cheersmas/jou/services"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all journal entries",
	Long:  `Display all journal entries sorted by date`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Initialize database
		db, err := database.NewDatabase()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer db.Close()

		// Initialize repository
		journalRepo, err := repositories.NewJournalRepository(ctx, db)
		if err != nil {
			log.Fatalf("Failed to initialize repository: %v", err)
		}

		// Initialize service
		journalService := services.NewJournalService(journalRepo)

		// Get all journals
		journals, err := journalService.ListAll(ctx)
		if err != nil {
			log.Fatalf("Failed to retrieve journals: %v", err)
		}

		// Display journals
		for _, journal := range journals {
			fmt.Printf("ID: %d\n", journal.Id)
			fmt.Printf("Date: %s\n", journal.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Content: %s\n", journal.Content)
			fmt.Println("-------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
