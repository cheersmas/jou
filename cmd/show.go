package cmd

import (
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all journal entries",
	Long:  `Display all journal entries sorted by date`,
	Run: func(cmd *cobra.Command, args []string) {
		// ctx := context.Background()

		// // Get the database instance
		// db := database.GetInstance()
		// if db == nil {
		// 	log.Fatal("Database connection not initialized")
		// }

		// // Initialize repository
		// journalRepo, err := repositories.NewJournalRepository(ctx, db)
		// if err != nil {
		// 	log.Fatalf("Failed to initialize repository: %v", err)
		// 	panic(err)
		// }

		// // Initialize service
		// journalService := services.NewJournalService(journalRepo)

		// // Get all journals
		// journals, err := journalService.ListAll(ctx)
		// if err != nil {
		// 	log.Fatalf("Failed to retrieve journals: %v", err)
		// 	panic(err)
		// }

		// // Display journals
		// for _, journal := range journals {
		// 	fmt.Printf("ID: %d\n", journal.Id)
		// 	fmt.Printf("Date: %s\n", journal.CreatedAt.Format("2006-01-02 15:04:05"))
		// 	fmt.Printf("Content: %s\n", journal.Content)
		// 	fmt.Println("-------------------")
		// }
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
