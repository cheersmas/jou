package main

import (
	"fmt"
	"log"

	"github.com/cheersmas/jou/cmd"
	"github.com/cheersmas/jou/database"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		panic(err)
	}
	fmt.Print("wtf is going on")
	defer db.Close()

	cmd.Execute()
}
