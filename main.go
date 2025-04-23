package main

import (
	"cyberus/provider-service/routes"
	"fmt"
	"net/http"
)

func main() {
	// Setup all routes
	routes.SetupRoutes()

	// Start the server on port 8080
	fmt.Println("Starting cyberus-provider-service server on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
