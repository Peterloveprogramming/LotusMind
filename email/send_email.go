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

func (maker *SendEmailMaker) SendChakaraResult(to []string, uniqueCode string, language string) error {
	// //testing SMTP connection
	// conn, connErr := net.Dial("tcp", maker.fromEmailSmtpAddress)
	// if connErr != nil {
	// 	log.Printf("❌ Failed to connect to SMTP server: %v", connErr)
	// } else {
	// 	log.Printf("✅ Connected to SMTP server: %s", maker.fromEmailSmtpAddress)
	// 	conn.Close()
	// }
	var subject string
	var tmpl *template.Template
	var err error

	log.Printf("Preparing to send Chakra report to: %v, Language: %s, UniqueCode: %s", to, language, uniqueCode)

	if language == French {
		subject = SubjectFrench
		tmpl, err = template.New("email").Parse(ChakaraReportTemplateFrench)
	} else if language == English {
		subject = Subject
		tmpl, err = template.New("email").Parse(ChakaraReportTemplate)
	} else {
		log.Printf("Error: Unsupported language: %s", language)
		return fmt.Errorf("unsupported language: %s", language)
	}

	if err != nil {
		log.Printf("Error parsing HTML template: %v", err)
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	reportUrl := maker.frontEndUrl
	if len(to) > 0 {
		reportUrl = fmt.Sprintf("%s/report/%s/%s", maker.frontEndUrl, to[0], uniqueCode)
	}
	log.Printf("Report URL: %s", reportUrl)

	data := SendChakaraResultParams{
		Subject:       Subject,
		ChakaraNumber: uniqueCode,
		DiscountCode:  "VD50R7V0BXZH",
		FrontEndUrl:   reportUrl,
	}

	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("From: %s\n", maker.fromEmail))
	body.WriteString(fmt.Sprintf("To: %s\n", strings.Join(to, ",")))
	body.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	body.WriteString("MIME-version: 1.0;\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n")
	body.WriteString("\n")

	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("Error executing email template: %v", err)
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	log.Printf("Email body prepared. Sending email...")

	err = smtp.SendMail(
		maker.fromEmailSmtpAddress,
		maker.smtpPlanAuth,
		maker.fromEmail,
		to,
		body.Bytes(),
	)

	if err != nil {
		log.Printf("Error sending email: %v", err)
		log.Printf("SMTP server: %s", maker.fromEmailSmtpAddress)
		log.Printf("SMTP sender: %s", maker.fromEmail)
		log.Printf("SMTP recipients: %v", to)
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Printf("Email sent successfully to: %v", to)
	return nil
}
