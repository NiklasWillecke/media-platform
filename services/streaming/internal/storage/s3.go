package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Client *minio.Client
}

func NewS3Client(endpoint, accessKey, secretKey string, useSSL bool) (*S3Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &S3Client{Client: client}, nil
}

func (m *S3Client) CreateBucket(bucketName string) error {
	// Make a new bucket called testbucket.
	ctx := context.Background()
	err := m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := m.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
		return err
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return nil
}

func (m *S3Client) GetObjekt(bucketName, objektName string) {
	// Make a new bucket called testbucket.
	ctx := context.Background()
	object, err := m.Client.GetObject(ctx, bucketName, objektName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Fehler beim Get Objekt S3 request:", err)
		return
	}
	defer object.Close()

	// Prüfen ob Objekt existiert (triggert ersten Read)
	stat, err := object.Stat()
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			fmt.Printf("Objekt '%s' existiert nicht im Bucket '%s'\n", objektName, bucketName)
			return
		}
		fmt.Println("Fehler beim Abrufen der Objekt-Info:", err)
		return
	}

	fmt.Printf("Lade Objekt herunter (%d Bytes)...\n", stat.Size)

	if err := os.MkdirAll("tmp/"+bucketName, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Datei anlegen
	localFile, err := os.Create("tmp/" + bucketName + "/" + objektName)
	if err != nil {
		fmt.Println("Fehler beim Erstellen der Datei:", err)
		return
	}
	defer localFile.Close()

	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("HLS-Datei erfolgreich gespeichert:", "tmp/"+bucketName+"/"+objektName)
}
