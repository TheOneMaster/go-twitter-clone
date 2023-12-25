package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	wd, err := os.Getwd()
	if err != nil {
		ServerError(w)
	}

	base_layout := filepath.Join(wd, "templates/layouts/base.html")
	index_layout := filepath.Join(wd, "templates/layouts/index.html")

	t, err := template.ParseFiles(base_layout, index_layout)
	if err != nil {
		fmt.Println(err)
		PageNotFound(w)
	}

	err = t.Execute(w, nil)
	if err != nil {
		ServerError(w)
	}

}

func PageNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404"))
}

func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}
