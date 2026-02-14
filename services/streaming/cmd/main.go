package main

import (
	"bytes"
	"fmt"

	"github.com/NiklasWillecke/media-platform/services/streaming/internal/cache"
)

func main() {

	fmt.Println("Lets Go")
	/*
		s3, err := storage.NewS3Client("192.168.178.115:9000", "minioadmin", "minioadmin", false)
		if err != nil {
			log.Fatal(err)
		}

		s3.GetObjekt("testbucket", "generated-image.png")

		handler.StartFileServer()*/

	cache := cache.NewLRUCache(1024 * 1024 * 10)
	cache.Set("1", []byte{'W', 'o', 'r', 'l', 'd'})

	value, ok := cache.Get("1")
	fmt.Println(string(value))

	if !ok || bytes.Equal(value, []byte{'W', 'o', 'r', 'l', 'd'}) {
		fmt.Println("Test")
	}

}
