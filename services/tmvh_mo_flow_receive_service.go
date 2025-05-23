package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"CyberusGolangShareLibrary/redis_db"
	"cyberus/provider-service/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func TmvhMoFlowReceiveProcessRequest(agency_id string, partner_id string, refid string, adsid string, client_ip string) map[string]string {
	// var payload map[string]interface{}
	redisConnection := os.Getenv("BN_REDIS_URL")
	dbConnection := os.Getenv("BN_DB_URL")
	res := map[string]string{}

	// // Get current time
	now := time.Now()
	// // Unix timestamp in nanoseconds
	timestamp := (now.UnixNano())
	nano_timestamp := strconv.FormatInt(timestamp, 10)

	// // Generate a random UUID (UUID v4)
	transaction_id := uuid.New().String()

	if agency_id == "" || partner_id == "" || refid == "" || adsid == "" {
		fmt.Println("Invalid param")
		res["code"] = "-1"
		res["message"] = "Invalid param"
		res["transaction_id"] = transaction_id
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
		fmt.Println("Json error")
		res["code"] = "-1"
		res["message"] = "Json error"
		res["transaction_id"] = transaction_id
		return res
	}

	payloadString := string(jsonData)
	//fmt.Println(payloadString)

	redis_db.ConnectRedis(redisConnection, "", 0)

	//Find partner service in Redis
	partner_service, getRedisErr := redis_db.GetValue("SERVICES:" + partner_id + ":" + adsid)
	fmt.Println("Redis Cache : " + partner_service)
	if getRedisErr != nil {

		// Init database
		postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dbConnection)
		if err != nil {
			panic(err)
		}
		// Test connection
		err = sqlConfig.Ping()
		if err != nil {
			fmt.Println(err)
		}

		var clientService models.ClientService
		queryResult := postgresDB.Where("client_partner_id = ? AND ads_id = ?", partner_id, adsid).First(&clientService)
		if queryResult.Error != nil {
			log.Printf("User not found or error: %v", queryResult.Error)
			res["code"] = "0"
			res["message"] = "Retrived"
			res["transaction_id"] = transaction_id
			return res
		} else {

			// Convert to JSON string
			// jsonData, err := json.Marshal(clientService)
			// if err != nil {
			// 	fmt.Println("Error marshaling JSON:", err)
			// 	res["code"] = "-1"
			// 	res["message"] = "Convert Json error"
			// 	res["transaction_id"] = transaction_id
			// 	return res
			// }

			// Convert the int field to string

			var counter = strconv.Itoa(clientService.PostbackCounter)
			cacheData := "{\"keyword\":\"" + clientService.Keyword + "\",\"shortcode\":\"" + clientService.Shortcode + "\",\"telcoid\":\"" + clientService.TelcoID + "\",\"ads_id\":\"" + clientService.AdsID + "\",\"client_partner_id\":\"" + clientService.ClientPartnerID + "\",\"wap_aoc_refid\":\"" + clientService.WapAocRefID + "\",\"wap_aoc_id\":\"" + clientService.WapAocID + "\",\"wap_aoc_media\":\"" + clientService.WapAocMedia + "\",\"postback_url\":\"" + clientService.PostbackURL + "\",\"dn_url\":\"" + clientService.DNURL + "\",\"postback_counter\":" + counter + "}"
			redis_key := "SERVICE:" + partner_id + ":" + adsid
			ttl := 240 * time.Hour // expires in 240 Hour
			// Set key with TTL
			if err := redis_db.SetWithTTL(redis_key, cacheData, ttl); err != nil {
				//write to file if Redis problem or forward request to AIS
				log.Fatalf("SetWithTTL error: %v , %v", err, payloadString)
			}
		}
	}

	redis_key := "TMVH-MO:" + partner_id + ":" + transaction_id

	ttl := 240 * time.Hour // expires in 240 Hour
	// Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, payloadString, ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v , %v", err, payloadString)
	}

	res["code"] = "302"
	res["partner_id"] = partner_id
	res["refid"] = refid
	res["adsid"] = adsid

	return res
}
