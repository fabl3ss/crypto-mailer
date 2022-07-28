package usecase

import (
	"genesis_test_case/src/config"
	"genesis_test_case/src/pkg/domain"
	"genesis_test_case/src/pkg/repository"
	"strconv"
)

type cryptoUsecase struct {
	repos *repository.Repositories
}

func NewCryptoUsecase(r *repository.Repositories) domain.CryptoUsecase {
	return &cryptoUsecase{
		repos: r,
	}
}

func (c *cryptoUsecase) GetConfigCurrencyRate() (int, error) {
	cfg := config.Get()
	rate, err := c.repos.Crypto.GetCurrencyRate(cfg.BaseCurrency, cfg.QuoteCurrency)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(rate.Price)
}
