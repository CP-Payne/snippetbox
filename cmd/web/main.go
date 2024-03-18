package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/CP-Payne/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// The below will return a pointer to the address value and not the value itself
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define a new command line flag for the MYSQL DSN string.
	dsn := flag.String("dsn", "web:CHANGETHISPASSWORD@/snippetbox?parseTime=true", "MySQL data source name")

	// If parse is not called then the addr flag will always be its default value
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Created separate function for creating conenction pool to keep main tidy.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// Set a new http.server struct in order to make use of our new errorLog logger
	// instead of the standard logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes method to get the servemux containing our routes
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Use ping to create a connection and check for errors
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
