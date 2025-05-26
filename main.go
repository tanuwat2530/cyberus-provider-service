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
	fmt.Println("Starting cyberus-provider-service server on port 4001...")
	if err := http.ListenAndServe(":4001", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
