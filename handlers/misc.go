package handlers

import (
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/templates"
)

func HumansHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/humans.txt")
}

func PageNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404"))
}

func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}

func ServeStaticPage(pageFile string, w http.ResponseWriter, r *http.Request) {
	t, err := templates.LoadFiles("base.html", pageFile)
	if err != nil {
		PageNotFound(w)
	}

	err = t.Execute(w, nil)
	if err != nil {
		ServerError(w)
	}
}

func ServeFragment(w http.ResponseWriter, fragment string, data any) {
	t, err := templates.LoadFragment(fragment)
	if err != nil {
		ServerError(w)
	}
	t.Execute(w, data)
}
