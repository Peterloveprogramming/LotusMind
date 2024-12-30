package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/lotusMind/meditation/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/meditation?sslmode=disable"
)

// can be either transaction or connection
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	//load config from where app.env is located with respective to the current file
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can not load config", err)
	}

	// create a new database connection
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
