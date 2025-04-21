package storage

import (
	"fmt"
	"strings"
)

type Maker interface {
	SaveChakaraReportAsText(email string, uniqueId string, content string) error
	GetChakaraReportByUniqueCode(email string, uniqueCode string) (string, error)
}

func StorageMaker(appEnvironment string, awsRegion, awsAccessKeyId, awsSecretAccessKey, awsBucketName string) (Maker, error) {
	envLower := strings.ToLower(appEnvironment)

	// Use Local Storage for "local" or "dev" environments
	if envLower == "local" || envLower == "dev" || envLower == "development" {
		fmt.Printf("Using Local Storage Maker for environment: %s\n", envLower)
		// Call NewLocalStorageMaker without rootDir
		maker, err := NewLocalStorageMaker()
		if err != nil {
			// Error from NewLocalStorageMaker is already descriptive
			return nil, fmt.Errorf("failed to create local storage maker: %w", err)
		}
		return maker, nil
	}
	maker, err := NewS3StorageMaker(awsRegion, awsAccessKeyId, awsSecretAccessKey, awsBucketName) // Pass the URL here
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 storage maker maker: %w", err)
	}
	return maker, nil
}
