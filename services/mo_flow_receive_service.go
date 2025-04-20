package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"CyberusGolangShareLibrary/redis_db"
	"log"
	"strconv"
	"time"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

// Struct to map the expected JSON fields
type MoFlowReceiveRequest struct {
	IdPartner    string `json:"id_partner"`
	RefIdPartner string `json:"refid_partner"`
	MediaPartner string `json:"media_partner"`
	NamePartner  string `json:"name_partner"`
}

func MoFlowReceiveProcessRequest(r *http.Request) map[string]string {

	// Get current time
	now := time.Now()
	// Unix timestamp in nanoseconds
	timestamp := (now.UnixNano())
	nano_timestamp := strconv.FormatInt(timestamp, 10)

	// Generate a random UUID (UUID v4)
	transaction_id := uuid.New().String()

	// Get a Client IP address
	ip := r.RemoteAddr

	fmt.Println("ClientIP : " + ip)

	dns := "host=localhost user=root password=11111111 dbname=cyberus_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@cyberus"
	// Init database
	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dns)
	if err != nil {
		panic(err)
	}
	// Test connection
	err = sqlConfig.Ping()
	if err != nil {
		fmt.Println(err)
	}
	postgresDB.DB()

	redis_db.ConnectRedis()

	redis_key := transaction_id
	redis_value := transaction_id
	ttl := 1 * time.Hour // expires in 1 Hour

	// Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, redis_value, ttl); err != nil {
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
	res := map[string]string{
		"code":           "0",
		"message":        "retrieved",
		"timestamp":      nano_timestamp,
		"transaction_id": transaction_id,
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//http.Error(w, "Failed to read body", http.StatusBadRequest)
		return res
	}
	defer r.Body.Close()

	return res
}
