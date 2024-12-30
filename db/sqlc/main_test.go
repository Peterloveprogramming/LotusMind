package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/meditation?sslmode=disable"
)

// can be either transaction or connection
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	// create a new database connection
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
