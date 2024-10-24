package fsservice

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3FileSystem struct {
	client *s3.Client
	bucket string
}

// Initialize the S3FileSystem with a new S3 client
func NewS3FileSystem(bucket string) (*S3FileSystem, error) {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	// Create a new S3 client
	s3Client := s3.NewFromConfig(cfg)

	return &S3FileSystem{
		client: s3Client,
		bucket: bucket,
	}, nil
}

// ListFiles lists objects (files) in the specified S3 path (prefix)
func (s3fs *S3FileSystem) ListFiles(path string) ([]string, error) {
	// Call S3 to list current objects
	output, err := s3fs.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s3fs.bucket),
		Prefix: aws.String(path), // Use the path as the prefix
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list files in bucket %q, %v", s3fs.bucket, err)
	}

	filelist := []string{}
	for _, item := range output.Contents {
		fmt.Println("Name: ", *item.Key)
		filelist = append(filelist, *item.Key)
	}

	return filelist, nil
}

// GetFileContents downloads the contents of the specified S3 object (file)
func (s3fs *S3FileSystem) GetFileContents(path string) ([]byte, error) {
	// Get the object from the S3 bucket
	output, err := s3fs.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3fs.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file %q from bucket %q, %v", path, s3fs.bucket, err)
	}
	defer output.Body.Close()

	// Read the object data
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file contents: %v", err)
	}

	return buf.Bytes(), nil
}
