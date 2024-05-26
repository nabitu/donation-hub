package projectfilestorage

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/isdzulqor/donation-hub/internal/core/model"
	"log"
	"strings"
	"time"
)

type Storage struct {
	container *model.Container
}

func (s Storage) RequestUploadUrl(mimeType string, fileSize int64) (*model.RequestUploadUrlStorage, error) {
	presignClient := s3.NewPresignClient(s.container.Connection.S3Client)
	bucketName := s.container.Config.AwsS3Bucket
	objectName := fmt.Sprintf("%d_%x.jpg", time.Now().Unix(), makeRandomBytes(8))
	duration := 15 * time.Minute
	expiredAt := time.Now().Add(duration).Unix()
	presignedUrl, err := presignClient.PresignPutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String(objectName),
			ACL:           types.ObjectCannedACLPublicRead,
			ContentType:   aws.String(mimeType),
			ContentLength: aws.Int64(fileSize),
		},
		s3.WithPresignExpires(duration),
	)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create presigned URL, %v", err))
	}

	url := strings.Replace(presignedUrl.URL, "localstack", "localhost", -1)
	return &model.RequestUploadUrlStorage{
		Url:       url,
		ExpiresAt: expiredAt,
	}, nil
}

func New(container *model.Container) *Storage {
	return &Storage{container: container}
}

func makeRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
