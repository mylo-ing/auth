package email

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Mailer struct {
	ses    *sesv2.Client
	sender string
}

func New() *Mailer {
	sender := "myLocal Software <info@mylocal.ing>"
	cfg, err := config.LoadDefaultConfig(context.TODO()) // reads env + ~/.aws/*
	if err != nil {
		log.Fatalf("ses init: %v", err)
	}
	return &Mailer{
		ses:    sesv2.NewFromConfig(cfg),
		sender: sender,
	}
}

func (m *Mailer) SendSignupConfirmation(toEmail, code string) error {
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
