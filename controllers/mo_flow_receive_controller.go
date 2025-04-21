package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/provider-service/services"
	"fmt"

	"net/http"
)

func MoFlowReceive(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := services.MoFlowReceiveProcessRequest(r)
	// fmt.Println(response["code"])
	// fmt.Println(response["partner_id"])
	// fmt.Println(response["refid"])
	// fmt.Println(response["adsid"])
	//utilities.ResponseWithJSON(w, http.StatusOK, response)

	// Build redirect URL with query params
	if response["code"] == "302" {
		target := fmt.Sprintf("https://portal.mmtcontent.com/landing?id=%s&refid=%s&media=%s", response["partner_id"], response["refid"], response["adsid"])
		// Send redirect (302 by default)
		http.Redirect(w, r, target, http.StatusFound)
	}
	if response["code"] == "-1" {
		utilities.ResponseWithJSON(w, http.StatusOK, response)
	}

}
