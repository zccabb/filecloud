package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/boltdb/boltd"
)
//boltdbweb --db-name=./filecloud.db -p 8085
func main() {
	// Validate parameters.
	var path = "../filecloud.db"

	// Open the database.
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Enable logging.
	log.SetFlags(log.LstdFlags)

	// Setup the HTTP handlers.
	http.Handle("/", boltd.NewHandler(db))

	// Start the HTTP server.
	go func() { log.Fatal(http.ListenAndServe(":9000", nil)) }()

	fmt.Printf("Listening on http://localhost:%s\n", "9000")
	select {}
}
