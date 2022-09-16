package usecase

import (
	"genesis_test_case/src/config"
	"genesis_test_case/src/pkg/domain"
	"genesis_test_case/src/pkg/types/errors"
	"genesis_test_case/src/pkg/usecase"
	mocks "genesis_test_case/src/pkg/usecase/mocks"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestSendCurrencyRate(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Error("Error loading .env file")
	}
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	chart := mocks.NewMockChartProvider(ctl)
	mailer := mocks.NewMockMailingRepository(ctl)
	banner := mocks.NewMockCryptoBannerProvider(ctl)
	storage := mocks.NewMockEmailStorage(ctl)
	exchanger := mocks.NewMockExchangeProvider(ctl)
	mockRepos := &usecase.CryptoMailingRepositories{
		Repositories: usecase.Repositories{
			Chart:     chart,
			Mailer:    mailer,
			Banner:    banner,
			Storage:   storage,
			Exchanger: exchanger,
		},
	}
	BTCUAHPair := &domain.CurrencyPair{
		BaseCurrency:  os.Getenv(config.EnvBaseCurrency),
		QuoteCurrency: os.Getenv(config.EnvQuoteCurrency),
	}

	mailingUsecase := NewCryptoMailingUsecase(
		"./../../../../"+os.Getenv(config.EnvCryptoHtmlMessagePath),
		BTCUAHPair,
		mockRepos,
	)
	mockStorageResp := []string{"example@example.com"}
	mockCryptoChartResp := []float64{0.0, 0.1, 0.2}
	mockCryptoRateResp := &domain.CurrencyRate{}
	mockBannerResp := "http://example.com/example"

	chart.EXPECT().GetWeekAverageChart(BTCUAHPair).Return(mockCryptoChartResp, nil)
	exchanger.EXPECT().GetCurrencyRate(BTCUAHPair).Return(mockCryptoRateResp, nil)
	banner.EXPECT().GetCryptoBannerUrl(mockCryptoChartResp, mockCryptoRateResp).Return(mockBannerResp, nil)
	storage.EXPECT().GetAllEmails().Return(mockStorageResp, nil)
	mailer.EXPECT().MultipleSending(gomock.Any(), mockStorageResp).Return(nil, nil)
	unsent, err := mailingUsecase.SendCurrencyRate()
	require.NoError(t, err)
	require.Nil(t, unsent)
}

func TestSendCurrencyRateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	chart := mocks.NewMockChartProvider(ctl)
	mailer := mocks.NewMockMailingRepository(ctl)
	banner := mocks.NewMockCryptoBannerProvider(ctl)
	storage := mocks.NewMockEmailStorage(ctl)
	exchanger := mocks.NewMockExchangeProvider(ctl)
	mockRepos := usecase.CryptoMailingRepositories{
		Repositories: usecase.Repositories{
			Chart:     chart,
			Mailer:    mailer,
			Banner:    banner,
			Storage:   storage,
			Exchanger: exchanger,
		},
	}
	BTCUAHPair := &domain.CurrencyPair{
		BaseCurrency:  os.Getenv(config.EnvBaseCurrency),
		QuoteCurrency: os.Getenv(config.EnvQuoteCurrency),
	}

	mailingUsecase := NewCryptoMailingUsecase(
		os.Getenv(config.EnvCryptoHtmlMessagePath),
		BTCUAHPair,
		&mockRepos,
	)

	chart.EXPECT().GetWeekAverageChart(BTCUAHPair).Return(nil, errors.ErrFailedParseHttpBody)
	_, err := mailingUsecase.SendCurrencyRate()
	require.Error(t, err)
}
