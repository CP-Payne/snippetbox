package main

import (
	"errors"
	"fmt"

	//"html/template"
	"net/http"
	"strconv"

	"github.com/CP-Payne/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/base.tmpl.html", // Must be first
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }
	//
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// log.Println(err.Error())
	// 	// app.errorLog.Println(err.Error()) // app.serverError helper now prints this
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	app.serverError(w, err)
	// 	return
	// }
	//
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	// log.Println(err.Error())
	// 	// app.errorLog.Println(err.Error())
	// 	app.serverError(w, err)
	// 	// http.Error(w, "Internal Server Error", 500)
	// }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use not found helper
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write to the http.ResponseWriter
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// User r.Method to check whether the request is post or not
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy data
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
