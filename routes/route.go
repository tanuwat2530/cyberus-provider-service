package routes

import (
	"cyberus/provider-service/controllers"
	"net/http"
)

// SetupRoutes registers all application routes
func SetupRoutes() {
	// Register routes using http.HandleFunc
	http.HandleFunc("/api/mo-flow-receive", controllers.MoFlowReceive)
	http.HandleFunc("/api/", HomeHandler)
}

// HomeHandler for root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to backend API power by GoLang ^_^"))
}
