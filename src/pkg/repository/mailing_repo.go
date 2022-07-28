package repository

import (
	"encoding/base64"
	"genesis_test_case/src/pkg/domain"
	"genesis_test_case/src/platform/gmail_api"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/api/gmail/v1"
)

type mailingRepository struct {
	gmailService *gmail.Service
}

func NewMailingRepository(srv *gmail.Service) domain.MailingRepository {
	return &mailingRepository{
		gmailService: srv,
	}
}

func (m *mailingRepository) updateOauthToken() error {
	service, err := gmail_api.UpdateService(
		os.Getenv("GMAIL_CREDENTIALS_PATH"),
		os.Getenv("GMAIL_TOKEN_PATH"),
	)
	if err != nil {
		return err
	}

	// Update for new service with fresh token
	m.gmailService = service
	return nil
}

func (m *mailingRepository) SendMessage(message *gmail.Message, body string) error {
	if err := m.updateOauthToken(); err != nil {
		return err
	}

	// Encode message
	message.Raw = base64.URLEncoding.EncodeToString([]byte(body))

	// Send email
	_, err := m.gmailService.Users.Messages.Send("me", message).Do()
	if err != nil {
		return errors.Wrap(err, "sending failed")
	}

	return nil
}
