package domain

import "google.golang.org/api/gmail/v1"

type Recipient struct {
	Email string `json:"email" validate:"required,email"`
}

type MailingUsecase interface {
	SendCurrencyRate() ([]string, error)
	Subscribe(recipient *Recipient) error
}

type MailingRepository interface {
	SendMessage(message *gmail.Message, body string) error
}
