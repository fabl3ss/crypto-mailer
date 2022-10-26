package exchangers

import (
	"genesis_test_case/src/pkg/application"
	"genesis_test_case/src/pkg/domain/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetCurrencyRateTest(exchanger application.ExchangeProvider, t *testing.T) {
	pair := models.NewCurrencyPair(
		"BTC",
		"UAH",
	)

	rate, err := exchanger.GetCurrencyRate(pair)

	require.NoError(t, err)
	require.Equal(t, pair.GetBaseCurrency(), rate.GetBaseCurrency())
	require.Equal(t, pair.GetQuoteCurrency(), rate.GetQuoteCurrency())
	require.NotEmpty(t, rate.Price)
}
