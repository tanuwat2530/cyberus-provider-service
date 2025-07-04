package routes

import (
	"cyberus/provider-service/controllers"
	"net/http"
)

// SetupRoutes registers all application routes
func SetupRoutes() {
	// Register routes using http.HandleFunc
	http.HandleFunc("/ais/receive", controllers.AisMoFlowReceive)
	http.HandleFunc("/dtac/receive", controllers.TmvhMoFlowReceive)
	http.HandleFunc("/tmvh/receive", controllers.TmvhMoFlowReceive)
	http.HandleFunc("/api/", HomeHandler)
}

// HomeHandler for root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to CYBERUS-PROVIDER-SERVICE API power by GoLang ^_^"))
}
