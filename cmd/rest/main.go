package main

import (
	"flag"
	"fmt"
	"github.com/emeli-frank/pick_go/pkg/database"
	"github.com/emeli-frank/pick_go/pkg/domain/product"
	"github.com/emeli-frank/pick_go/pkg/domain/user"
	http2 "github.com/emeli-frank/pick_go/pkg/http"
	"github.com/emeli-frank/pick_go/pkg/storage/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	//infoLog *log.Logger
}

func main() {
	dsn := flag.String("dsn", "pick:pick@/pick?parseTime=true", "MySQL database connection info")
	addr := flag.String("addr", ":4242", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := database.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := application{
		//infoLog:     infoLog,
	}
	_ = app


	response := http2.NewResponse(errorLog)

	userRepo := mysql.NewUserStorage(db)
	var userService = user.New(userRepo)

	productRepo := mysql.NewProductStorage(db)
	productService := product.New(productRepo)

	server := http2.NewServer(response, userService, productService, infoLog)

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: server.Routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server is now running at localhost%s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
