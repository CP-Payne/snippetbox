package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// The below will return a pointer to the address value and not the value itself
	addr := flag.String("addr", ":4000", "HTTP network address")

	// If parse is not called then the addr flag will always be its default value
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
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
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
