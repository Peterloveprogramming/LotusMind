package sendemail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"text/template"
)

type SendEmailMaker struct {
	fromEmail            string
	fromEmailSmtpAddress string
	smtpPlanAuth         smtp.Auth
	frontEndUrl          string
}

func NewSendEmailMaker(frontendUrl string, fromEmail string, fromEmailPassword string, fromEmailSmtp string, fromEmailSmtpAddress string) (Maker, error) {
	// Check if any required parameter is empty
	if fromEmail == "" || fromEmailPassword == "" || fromEmailSmtp == "" || fromEmailSmtpAddress == "" {
		return nil, fmt.Errorf("missing required email configuration: email,emailPassword, emailSmtp, or emailSmtpAddress cannot be empty")
	}

	// Create email auth
	auth := smtp.PlainAuth(
		"",
		fromEmail,
		fromEmailPassword,
		fromEmailSmtp,
	)

	// Create the maker instance
	maker := &SendEmailMaker{
		fromEmail:            fromEmail,
		fromEmailSmtpAddress: fromEmailSmtpAddress,
		smtpPlanAuth:         auth,
		frontEndUrl:          frontendUrl,
	}
	return maker, nil
}

type SendChakaraResultParams struct {
	Subject       string
	ChakaraNumber string
	DiscountCode  string
	FrontEndUrl   string
}

func (maker *SendEmailMaker) SendChakaraResult(to []string, uniqueCode string) error {
	subject := Subject
	tmpl, err := template.New("email").Parse(ChakaraReportTemplate)
	if err != nil {
		log.Printf("Error parsing HTML template: %v", err)
		return fmt.Errorf("failed to parse email template: %w", err)
	}
	reportUrl := maker.frontEndUrl // Start with the base URL
	fmt.Println("reportUrl", reportUrl)
	if len(to) > 0 {
		reportUrl = fmt.Sprintf("%s/report/%s/%s", maker.frontEndUrl, to[0], uniqueCode)
	}

	data := SendChakaraResultParams{
		Subject:       Subject,
		ChakaraNumber: uniqueCode,
		DiscountCode:  "CHAKRATEST15OFF",
		FrontEndUrl:   reportUrl,
	}

	// 4. Execute the template into a buffer
	var body bytes.Buffer
	// Set headers first, including Content-Type
	body.WriteString(fmt.Sprintf("From: %s\n", maker.fromEmail))
	body.WriteString(fmt.Sprintf("To: %s\n", strings.Join(to, ",")))
	body.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	body.WriteString("MIME-version: 1.0;\n")                          // Specify MIME version
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n") // Set Content-Type to HTML
	body.WriteString("\n")

	// Execute template and write HTML body
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("failed to parse email template: %v", err)
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	err = smtp.SendMail(
		maker.fromEmailSmtpAddress,
		maker.smtpPlanAuth,
		maker.fromEmail,
		to,
		body.Bytes(),
	)

	if err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("Error sending email:: %w", err)
	}

	return nil
}
