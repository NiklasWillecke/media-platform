package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Queue struct {
	Client valkey.Client
	Name   string
}

func Init(name string, address string) *Queue {

	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{address}})
	if err != nil {
		panic(err)
	}

	return &Queue{
		Client: client,
		Name:   name,
	}

}

func (q *Queue) Produce() {
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		job := fmt.Sprintf("Job-%d", i)
		err := q.Client.Do(
			ctx,
			q.Client.B().Lpush().Key(q.Name).Element(job).Build(),
		).Error()
		if err != nil {
			log.Fatalf("Fehler beim LPUSH: %v", err)
		}
		fmt.Println("Gesendet:", job)
	}
}

func (q *Queue) Consume() {

	ctx := context.Background()
	fmt.Println("Warte auf Jobs ...")
	for {
		resp := q.Client.Do(ctx, q.Client.B().Brpop().Key(q.Name).Timeout(0).Build())
		values, err := resp.AsStrSlice()
		if err != nil {
			log.Fatalf("Fehler beim BLPOP: %v", err)
		}

		// Antwort besteht aus [queueName, value]
		fmt.Printf("Bearbeite: %s\n", values[1])
		time.Sleep(500 * time.Millisecond) // simulierte Arbeit
	}
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
