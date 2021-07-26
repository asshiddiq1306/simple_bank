package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/asshiddiq1306/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQuery *Query
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQuery = NewQuery(testDB)

	os.Exit(m.Run())
}
