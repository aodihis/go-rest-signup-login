package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aodihis/go-rest-signup-login/config"
	"github.com/aodihis/go-rest-signup-login/database"
	"github.com/aodihis/go-rest-signup-login/internal/handlers"
)

func main() {

	config.LoadEnv()
	port := config.GetEnv("PORT")

	if err := database.ConnectDb(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	authHandler := handlers.NewAuthHandler()

	http.HandleFunc("/signup", authHandler.SignUp)
	http.HandleFunc("/login", authHandler.Login)

	if port == "" {
		port = "80"
	}
	fmt.Printf("Server run on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
