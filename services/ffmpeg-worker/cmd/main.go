package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

func main() {

	ctx := context.Background()

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	queue := "jobs"

	// CONSUMER – liest sie mit BLPOP (blockierend)
	fmt.Println("Warte auf Jobs ...")
	for {
		resp := client.Do(ctx, client.B().Brpop().Key(queue).Timeout(0).Build())
		values, err := resp.AsStrSlice()
		if err != nil {
			log.Fatalf("Fehler beim BLPOP: %v", err)
		}

		// Antwort besteht aus [queueName, value]
		fmt.Printf("Bearbeite: %s\n", values[1])
		time.Sleep(500 * time.Millisecond) // simulierte Arbeit
	}

}
