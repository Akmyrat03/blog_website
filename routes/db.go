package routes

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "user_management"
)

var db *sql.DB

// InitDB initializes the database connection and assigns it to the global db variable
func InitDB() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Verify that the connection is actually established
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	log.Println("Connected to the database!")
}
