package usecase

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"genesis_test_case/src/config"
	domain2 "genesis_test_case/src/pkg/domain"
	"genesis_test_case/src/pkg/repository"
	"genesis_test_case/src/pkg/types/errors"
	"genesis_test_case/src/pkg/utils"
	"google.golang.org/api/gmail/v1"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type mailingUsecase struct {
	repos *repository.Repositories
}

func NewMailingUsecase(r *repository.Repositories) domain2.MailingUsecase {
	return &mailingUsecase{
		repos: r,
	}
}

func (m *mailingUsecase) Subscribe(recipient *domain2.Recipient) error {
	cfg := config.Get()

	emails, err := utils.ExistsInCsv(cfg.StorageFile, recipient.Email)
	if err != nil {
		return err
	} else if emails == nil {
		return errors.ErrAlreadyExists
	}

	f, err := os.OpenFile(
		config.Get().StorageFile,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, v := range emails {
		err := w.Write([]string{v})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mailingUsecase) SendCurrencyRate() ([]string, error) {
	cfg := config.Get()

	rate, err := m.repos.Crypto.GetCurrencyRate(cfg.BaseCurrency, cfg.QuoteCurrency)
	if err != nil {
		return nil, err
	}

	nowUtc := time.Now().UTC()
	candleProps := &domain2.CandleProps{
		Base:        cfg.BaseCurrency,
		Granularity: strconv.Itoa(int(time.Hour.Seconds())),
		Start:       nowUtc.AddDate(0, 0, -7).Format(time.RFC3339),
		End:         nowUtc.Format(time.RFC3339),
	}

	candles, err := m.repos.Crypto.GetCandles(cfg, candleProps)
	if err != nil {
		return nil, err
	}

	// Compute avg of each candle for chart
	var chartCandles []float64

	// Do reverse loop because GetCandle returns slice
	// in descending order by date
	for i := len(candles) - 1; i >= 0; i-- {
		// candles[i][3] -> opening price (first trade) in the bucket interval
		chartCandles = append(chartCandles, candles[i][3])
	}

	url, err := m.repos.Image.GetCryptoBannerUrl(
		os.Getenv("BANNER_API_URL"),
		os.Getenv("BANNER_API_TOKEN"),
		chartCandles,
		rate,
	)
	if err != nil {
		return nil, err
	}

	// New message for our gmail service to send
	var message gmail.Message

	v := struct {
		Chart string
	}{Chart: url}

	var tpl bytes.Buffer

	t, _ := template.ParseFiles(
		os.Getenv("CRYPTO_MESSAGE_HTML"),
	)
	if err := t.Execute(&tpl, v); err != nil {
		return nil, err
	}

	// Gmail message declaration
	to := "To: %s\r\n"
	messageStr :=
		"Subject: Crypto Newsletter\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"Content-Transfer-Encoding: base64\r\n\r\n" +
			tpl.String()

	f, err := os.Open(cfg.StorageFile)
	if err != nil {
		return nil, err
	}

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)

	var unsent []string
	for {
		rec, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		err = m.repos.Mailing.SendMessage(&message, fmt.Sprintf(to, rec[0])+messageStr)
		if err != nil {
			log.Printf("Sending error!: %v\n", err)
			unsent = append(unsent, rec[0])
		} else {
			log.Printf("Message sent! to %s\n", rec[0])
		}
	}
	return unsent, nil
}
