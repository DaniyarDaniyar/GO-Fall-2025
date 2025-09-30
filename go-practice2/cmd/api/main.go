package main

import (
	"fmt"
	"log"
	"net/http"

	"go-practice2/internal/handlers"
	"go-practice2/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/user", middleware.APIKey(http.HandlerFunc(handlers.UserHandler)))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
