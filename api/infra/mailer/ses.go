package mailer

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type sesMailer struct {
	ses    *sesv2.Client
	sender string
}

func NewSES() EmailSender {
	sender := "myLocal Software <info@mylocal.ing>"
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("ses init: %v", err)
	}
	return &sesMailer{ses: sesv2.NewFromConfig(cfg), sender: sender}
}

func (m *sesMailer) SendSignupConfirmation(toEmail, code string) error {
	subject := "Welcome to myLocal!"
	text := fmt.Sprintf(`Verify your account with this code: %s`, code)
	html := fmt.Sprintf("<strong>Verify your account with this code: %s</strong>", code)

	_, err := m.ses.SendEmail(context.TODO(), &sesv2.SendEmailInput{
		FromEmailAddress: &m.sender,
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Html: &types.Content{Data: aws.String(html)},
					Text: &types.Content{Data: aws.String(text)},
				},
			},
		},
	})
	return err
}
