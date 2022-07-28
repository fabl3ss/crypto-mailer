package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	domain2 "genesis_test_case/src/pkg/domain"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

type imageRepository struct{}

func NewImageRepository() domain2.ImageRepository {
	return &imageRepository{}
}

func (i *imageRepository) extractBannerUrl(url string, bearer string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+bearer)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	jpg := struct {
		Url string `json:"image_url_jpg"`
	}{}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &jpg)
	if err != nil {
		return "", errors.Wrap(err, "can't unmarshal JSON")
	} else if jpg.Url == "" {
		return "", nil
	}
	return jpg.Url, nil
}

func (i *imageRepository) GetCryptoBannerUrl(bannerUrl, bearer string, chart []float64, rate *domain2.CurrencyRate) (string, error) {
	postBody, _ := json.Marshal(&domain2.BannerRequest{
		Template:    os.Getenv("CRYPTO_TEMPLATE"),
		Transparent: false,
		WebHookUrl:  nil,
		Metadata:    nil,
		Modifications: []domain2.BannerModifications{
			{
				Name:       "bg",
				Background: nil,
			},
			{
				Name:       "title",
				Text:       fmt.Sprintf("%s/%s", rate.BaseCurrency, rate.QuoteCurrency),
				Background: nil,
			},
			{
				Name:       "subtitle",
				Text:       fmt.Sprintf("1 %s = %s %s", rate.BaseCurrency, rate.Price, rate.QuoteCurrency),
				Background: nil,
			},
			{
				Name:      "chart",
				ChartData: chart,
			},
		},
	})

	req, err := http.NewRequest("POST", bannerUrl, bytes.NewBuffer(postBody))
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	b := struct {
		Self string `json:"self"`
	}{}

	err = json.Unmarshal(body, &b)
	if err != nil {
		return "", err
	}

	var jpgUrl string
	for jpgUrl == "" {
		// Wait for image generation
		time.Sleep(100 * time.Millisecond)
		jpgUrl, err = i.extractBannerUrl(b.Self, bearer)
		if err != nil {
			return "", err
		}
	}

	return jpgUrl, nil
}
