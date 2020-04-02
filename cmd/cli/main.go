package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/emeli-frank/pick_go/pkg/database"
	"log"
	"os"
)

type application struct {
	DB               *sql.DB
	errorLog         *log.Logger
	infoLog          *log.Logger
}

func main() {
	dsn := flag.String("dsn", "pick:pick@/pick?parseTime=true&multiStatements=true", "MySQL database connection info")
	action := flag.String("action", "", "action")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := database.OpenDB(*dsn)
	if err != nil {
		fmt.Printf("database: An error occured, %s", err)
	}
	defer db.Close()

	app := application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		DB:               db,
	}

	switch *action {
	case "seed-core":
		if err := app.seedCore(); err != nil {
			app.errorLog.Fatal()
		}
		/*case "seed-mock":
		fmt.Println("creating mock data...")

		checkErr(app.mock())

		fmt.Println("mocking competed")*/
	}
}

func (a application) seedCore() error {
	scriptPaths := []string{
		"./pkg/storage/mysql/.db_setup/teardown.sql",
		"./pkg/storage/mysql/.db_setup/tables.sql",
		"./pkg/storage/mysql/.db_setup/data.sql",
	}

	err := database.ExecScripts(a.DB, scriptPaths...)
	if err != nil {
		panic(err)
	}

	return err
}

func (app *application) mock(db *sql.DB) error {
	var err error
	err = database.ExecScripts(
		db,
		"./pkg/storage/mysql/.db_setup/teardown.sql",
		"./pkg/storage/mysql/.db_setup/tables.sql",
		"./pkg/storage/mysql/.db_setup/data.sql")
	if err != nil {
		return err
	}



	return nil
}


func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}


