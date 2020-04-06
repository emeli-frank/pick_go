package main

import (
	"bufio"
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
	scriptPath := flag.String("scriptPath", "./scripts", "Scripts path")
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
	default :
		fmt.Println("Creating tables and seeding core data")
		if err := app.seedCore(*scriptPath); err != nil {
			app.errorLog.Fatal()
		}
		fmt.Println("Done!")
	}

	fmt.Println()
	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (a application) seedCore(scriptPath string) error {
	//scriptPath = "./pkg/storage/mysql/.db_setup"
	//fmt.Println(fmt.Sprintf("%s/teardown.sql", scriptPath))
	scriptPaths := []string{
		fmt.Sprintf("%s/teardown.sql", scriptPath),
		fmt.Sprintf("%s/tables.sql", scriptPath),
		fmt.Sprintf("%s/data.sql", scriptPath),
	}

	err := database.ExecScripts(a.DB, scriptPaths...)
	if err != nil {
		panic(err)
	}

	return err
}

func (a *application) mock(db *sql.DB) error {
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


