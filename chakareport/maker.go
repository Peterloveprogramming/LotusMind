package chakareport

import (
	"fmt"
	"strings"
)

type ChakraInfo struct {
	ChakraName   string `json:"chakra_name" binding:"required"`
	ChakraScore  int32  `json:"chakra_score" binding:"required"`
	ChakraStatus string `json:"chakra_status" binding:"required"`
}

// Maker is an interface for managing tokens
type Maker interface {
	GenerateChakaraReport(chakaraInfo []ChakraInfo, language string) ([]byte, error)
}

func ChakraMaker(appEnvironment string, apiUrl string) (Maker, error) {
	envLower := strings.ToLower(appEnvironment)

	if envLower == "dev" || envLower == "development" {
		// Dummy maker doesn't need the URL, so no change here
		maker, err := NewDummyChakaraReportMaker()
		if err != nil {
			return nil, fmt.Errorf("failed to create dummy chakra report maker: %w", err)
		}
		return maker, nil
	}

	// Pass the apiUrl to the real maker's constructor
	maker, err := NewChakaraReportMaker(apiUrl) // Pass the URL here
	if err != nil {
		// Error message already includes context from NewChakaraReportMaker
		return nil, fmt.Errorf("failed to create chakra report maker: %w", err)
	}
	return maker, nil
}
