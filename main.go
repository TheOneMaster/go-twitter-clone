package main

import (
	"fmt"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const PortNumber = 8080

func main() {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("static"))

	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", handlers.IndexHandler)

	fmt.Printf("Running on http://localhost:%d\n", PortNumber)
	http.ListenAndServe(fmt.Sprintf(":%d", PortNumber), r)
}
