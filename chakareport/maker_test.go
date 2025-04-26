package chakareport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChakraMakerGenerateChakaraReportMaker(t *testing.T) {
	maker, err := CreateChakraMakerForTesting(ChakaraReportMakerType)
	require.NoError(t, err)
	require.NotNil(t, maker)
	require.Equal(t, maker.GetType(), ChakaraReportMakerType)
}

func TestChakraMakerGenerateDummyChakaraReportMaker(t *testing.T) {
	maker, err := CreateChakraMakerForTesting(DummyChakaraReportMakerType)
	require.NoError(t, err)
	require.NotNil(t, maker)
	require.Equal(t, maker.GetType(), DummyChakaraReportMakerType)
}
