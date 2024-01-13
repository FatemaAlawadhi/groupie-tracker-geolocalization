package Handlers

import (
	"html/template"
	"net/http"
	"groupie-tracker-filter/Error"
)

var err error


func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/about" && r.URL.Path != "/stars" {
		Error.RenderErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("WebPages/Home.html")
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		}
	} else {
		Error.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
}

func HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("WebPages/About.html")
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		}
	} else {
		Error.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
}
