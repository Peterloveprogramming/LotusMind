package storage

import "fmt"

const TimeStampFormat = "200601021504"
const ReportFolderName = "chakara-report"
const ReportAnswerFolderName = "chakara-report-answers"

const ReportFileExtension = ".txt"

// 0755 provides standard permissions (rwxr-xr-x).
const FilePermission = 0755

const LocalStorageMakerType = "LocalStorageMaker"
const S3StorageMakerType = "S3StorageMaker"

func CreateStorageMakerForTesting(makerType string) (Maker, error) {
	prod := "prod"
	dev := "dev"
	awsRegion := "test"
	awsAccessKeydId := "test"
	awsAccessKey := "test"
	awsBucketName := "test"

	switch makerType {
	case LocalStorageMakerType:

		Localtorage, err := StorageMaker(dev, awsRegion, awsAccessKeydId, awsAccessKey, awsBucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to create Localtorage for testing: %w", err)
		}
		return Localtorage, nil
	// do not use S3StorageMakerType. unless providing valid aws key and bucket name
	case S3StorageMakerType:
		S3Storage, err := StorageMaker(prod, awsRegion, awsAccessKeydId, awsAccessKey, awsBucketName)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3Storage for testing: %w", err)
		}
		return S3Storage, nil
	default:
		// Return an error if the provided type string is not recognized
		return nil, fmt.Errorf("invalid maker type provided for testing: %s", makerType)
	}
}
