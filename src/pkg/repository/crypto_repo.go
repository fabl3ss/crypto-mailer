package repository

import (
	"encoding/json"
	"fmt"
	"genesis_test_case/src/config"
	"genesis_test_case/src/pkg/domain"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type cryptoRepository struct{}

func NewCryptoRepository() domain.CryptoRepository {
	return &cryptoRepository{}
}

func (c *cryptoRepository) GetCandles(cfg *config.Config, candleProps *domain.CandleProps) ([][]float64, error) {
	url := fmt.Sprintf(
		cfg.CryptoApiCandlesUrl,
		candleProps.Base,
		candleProps.Granularity,
		candleProps.Start,
		candleProps.End,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}

	var candles [][]float64

	err = json.Unmarshal(body, &candles)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal JSON")
	}
	return candles, nil
}

func (c *cryptoRepository) GetCurrencyRate(base string, quoted string) (*domain.CurrencyRate, error) {
	cfg := config.Get()
	url := fmt.Sprintf(
		cfg.CryptoApiFormatUrl,
		base,
		quoted,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "can't make http request")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}

	data := struct {
		Rate domain.CurrencyRate `json:"data"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.Wrap(err, "can't unmarshal JSON")
	}

	data.Rate.Price = strings.Split(data.Rate.Price, ".")[0]
	return &data.Rate, nil
}
