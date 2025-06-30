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

	log.Println("[DEBUG] Starting GetChakaraReportByUniqueCode")
	log.Printf("[DEBUG] Inputs - Email: %s, UniqueCode: %s\n", email, uniqueCode)

	ctx := context.TODO() // Consider replacing with real context

	prefix := ReportFolderName + "/" + email + "/"
	filePrefix := prefix + uniqueCode + "-"

	log.Printf("[DEBUG] S3 bucket: %s | Prefix: %s | FilePrefix: %s\n", maker.awsBucketName, prefix, filePrefix)

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(maker.awsBucketName),
		Prefix: aws.String(prefix),
	}

	var matchingObjects []types.Object
	paginator := s3.NewListObjectsV2Paginator(maker.s3Client, listInput)

	log.Println("[DEBUG] Beginning S3 pagination...")

	pageCount := 0
	for paginator.HasMorePages() {
		log.Printf("[DEBUG] Fetching page #%d\n", pageCount+1)
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Printf("[ERROR] Failed fetching S3 page: %v\n", err)
			return "", fmt.Errorf("failed to list objects in S3 bucket '%s' with prefix '%s': %w", maker.awsBucketName, prefix, err)
		}
		log.Printf("[DEBUG] Page #%d contains %d objects\n", pageCount+1, len(page.Contents))

		for _, obj := range page.Contents {
			if obj.Key != nil && strings.HasPrefix(*obj.Key, filePrefix) && strings.HasSuffix(*obj.Key, ReportFileExtension) {
				log.Printf("[DEBUG] Matching object found: %s\n", *obj.Key)
				matchingObjects = append(matchingObjects, obj)
			}
		}
		pageCount++
	}

	if len(matchingObjects) == 0 {
		log.Println("[WARN] No matching objects found.")
		return "", fmt.Errorf("report with unique code '%s' not found for email '%s' in S3 bucket '%s'", uniqueCode, email, maker.awsBucketName)
	}

	log.Printf("[DEBUG] %d matching object(s) found\n", len(matchingObjects))

	if len(matchingObjects) > 1 {
		log.Println("[DEBUG] Sorting matching objects by key (descending)...")
		sort.Slice(matchingObjects, func(i, j int) bool {
			keyI := ""
			if matchingObjects[i].Key != nil {
				keyI = *matchingObjects[i].Key
			}
			keyJ := ""
			if matchingObjects[j].Key != nil {
				keyJ = *matchingObjects[j].Key
			}
			return keyI > keyJ
		})
	}

	targetKey := matchingObjects[0].Key
	if targetKey == nil {
		log.Println("[ERROR] Target key is nil after sorting")
		return "", fmt.Errorf("internal error: found object for unique code '%s' has nil key", uniqueCode)
	}

	log.Printf("[DEBUG] Getting object from S3: %s\n", *targetKey)
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(maker.awsBucketName),
		Key:    targetKey,
	}

	resp, err := maker.s3Client.GetObject(ctx, getObjectInput)
	if err != nil {
		log.Printf("[ERROR] Failed to get object: %v\n", err)
		return "", fmt.Errorf("failed to get object '%s' from S3 bucket '%s': %w", *targetKey, maker.awsBucketName, err)
	}
	defer resp.Body.Close()

	log.Println("[DEBUG] Successfully got object, starting to read body...")
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read S3 object body: %v\n", err)
		return "", fmt.Errorf("failed to read content of object '%s' from S3: %w", *targetKey, err)
	}

	log.Printf("[INFO] Successfully retrieved file '%s' from bucket '%s'\n", *targetKey, maker.awsBucketName)
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
