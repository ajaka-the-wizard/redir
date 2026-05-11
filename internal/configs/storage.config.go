package configs

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func initWhatEverBucket(ctx context.Context, cfg *EnvData) (*s3.Client, error) {

	staticProvider := credentials.NewStaticCredentialsProvider(cfg.STORAGE_SERVICE_ACCESS_KEY_ID, cfg.STORAGE_SERVICE_SECRET_ACCESS_KEY, "")

	cfgs, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithCredentialsProvider(staticProvider))

	client := s3.NewFromConfig(cfgs, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.STORAGE_SERVICE_ENDPOINT)
		o.UsePathStyle = true
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})
	return client, err
}

func ensureBucketExists(ctx context.Context, client *s3.Client, bucket string, logger *slog.Logger) {
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var ownedErr *types.BucketAlreadyOwnedByYou
		var existsErr *types.BucketAlreadyExists

		if errors.As(err, &ownedErr) {
			logger.Info("bucket already exists, skipping creation", "bucket", bucket)
			return
		}
		if errors.As(err, &existsErr) {
			logger.Error("bucket already exists but owned by different account", "bucket", bucket)
			return
		}
		logger.Error("failed to create bucket", "bucket", bucket, "error", err.Error())
		return
	}
	logger.Info("bucket created successfully", "bucket", bucket)
}

func PerformAllNecessaryActivationStep(ctx context.Context, cfg *EnvData, logger *slog.Logger) (*s3.Client, *s3.PresignClient, *transfermanager.Client) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := initWhatEverBucket(ctx, cfg)
	if err != nil {
		logger.Error("could not load default storage configuration", "error", err.Error())
		os.Exit(1)
	}
	ensureBucketExists(ctx, client, cfg.BUCKET_NAME, logger)
	presignedClient := s3.NewPresignClient(client)
	tm := transfermanager.New(client)
	return client, presignedClient, tm
}
