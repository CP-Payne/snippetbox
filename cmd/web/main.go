package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// The below will return a pointer to the address value and not the value itself
	addr := flag.String("addr", ":4000", "HTTP network address")

	// If parse is not called then the addr flag will always be its default value
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
