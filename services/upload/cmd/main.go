package main

import (
	"fmt"
	"log"
	"os"
	dataStore "streaming-platform/services/upload/pkg/s3"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	region := os.Getenv("REGION")
	accessKeyID := os.Getenv("ACCESKEYID")
	secretAccessKey := os.Getenv("SECRETACCESKEY")
	endpoint := os.Getenv("ENDPOINT")

	fmt.Println(region)
	fmt.Println(accessKeyID)
	fmt.Println(secretAccessKey)
	fmt.Println(endpoint)

	if accessKeyID == "" || secretAccessKey == "" || region == "" ||
		endpoint == "" {
		log.Fatal(
			"missing env: RUSTFS_ACCESS_KEY_ID / " +
				"RUSTFS_SECRET_ACCESS_KEY / RUSTFS_REGION / " +
				"RUSTFS_ENDPOINT_URL",
		)
	}

	client := dataStore.Init(accessKeyID, secretAccessKey, region, endpoint)

	client.CreateBucket("peter")

}

/*
	ctx := context.Background()

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	queue := "jobs"

	// PRODUCER – schreibt 3 Aufträge
	for i := 1; i <= 5; i++ {
		job := fmt.Sprintf("Job-%d", i)
		err := client.Do(
			ctx,
			client.B().Lpush().Key(queue).Element(job).Build(),
		).Error()
		if err != nil {
			log.Fatalf("Fehler beim LPUSH: %v", err)
		}
		fmt.Println("Gesendet:", job)
	}

*/
