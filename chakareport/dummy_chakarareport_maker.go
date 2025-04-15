package chakareport

type DummyChakaraReportMake struct {
}

// this creates a new pasetoMaker interface
func NewDummyChakaraReportMaker() (Maker, error) {
	maker := &DummyChakaraReportMake{}
	return maker, nil
}

func (maker *DummyChakaraReportMake) GenerateChakaraReport(language string) ([]byte, error) {
	println(language)
	reportBytes := []byte("dummy report")
	return reportBytes, nil
}
