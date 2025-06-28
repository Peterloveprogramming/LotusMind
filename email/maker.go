package sendemail

import (
	"fmt"
)

type Maker interface {
	SendChakaraResult(to []string, uniqueCode string, language string) error
	SendAppReviewSurvey(to []string, language string) error
}

func EmailMaker(frontendUrl string, sendEmail string, emailPassword string, emailSmtp string, emailSmtpAddress string) (Maker, error) {

	maker, err := NewSendEmailMaker(frontendUrl, sendEmail, emailPassword, emailSmtp, emailSmtpAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewSendEmailMaker: %w", err)
	}
	return maker, nil
}
