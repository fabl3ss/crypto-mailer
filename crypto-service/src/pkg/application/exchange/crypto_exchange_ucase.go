package usecase

import (
	"errors"
	"genesis_test_case/src/config"
	"genesis_test_case/src/pkg/application"
	"genesis_test_case/src/pkg/domain/models"
	"genesis_test_case/src/pkg/domain/usecases"
	myerr "genesis_test_case/src/pkg/types/errors"
)

type CryptoExchangerUsecase struct {
	cryptoProvider application.ExchangeProvider
	cache          application.CryptoCache
}

func NewCryptoExchangeUsecase(
	crypto application.ExchangeProvider,
	cache application.CryptoCache,
) usecases.CryptoExchangerUsecase {
	return &CryptoExchangerUsecase{
		cryptoProvider: crypto,
		cache:          cache,
	}
}

func (c *CryptoExchangerUsecase) GetCurrentExchangePrice(pair *models.CurrencyPair) (float64, error) {
	cacheRate, err := c.cache.GetCurrencyCache(config.CryptoCacheKey)
	if err != nil {
		if errors.Is(err, myerr.ErrNoCache) {
			rate, err := c.cryptoProvider.GetCurrencyRate(pair)
			if err != nil {
				return -1, err
			}
			err = c.cache.SetCurrencyCache(config.CryptoCacheKey, rate)
			if err != nil {
				return -1, err
			}
			return rate.Price, nil
		}

		return -1, err
	}

	return cacheRate.Price, nil
}
