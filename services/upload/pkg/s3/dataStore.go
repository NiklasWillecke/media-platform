package dataStore

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client    *s3.Client
	Presigner *s3.PresignClient
}

func Init(accessKeyID string, secretAccessKey string, region string, endpoint string) *S3Service {

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKeyID,
				secretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &S3Service{
		Client:    client,
		Presigner: s3.NewPresignClient(client),
	}
}

// CreateBucket creates a new S3 bucket with the given name using the S3 client.
func (c *S3Service) CreateBucket(name string) {
	_, err := c.Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		log.Fatalf("create bucket failed: %v", err)
	}
}

func (c *S3Service) CreatePresignedUrl(bucket string, key string) string {
	presignResult, err := c.Presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	s3.WithPresignExpires(15 * time.Minute)

	if err != nil {
		panic("Couldn't get presigned URL for PutObject")
	}

	return presignResult.URL
}

/*
func (c *S3Service) IsVideoUploaded() (bool, error) {
	out, err := c.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NotFound" {
				return &UploadCheckResult{Exists: false}, nil
			}
		}
		return nil, fmt.Errorf("head object fehlgeschlagen: %w", err)
	}
}

*/
