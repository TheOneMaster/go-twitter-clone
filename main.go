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
	r.Use(middleware.Logger)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// GET Routes
	r.Get("/", handlers.IndexHandler)
	r.Get("/humans.txt", handlers.HumansHandler)

	// POST Routes

	fmt.Printf("Running on http://localhost:%d\n", PortNumber)
	http.ListenAndServe(fmt.Sprintf(":%d", PortNumber), r)
}
