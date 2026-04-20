package configs

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func initWhatEverBucket(ctx context.Context, cfg *EnvData) (*s3.Client, error) {

	staticProvider := credentials.NewStaticCredentialsProvider(cfg.STORAGE_SERVICE_ACCESS_KEY_ID, cfg.STORAGE_SERVICE_SECRET_ACCESS_KEY, "")

	cfgs, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithCredentialsProvider(staticProvider))

	client := s3.NewFromConfig(cfgs, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.STORAGE_SERVICE_ENDPOINT)
		o.UsePathStyle = true
	})
	return client, err
}

func ensureBucketExists(ctx context.Context, client *s3.Client, bucket string) {
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var ownedErr *types.BucketAlreadyOwnedByYou
		var existsErr *types.BucketAlreadyExists

		if errors.As(err, &ownedErr) {
			log.Printf("Bucket %s already exists,skipping creation.", bucket)
			return
		}
		if errors.As(err, &existsErr) {
			log.Fatalf("Bucket %s already exists but owned by a different account", bucket)
		}
		log.Fatalf("Failed to create bucket: %v", err)
	}
	log.Printf("Bucket %s Created Successfully", bucket)
}

func PerformAllNecessaryActivationStep(ctx context.Context, cfg *EnvData) (*s3.Client, *s3.PresignClient) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := initWhatEverBucket(ctx, cfg)
	if err != nil {
		log.Fatal("Couldnt load default configurations")
	}
	ensureBucketExists(ctx, client, cfg.BUCKET_NAME)
	presignedClient := s3.NewPresignClient(client)
	return client, presignedClient
}
