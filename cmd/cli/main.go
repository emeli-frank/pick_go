package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"github.com/emeli-frank/pick_go/pkg/database"
	"github.com/emeli-frank/pick_go/pkg/domain/product"
	"github.com/emeli-frank/pick_go/pkg/storage/mysql"
	"log"
	"os"
	"syreclabs.com/go/faker"
)

type application struct {
	DB               *sql.DB
	errorLog         *log.Logger
	infoLog          *log.Logger
	productService   product.Service
}

func main() {
	dsn := flag.String("dsn", "pick:pick@/pick?parseTime=true&multiStatements=true", "MySQL database connection info")
	action := flag.String("action", "", "action")
	scriptPath := flag.String("scriptPath", "./scripts", "Scripts path")
	flag.Parse()

	db, err := database.OpenDB(*dsn)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	productRepo := mysql.NewProductStorage(db)
	productService := product.New(productRepo)

	if err != nil {
		fmt.Printf("database: An error occured, %s", err)
	}
	defer db.Close()

	app := application{
		errorLog:         errorLog,
		infoLog:          infoLog,
		DB:               db,
		productService: productService,
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
	scriptPath = "./pkg/storage/mysql/.db_setup"
	scriptPaths := []string{
		fmt.Sprintf("%s/teardown.sql", scriptPath),
		fmt.Sprintf("%s/tables.sql", scriptPath),
		fmt.Sprintf("%s/data.sql", scriptPath),
	}

	err := database.ExecScripts(a.DB, scriptPaths...)
	if err != nil {
		return err
	}

	err = a.mock()
	if err != nil {
		return err
	}

	err = a.mockProducts()
	if err != nil {
		return err
	}

	return nil
}

func (a *application) mock() error {
	var err error
	err = database.ExecScripts(
		a.DB,
		"./pkg/storage/mysql/.db_setup/teardown.sql",
		"./pkg/storage/mysql/.db_setup/tables.sql",
		"./pkg/storage/mysql/.db_setup/data.sql")
	if err != nil {
		return err
	}

	return nil
}

func (a *application) mockProducts() error {
	price := faker.RandomInt(50, 500)
	for i := 0; i < 100; i++ {
		p := product.Product{
			Name: faker.Commerce().ProductName(),
			Description: faker.Lorem().Paragraph(5),
			Quantity: faker.RandomInt(0, 100),
			RegularPrice: float32(price),
			DiscountPrice: float32(faker.RandomInt(50, price)),
		}

		_, err := a.productService.CreateProduct(&p)
		if err != nil {
			return err
		}
	}

	return nil
}


func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}


