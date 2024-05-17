package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/saratkumar-yb/infinityapi/config"
	"github.com/saratkumar-yb/infinityapi/db"
	"github.com/saratkumar-yb/infinityapi/router"
)

func main() {
	config.LoadConfig()

	migrateCmd := flag.Bool("migrate", false, "Migrate the database")
	startServerCmd := flag.Bool("startserver", false, "Start the API server")
	flag.Parse()

	if *migrateCmd {
		err := db.Migrate()
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration successful")
	} else if *startServerCmd {
		database, err := db.Connect()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer database.Close()

		r := router.NewRouter()
		address := fmt.Sprintf("%s:%d", config.AppConfig.HTTPListener, config.AppConfig.HTTPPort)
		log.Fatal(http.ListenAndServe(address, r))
	} else {
		log.Println("No command provided. Use -migrate or -startserver.")
	}
}
