package gmail_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client
func getClient(config *oauth2.Config) *http.Client {
	tokFile := os.Getenv("GMAIL_TOKEN_PATH")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		saveToken(tokFile, getTokenFromWeb(config))
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode, oauth2.AccessTypeOffline)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func UpdateService(credentialsPath string, tokenPath string) (*gmail.Service, error) {
	config, err := getClientFromFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	tok, err := tokenFromFile(tokenPath)
	tokenSource := config.TokenSource(context.TODO(), tok)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if newToken.AccessToken != tok.AccessToken {
		saveToken(tokenPath, newToken)
	}
	return GetGmailService()
}

func getClientFromFile(path string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read client secret file")
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse client secret file to config")
	}
	return config, nil
}

func GetGmailService() (*gmail.Service, error) {
	config, err := getClientFromFile(os.Getenv("GMAIL_CREDENTIALS_PATH"))
	if err != nil {
		return nil, err
	}

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(getClient(config)))
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve Gmail client")
	}
	return srv, nil
}
