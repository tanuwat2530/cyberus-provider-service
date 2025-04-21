package services

import (
	"CyberusGolangShareLibrary/redis_db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Struct to map the expected JSON fields
// type MoFlowReceiveRequest struct {
// 	AgencyId  string `json:"agency_id"`
// 	PartnerId string `json:"partner_id"`
// 	RefId     string `json:"refid"`
// 	AdsId     string `json:"adsid"`
// }

func MoFlowReceiveProcessRequest(r *http.Request) map[string]string {

	// var payload map[string]interface{}
	res := map[string]string{}

	// // Get current time
	now := time.Now()
	// // Unix timestamp in nanoseconds
	timestamp := (now.UnixNano())
	nano_timestamp := strconv.FormatInt(timestamp, 10)

	// // Generate a random UUID (UUID v4)
	transaction_id := uuid.New().String()

	// // Get a Client IP address
	client_ip := r.RemoteAddr

	agency_id := r.URL.Query().Get("agency_id")
	partner_id := r.URL.Query().Get("partner_id")
	refid := r.URL.Query().Get("refid")
	adsid := r.URL.Query().Get("adsid")

	if agency_id == "" || partner_id == "" || refid == "" || adsid == "" {
		res["code"] = "-1"
		res["message"] = "Invalid param"
		return res
	}

	// Parse query parameters into a map
	payload := map[string]string{
		"agency_id":  agency_id,
		"partner_id": partner_id,
		"refid":      refid,
		"adsid":      adsid,
	}

	// ✅ Add a new JSON node (e.g., tracking_source)
	payload["client_ip"] = client_ip
	payload["timestamp"] = nano_timestamp
	payload["transaction_id"] = transaction_id

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		//http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
		res["code"] = "-1"
		res["message"] = "Json error"
		res["transaction_id"] = transaction_id
		return res
	}

	payloadString := string(jsonData)
	//fmt.Println(payloadString)

	redis_db.ConnectRedis()
	redis_key := "MO:" + partner_id + ":" + transaction_id

	ttl := 240 * time.Hour // expires in 240 Hour

	// Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, payloadString, ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v", err)
	}
	//fmt.Println("Key set successfully with TTL")

	// errPayload := json.NewDecoder(r.Body).Decode(&payload)
	// if errPayload != nil {
	// 	// Example: print the values
	// 	//fmt.Println("Error decode Json to map[string]interface{} :", errPayload.Error())
	// 	res["code"] = "-1"
	// 	res["desc"] = "JSON Decode error"
	// 	res["ref-id"] = "undefined"
	// 	res["tran-ref"] = "undefined"
	// 	res["media"] = "undefined"
	// 	return res
	// }

	// jsonData, err := json.MarshalIndent(payload, "", "  ")
	// if err != nil {
	// 	//fmt.Println("Error marshalling JSON:", err.Error())
	// 	res["code"] = "-2"
	// 	res["desc"] = "JSON Marshalling error"
	// 	res["ref-id"] = "undefined"
	// 	res["tran-ref"] = "undefined"
	// 	res["media"] = "undefined"
	// 	return res
	// }

	// // // Unmarshal JSON into struct
	// var requestData MoFlowReceiveRequest
	// err = json.Unmarshal(jsonData, &requestData)
	// if err != nil {
	// 	//fmt.Println("Error map Json to Struct :" + err.Error())
	// 	//fmt.Println("Error marshalling JSON:", err.Error())
	// 	res["code"] = "-3"
	// 	res["desc"] = "JSON Struct error"
	// 	res["ref-id"] = "undefined"
	// 	res["tran-ref"] = "undefined"
	// 	res["media"] = "undefined"
	// 	return res
	// }

	// // // Add ClientIP to Payload
	// payload["client_ip"] = ip
	// payload["timestamp"] = nano_timestamp
	// payload["transaction_id"] = transaction_id
	// payload["agency_id"] = requestData.AgencyId
	// payload["partner_id"] = requestData.PartnerId
	// payload["ads_id"] = requestData.AdsId
	// payload["ref_id"] = requestData.RefId
	// // Convert the struct to JSON string
	// payloadBytes, err := json.Marshal(payload)
	// if err != nil {
	// 	fmt.Printf("Failed to convert payload to JSON:  %+v\n ", http.StatusInternalServerError)
	// 	fmt.Println("Error marshalling JSON:", err.Error())
	// 	res["code"] = "-4"
	// 	res["desc"] = "Additional value error"
	// 	return res
	// }

	// payloadString := string(payloadBytes)
	// fmt.Println(payloadString)
	// redis_db.ConnectRedis()
	// redis_key := "MO:" + requestData.PartnerId + ":" + transaction_id

	// ttl := 240 * time.Hour // expires in 240 Hour

	// // Set key with TTL
	// if err := redis_db.SetWithTTL(redis_key, payloadString, ttl); err != nil {
	// 	//write to file if Redis problem or forward request to AIS
	// 	log.Fatalf("SetWithTTL error: %v", err)
	// }
	// fmt.Println("Key set successfully with TTL")
	res["code"] = "302"
	res["partner_id"] = partner_id
	res["refid"] = refid
	res["adsid"] = adsid

	defer r.Body.Close()

	return res
}
