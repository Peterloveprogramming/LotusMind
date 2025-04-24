package chakareport

import "fmt"

const ChakaraReportMakerType = "ChakaraReportMaker"
const DummyChakaraReportMakerType = "DummyChakaraReportMaker"

func CreateChakraMakerForTesting(makerType string) (Maker, error) {
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
