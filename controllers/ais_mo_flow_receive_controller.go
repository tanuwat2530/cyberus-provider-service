package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/provider-service/services"

	"net/http"
)

func AisMoFlowReceive(w http.ResponseWriter, r *http.Request) {
	//redirectUrl := os.Getenv("BN_DB_URL")
	// Check if the method is POST
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client_ip := r.RemoteAddr
	agency_id := r.URL.Query().Get("agency_id")
	partner_id := r.URL.Query().Get("partner_id")
	refid := r.URL.Query().Get("refid")
	adsid := r.URL.Query().Get("adsid")

	response := services.AisMoFlowReceiveProcessRequest(agency_id, partner_id, refid, adsid, client_ip)

	// Build redirect URL with query params
	if response["code"] == "302" {
		url := "http://www.bigfunarea.com/?id=" + response["shortcode"] + "&refid=" + response["refid"] + "&media=" + response["adsid"]
		// For demonstration, let's use 302 Found (temporary redirect)
		http.Redirect(w, r, url, http.StatusFound)

	}
	if response["code"] == "-1" {
		utilities.ResponseWithJSON(w, http.StatusOK, response)
	}

	if response["code"] == "0" {
		utilities.ResponseWithJSON(w, http.StatusOK, response)
	}
	defer r.Body.Close()
}
