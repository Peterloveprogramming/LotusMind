package storage

import "fmt"

const TimeStampFormat = "200601021504"
const ReportFolderName = "chakara-report"
const ReportFileExtension = ".txt"

// 0755 provides standard permissions (rwxr-xr-x).
const FilePermission = 0755

const LocalStorageMakerType = "LocalStorageMaker"
const S3StorageMakerType = "S3StorageMaker"

func CreateStorageMakerForTesting(makerType string) (Maker, error) {
	appUrl := "http://api:8080"
	prod := "prod"
	dev := "dev"

	switch makerType {
	case ChakaraReportMakerType:

		ChakararReportMaker, err := ChakraMaker(prod, appUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to create ChakaraReportMaker for testing: %w", err)
		}
		return ChakararReportMaker, nil
	case DummyChakaraReportMakerType:
		DummyChakararReportMaker, err := ChakraMaker(dev, appUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to create DummyChakararReportMaker for testing: %w", err)
		}
		return DummyChakararReportMaker, nil
	default:
		// Return an error if the provided type string is not recognized
		return nil, fmt.Errorf("invalid maker type provided for testing: %s", makerType)
	}
}
