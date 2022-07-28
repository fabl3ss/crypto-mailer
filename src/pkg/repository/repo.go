package repository

import (
	domain2 "genesis_test_case/src/pkg/domain"
	"genesis_test_case/src/platform/gmail_api"
)

type Repositories struct {
	Image   domain2.ImageRepository
	Crypto  domain2.CryptoRepository
	Mailing domain2.MailingRepository
}

func NewRepositories() (*Repositories, error) {
	gmailClient, err := gmail_api.GetGmailService()
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Image:   NewImageRepository(),
		Crypto:  NewCryptoRepository(),
		Mailing: NewMailingRepository(gmailClient),
	}, nil
}
