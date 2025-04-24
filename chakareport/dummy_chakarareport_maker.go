package chakareport

import (
	"encoding/json"
	"fmt"
)

type DummyChakaraReportMake struct {
}

// this creates a new pasetoMaker interface
func NewDummyChakaraReportMaker() (Maker, error) {
	maker := &DummyChakaraReportMake{}
	return maker, nil
}

func (maker *DummyChakaraReportMake) GenerateChakaraReport(chakaraInfo []ChakraInfo, language string) ([]byte, error) {
	// Create a map or an anonymous struct to hold the desired output structure
	reportData := map[string]interface{}{
		"chakarinfo": chakaraInfo,
		"language":   language,
	}

	// Serialize the combined data structure into a JSON byte slice
	reportBytes, err := json.Marshal(reportData)
	if err != nil {
		// Handle the error appropriately
		return nil, fmt.Errorf("failed to marshal chakra report data to JSON: %w", err)
	}

	// reportBytes now contains the JSON representation of the combined data
	return reportBytes, nil
}

func (maker *DummyChakaraReportMake) GetType() string {
	return DummyChakaraReportMakerType
}
