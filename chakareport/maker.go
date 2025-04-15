package chakareport

// Maker is an interface for managing tokens
type Maker interface {
	// //CreateToken creates a new token for a specific email and duration
	// CreateToken(email string, duration time.Duration) (string, error)
	GenerateChakaraReport(language string) ([]byte, error)
	// //VerifyToken checks if the token is valid or not
	// VerifyToken(token string) (*Payload, error)
}

// generate chakara report

// i want to be able to just do server.reportMaker.generateChakaraReport();
