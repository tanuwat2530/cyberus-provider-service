package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/provider-service/services"

	"net/http"
)

func MoFlowReceive(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := services.MoFlowReceiveProcessRequest(r)

	utilities.ResponseWithJSON(w, http.StatusOK, response)
}
