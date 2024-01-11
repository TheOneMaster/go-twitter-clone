package main

import (
	"fmt"
	"net/http"

	"github.com/TheOneMaster/go-twitter-clone/db"
	"github.com/TheOneMaster/go-twitter-clone/env"
	"github.com/TheOneMaster/go-twitter-clone/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	env.LoadEnv()

	db.InitConnection(env.Environment.DatabaseUrl)

	portNumber := env.Environment.PortNumber

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(handlers.ValidateSession)
	// r.Use(handlers.Validate)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Static Routes
	r.Get("/", handlers.IndexPage)
	r.Get("/login", handlers.LoginPage)
	r.Get("/register", handlers.RegisterPage)
	r.Get("/humans.txt", handlers.HumansHandler)
	r.Get("/logout", handlers.LogOut)

	r.Post("/login", handlers.LoginRequest)
	r.Post("/register", handlers.RegisterRequest)
	r.Post("/sendMessage", handlers.MessageHandler)

	// Dynamic routes
	r.Get("/message/*", handlers.GetMessage)
	r.Get("/replies/*", handlers.ReplyHandler)

	r.Post("/like/*", handlers.LikeMessage)

	fmt.Printf("Running on http://localhost:%d\n", portNumber)
	http.ListenAndServe(fmt.Sprintf(":%d", portNumber), r)
}
