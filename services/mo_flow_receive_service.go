package services

import (
	"CyberusGolangShareLibrary/redis_db"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// Struct to map the expected JSON fields
type MoFlowReceiveRequest struct {
	AgencyId  string `json:"agency_id"`
	PartnerId string `json:"partner_id"`
	RefId     string `json:"refid"`
	AdsId     string `json:"adsid"`
}

func MoFlowReceiveProcessRequest(r *http.Request) map[string]string {

	res := map[string]string{}
	var payload map[string]interface{}
	// Get current time
	now := time.Now()
	// Unix timestamp in nanoseconds
	timestamp := (now.UnixNano())
	nano_timestamp := strconv.FormatInt(timestamp, 10)

	// Generate a random UUID (UUID v4)
	transaction_id := uuid.New().String()

	// Get a Client IP address
	ip := r.RemoteAddr

	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		// Example: print the values
		//fmt.Println("Error decode Json to map[string]interface{} :", errPayload.Error())

		res["code"] = "-1"
		res["desc"] = "JSON Decode error"
		res["ref-id"] = "undefined"
		res["tran-ref"] = "undefined"
		res["media"] = "undefined"
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		//fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-2"
		res["desc"] = "JSON Marshalling error"
		res["ref-id"] = "undefined"
		res["tran-ref"] = "undefined"
		res["media"] = "undefined"
		return res
	}

	// // Unmarshal JSON into struct
	var requestData MoFlowReceiveRequest
	err = json.Unmarshal(jsonData, &requestData)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		//fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-3"
		res["desc"] = "JSON Struct error"
		res["ref-id"] = "undefined"
		res["tran-ref"] = "undefined"
		res["media"] = "undefined"
		return res
	}

	// // Add ClientIP to Payload
	payload["client_ip"] = ip
	payload["timestamp"] = nano_timestamp
	payload["transaction_id"] = transaction_id
	payload["agency_id"] = requestData.AgencyId
	payload["partner_id"] = requestData.PartnerId
	payload["ads_id"] = requestData.AdsId
	payload["ref_id"] = requestData.RefId
	// Convert the struct to JSON string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to convert payload to JSON:  %+v\n ", http.StatusInternalServerError)
		fmt.Println("Error marshalling JSON:", err.Error())
		res["code"] = "-4"
		res["desc"] = "Additional value error"
		return res
	}

	payloadString := string(payloadBytes)
	fmt.Println(payloadString)
	redis_db.ConnectRedis()
	redis_key := "MO:" + requestData.PartnerId + ":" + transaction_id

	ttl := 1 * time.Hour // expires in 1 Hour

	// Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, payloadString, ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v", err)
	}
	fmt.Println("Key set successfully with TTL")

	// Get the key
	// val, err := redis_db.GetValue(key)
	// if err != nil {
	// 	log.Printf("GetValue error: %v", err)
	// } else {
	// 	fmt.Printf("Retrieved value: %s\n", val)
	// }

	//redis_db.Set("aaa", "AAA", 300)

	res["code"] = "200"
	res["message"] = "OK"
	res["transaction_id"] = transaction_id

	defer r.Body.Close()

	return res
}
