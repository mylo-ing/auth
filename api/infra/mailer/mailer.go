package mailer

// EmailSender is what the controllers rely on.
type EmailSender interface {
	SendSignupConfirmation(toEmail, code string) error
}
