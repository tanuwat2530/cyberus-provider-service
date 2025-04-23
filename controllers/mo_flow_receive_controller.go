package controllers

import (
	"CyberusGolangShareLibrary/utilities"
	services "cyberus/provider-service/services"
	"fmt"
	"io/ioutil"

	"net/http"
)

func MoFlowReceive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call controller")
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

	response := services.MoFlowReceiveProcessRequest(agency_id, partner_id, refid, adsid, client_ip)
	// fmt.Println(response["code"])
	// fmt.Println(response["partner_id"])
	// fmt.Println(response["refid"])
	// fmt.Println(response["adsid"])
	//utilities.ResponseWithJSON(w, http.StatusOK, response)

	// Build redirect URL with query params
	if response["code"] == "302" {
		url := "https://portal.mmtcontent.com/landing?id" + response["partner_id"] + "s&refid=" + response["refid"] + "&media=" + response["adsid"]
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		// Display it in a simple HTML template
		fmt.Fprintf(w, `%s`, body)

	}
	if response["code"] == "-1" {
		utilities.ResponseWithJSON(w, http.StatusOK, response)
	}

	if response["code"] == "0" {
		utilities.ResponseWithJSON(w, http.StatusOK, response)
	}
	defer r.Body.Close()
}
