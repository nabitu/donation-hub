package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func InitContainer() *model.Container {
	cfg := InitConfig()
	s3Client, err := InitS3Client(cfg)
	if err != nil {
		panic(err)
	}

	sqlClient, err := InitDatabase(cfg)
	if err != nil {
		panic(err)
	}

	return &model.Container{
		Connection: &model.Connection{
			S3Client: s3Client,
			DB:       sqlClient,
		},
		Config: cfg,
	}
}

func InitConfig() *model.Config {
	appPort, ok := os.LookupEnv("APP_PORT")
	if !ok {
		panic("APP_PORT not provided")
	}

	dbDriverName, ok := os.LookupEnv("DATABASE_DRIVER_NAME")
	if !ok {
		panic("DATABASE_DRIVER_NAME not provided")
	}

	dbDataSource, ok := os.LookupEnv("DATABASE_DATA_SOURCE")
	if !ok {
		panic("DATABASE_DATA_SOURCE not provided")
	}

	awsRegion, ok := os.LookupEnv("AWS_DEFAULT_REGION")
	if !ok {
		panic("AWS_DEFAULT_REGION not provided")
	}

	awsAccessKeyId, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		panic("AWS_ACCESS_KEY_ID not provided")
	}

	awsSecretAccessKeyId, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		panic("AWS_SECRET_ACCESS_KEY not provided")
	}

	awsEndpoint, ok := os.LookupEnv("LOCALSTACK_ENDPOINT")
	if !ok {
		panic("LOCALSTACK_ENDPOINT not provided")
	}

	awsUseStylePath, ok := os.LookupEnv("AWS_USE_PATH_STYLE_ENDPOINT")
	if !ok {
		panic("AWS_USE_PATH_STYLE_ENDPOINT not provided")
	}
	awsBucketName, ok := os.LookupEnv("AWS_BUCKET")
	if !ok {
		panic("AWS_BUCKET not provided")
	}

	tokenSecretKey := "supersecrethehehe"
	tokenIssuer := "Donation Hub"

	return &model.Config{
		AppPort:                 appPort,
		DBDriverName:            dbDriverName,
		DBDataSource:            dbDataSource,
		AwsDefaultRegion:        awsRegion,
		AwsAccessKey:            awsAccessKeyId,
		AwsSecretAccessKey:      awsSecretAccessKeyId,
		AwsEndpoint:             awsEndpoint,
		AwsUsePathStyleEndpoint: awsUseStylePath == "1",
		AwsS3Bucket:             awsBucketName,
		TokenSecretKey:          tokenSecretKey,
		TokenIssuer:             tokenIssuer,
	}
}

func InitDatabase(cfg *model.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.DBDriverName, cfg.DBDataSource)
	log.Println("Get Database Connection")
	log.Println(cfg.DBDriverName)
	log.Println(cfg.DBDataSource)
	if err != nil {
		log.Fatalln(err)

		return nil, err
	}

	fmt.Println("Database Connected")

	return db, nil
}

func InitS3Client(cfg *model.Config) (s3Client *s3.Client, err error) {
	s3Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AwsDefaultRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AwsAccessKey,
			cfg.AwsSecretAccessKey,
			os.Getenv("AWS_SESSION_TOKEN"), // optional tapi wajib di definisikan
		)),
	)

	if err != nil {
		return s3Client, fmt.Errorf("failed to load configuration, %w", err)
	}

	s3Client = s3.NewFromConfig(s3Config, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(cfg.AwsEndpoint)
		options.UsePathStyle = cfg.AwsUsePathStyleEndpoint
	})

	return s3Client, err
}
