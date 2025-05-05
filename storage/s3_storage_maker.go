package storage

import (
	"context"
	"encoding/json"
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

// GetChakaraReportByUniqueCode retrieves the content of the chakra report
// matching the uniqueCode for a given email from the S3 bucket.
// If multiple matches exist, it returns the latest one based on the timestamp.
func (maker *S3StorageMaker) GetChakaraReportByUniqueCode(email string, uniqueCode string) (string, error) {
	if uniqueCode == "" {
		return "", fmt.Errorf("uniqueCode cannot be empty")
	}
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	ctx := context.TODO() // Use context.TODO() for now

	// Define the prefix for listing objects in the user's "folder"
	prefix := ReportFolderName + "/" + email + "/"
	// Define the specific prefix for the unique code within the user's folder
	filePrefix := prefix + uniqueCode + "-"

	// List objects in the bucket with the user's folder prefix
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(maker.awsBucketName),
		Prefix: aws.String(prefix), // List the whole user folder first
	}

	var matchingObjects []types.Object // Store objects matching the uniqueCode

	// Paginate through results
	paginator := s3.NewListObjectsV2Paginator(maker.s3Client, listInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to list objects in S3 bucket '%s' with prefix '%s': %w", maker.awsBucketName, prefix, err)
		}

		// Filter for files matching the specific uniqueCode prefix and extension
		for _, obj := range page.Contents {
			if obj.Key != nil && strings.HasPrefix(*obj.Key, filePrefix) && strings.HasSuffix(*obj.Key, ReportFileExtension) {
				matchingObjects = append(matchingObjects, obj)
			}
		}
	}

	// Check if any matching objects were found
	if len(matchingObjects) == 0 {
		return "", fmt.Errorf("report with unique code '%s' not found for email '%s' in S3 bucket '%s'", uniqueCode, email, maker.awsBucketName)
	}

	// If multiple matches, sort descending by key (filename) to get the latest timestamp
	if len(matchingObjects) > 1 {
		sort.Slice(matchingObjects, func(i, j int) bool {
			// Dereference pointers safely
			keyI := ""
			if matchingObjects[i].Key != nil {
				keyI = *matchingObjects[i].Key
			}
			keyJ := ""
			if matchingObjects[j].Key != nil {
				keyJ = *matchingObjects[j].Key
			}
			// Sort descending
			return keyI > keyJ
		})
	}

	// Get the key of the target object (the first one after sorting, which is the latest)
	targetKey := matchingObjects[0].Key
	if targetKey == nil {
		return "", fmt.Errorf("internal error: found object for unique code '%s' has nil key", uniqueCode)
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

// SaveChakaraReportAnswersAsText saves the chakra report answers map as a JSON file
// in the specified S3 bucket under the "chakara-report-answers" folder.
// The filename includes the uniqueId and a timestamp (YYYYMMDDHHMM).
func (maker *S3StorageMaker) SaveChakaraReportAnswersAsText(email string, uniqueId string, answers map[string]string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if uniqueId == "" {
		return fmt.Errorf("uniqueId cannot be empty")
	}
	if len(answers) == 0 {
		log.Printf("Warning: Attempting to save empty answers map for email '%s', uniqueId '%s' to S3. Skipping.\n", email, uniqueId)
		return nil // Or return an error if empty answers are invalid
	}

	ctx := context.TODO() // Use context.TODO() for now

	// Get the current time and format it
	currentTime := time.Now()
	timestamp := currentTime.Format(TimeStampFormat) // Go's reference time format

	// Define the folder structure and file name within the bucket for answers
	userAnswersFolder := ReportAnswerFolderName + "/" + email + "/"
	// Construct the answers filename with uniqueId and timestamp
	answersFilename := fmt.Sprintf("%s-%s%s", uniqueId, timestamp, ReportFileExtension)
	objectKey := userAnswersFolder + answersFilename

	// Marshal the answers map into JSON format
	jsonData, err := json.MarshalIndent(answers, "", "  ") // Use MarshalIndent for readability
	if err != nil {
		return fmt.Errorf("failed to marshal answers map to JSON for S3 object '%s': %w", objectKey, err)
	}

	// Upload the JSON data to S3
	_, err = maker.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(maker.awsBucketName),
		Key:         aws.String(objectKey),
		Body:        strings.NewReader(string(jsonData)), // Upload the JSON string
		ContentType: aws.String("application/json"),      // Set content type to JSON
	})

	if err != nil {
		log.Printf("ERROR: Couldn't upload answers file '%s' to bucket '%s'. Reason: %v\n", objectKey, maker.awsBucketName, err)
		return fmt.Errorf("failed to upload answers to S3 bucket '%s': %w", maker.awsBucketName, err)
	}

	log.Printf("Successfully uploaded answers file '%s' to bucket '%s'\n", objectKey, maker.awsBucketName)
	return nil
}
