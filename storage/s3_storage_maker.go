package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3StorageMaker struct {
	awsRegion          string
	awsAccessKeyId     string
	awsSecretAccessKey string // Renamed for clarity and consistency with AWS SDK
	awsBucketName      string
	s3Client           *s3.Client // Store the S3 client
}

// NewS3StorageMaker creates a new S3 storage maker instance.
// It requires AWS region, access key ID, secret access key, and bucket name.
func NewS3StorageMaker(awsRegion, awsAccessKeyId, awsSecretAccessKey, awsBucketName string) (Maker, error) {
	// Check if any required parameter is empty
	if awsRegion == "" || awsAccessKeyId == "" || awsSecretAccessKey == "" || awsBucketName == "" {
		return nil, fmt.Errorf("missing required AWS configuration: region, access key ID, secret access key, or bucket name cannot be empty")
	}
	println("aws region", awsRegion)
	println("aws access key id", awsAccessKeyId)
	println("aws secret access key", awsSecretAccessKey)
	println("aws bucket name", awsBucketName)

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion), // Set the region
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKeyId, awsSecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create the maker instance
	maker := &S3StorageMaker{
		awsRegion:          awsRegion,
		awsAccessKeyId:     awsAccessKeyId,
		awsSecretAccessKey: awsSecretAccessKey,
		awsBucketName:      awsBucketName,
		s3Client:           s3Client, // Store the initialized client
	}
	return maker, nil
}

// SaveChakaraReportAsText saves the chakra report content as a text file in the specified S3 bucket.
// The filename includes the uniqueId and a timestamp (YYYYMMDDHHMM).
func (maker *S3StorageMaker) SaveChakaraReportAsText(email string, uniqueId string, content string) error {
	ctx := context.TODO() // Use context.TODO() for now, consider passing a context if needed

	// Get the current time and format it
	currentTime := time.Now()
	// Format: YYYYMMDDHHMM (e.g., 202310271504)
	timestamp := currentTime.Format(TimeStampFormat) // Go's reference time format

	// Define the folder structure and file name within the bucket
	userFolder := "chakara-report/" + email + "/"
	// Construct the report name with uniqueId and timestamp
	reportName := fmt.Sprintf("%s-%s.txt", uniqueId, timestamp)
	objectKey := userFolder + reportName

	// Upload the file content to S3
	_, err := maker.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(maker.awsBucketName), // Use bucket name from the maker struct
		Key:    aws.String(objectKey),
		Body:   strings.NewReader(content), // Use the content passed to the function
		// ContentType: aws.String("text/plain"), // Optional: Set content type if desired
	})

	if err != nil {
		log.Printf("ERROR: Couldn't upload file '%s' to bucket '%s'. Reason: %v\n", objectKey, maker.awsBucketName, err)
		return fmt.Errorf("failed to upload report to S3 bucket '%s': %w", maker.awsBucketName, err)
	}

	log.Printf("Successfully uploaded file '%s' to bucket '%s'\n", objectKey, maker.awsBucketName)
	return nil // Return nil on success
}

// GetChakaraReportByTestNum retrieves the content of the Nth oldest chakra report
// for a given email from the S3 bucket.
// testNum = 1 retrieves the oldest, testNum = 2 retrieves the second oldest, etc.
func (maker *S3StorageMaker) GetChakaraReportByTestNum(email string, testNum int) (string, error) {
	if testNum <= 0 {
		return "", fmt.Errorf("testNum must be greater than 0")
	}

	ctx := context.TODO() // Use context.TODO() for now

	// Define the prefix for listing objects in the user's "folder"
	prefix := ReportFolderName + "/" + email + "/"

	// List objects in the bucket with the specified prefix
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(maker.awsBucketName),
		Prefix: aws.String(prefix),
	}

	var reportObjects []types.Object // Use types.Object

	// Paginate through results if necessary (though unlikely for single user reports)
	paginator := s3.NewListObjectsV2Paginator(maker.s3Client, listInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to list objects in S3 bucket '%s' with prefix '%s': %w", maker.awsBucketName, prefix, err)
		}
		// Filter for actual files (not the prefix "folder" itself) and correct extension
		for _, obj := range page.Contents {
			// Ensure obj.Key is not nil and doesn't represent the directory itself
			if obj.Key != nil && *obj.Key != prefix && strings.HasSuffix(*obj.Key, ReportFileExtension) {
				reportObjects = append(reportObjects, obj)
			}
		}
	}

	// Sort the objects by key (filename). Since the timestamp is YYYYMMDDHHMM,
	// this effectively sorts them chronologically (oldest first).
	sort.Slice(reportObjects, func(i, j int) bool {
		// Dereference pointers safely
		keyI := ""
		if reportObjects[i].Key != nil {
			keyI = *reportObjects[i].Key
		}
		keyJ := ""
		if reportObjects[j].Key != nil {
			keyJ = *reportObjects[j].Key
		}
		return keyI < keyJ
	})

	// Calculate the index (0-based)
	index := testNum - 1

	// Check if the requested testNum is valid
	if index < 0 || index >= len(reportObjects) {
		return "", fmt.Errorf("report number %d not found for email %s in S3 bucket '%s' (only %d reports exist)", testNum, email, maker.awsBucketName, len(reportObjects))
	}

	// Get the key for the requested report number
	targetKey := reportObjects[index].Key
	if targetKey == nil {
		// This should ideally not happen if filtering worked, but good to check
		return "", fmt.Errorf("internal error: found object at index %d has nil key", index)
	}

	// Get the object content from S3
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(maker.awsBucketName),
		Key:    targetKey,
	}

	resp, err := maker.s3Client.GetObject(ctx, getObjectInput)
	if err != nil {
		return "", fmt.Errorf("failed to get object '%s' from S3 bucket '%s': %w", *targetKey, maker.awsBucketName, err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the content
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read content of object '%s' from S3: %w", *targetKey, err)
	}

	log.Printf("Successfully retrieved file '%s' from bucket '%s'\n", *targetKey, maker.awsBucketName)
	return string(bodyBytes), nil
}
