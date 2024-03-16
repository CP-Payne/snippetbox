package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// The below will return a pointer to the address value and not the value itself
	addr := flag.String("addr", ":4000", "HTTP network address")

	// If parse is not called then the addr flag will always be its default value
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Set a new http.server struct in order to make use of our new errorLog logger
	// instead of the standard logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
