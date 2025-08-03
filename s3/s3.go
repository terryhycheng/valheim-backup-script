package s3

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Helper interface {
	// Getters
	Bucket() string
	Folder() string

	// Actions
	Upload(filePath string) error
	ListOldObjects(daysOld int) ([]string, error)
	Delete(keys []string) error
}

type s3Helper struct {
	client *s3.Client
	bucket string
	folder string
}

func New(bucket string, folder string) S3Helper {
	baseEndpoint := os.Getenv("S3_BASE_ENDPOINT")

	return &s3Helper{
		client: s3.New(s3.Options{
			BaseEndpoint: aws.String(baseEndpoint),
			Region:       "us-east-1",
			Credentials: aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					"",
				),
			),
		}),
		bucket: bucket,
		folder: folder,
	}
}

func (s *s3Helper) Bucket() string { return s.bucket }
func (s *s3Helper) Folder() string { return s.folder }

// Upload uploads a file to S3
// Example:
//
//	s3 := s3Helper.New("my-bucket", "my-folder")
//	s3.Upload("/tmp/valheim-backup-2025-08-03.tar.gz")
func (s *s3Helper) Upload(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	key := path.Join(s.folder, filepath.Base(filePath)) // keep filename in key

	s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          f,
		ContentLength: aws.Int64(stat.Size()),
	})
	return nil
}

// ListOldObjects returns keys of objects older than `daysOld`
func (s *s3Helper) ListOldObjects(daysOld int) ([]string, error) {
	ctx := context.Background()
	cutoff := time.Now().AddDate(0, 0, -daysOld)

	var keys []string
	var continuationToken *string

	for {
		resp, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.bucket),
			Prefix:            aws.String(path.Clean(s.folder) + "/"),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}

		for _, obj := range resp.Contents {
			if obj.LastModified.Before(cutoff) {
				keys = append(keys, *obj.Key)
			}
		}

		if *resp.IsTruncated && resp.NextContinuationToken != nil {
			continuationToken = resp.NextContinuationToken
		} else {
			break
		}
	}
	return keys, nil
}

func (s *s3Helper) Delete(keys []string) error {
	ctx := context.Background()

	for _, key := range keys {
		fmt.Printf("Deleting %s...\n", key)
		_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
