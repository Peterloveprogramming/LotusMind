package chakareport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateNewDummyChakaraReportMaker(t *testing.T) {
	maker, err := CreateChakraMakerForTesting(DummyChakaraReportMakerType)
	require.NoError(t, err)
	require.NotNil(t, maker)
	require.Equal(t, maker.GetType(), DummyChakaraReportMakerType)
}

func TestGenerateReportWithDummyChakaraReportMaker(t *testing.T) {
	maker, err := CreateChakraMakerForTesting(DummyChakaraReportMakerType)
	require.NoError(t, err)
	require.NotNil(t, maker)
	require.Equal(t, maker.GetType(), DummyChakaraReportMakerType)

	testChakaraInfo := []ChakraInfo{
		{ChakraName: "Root", ChakraScore: 75, ChakraStatus: "Balanced"},
		{ChakraName: "Sacral", ChakraScore: 40, ChakraStatus: "Blocked"},
	}
	testLanguage := "en"

	reportBytes, err := maker.GenerateChakaraReport(testChakaraInfo, testLanguage)
	require.NoError(t, err)
	require.NotNil(t, reportBytes)
}
