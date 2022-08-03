package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

const tokenFileName = "token.json"
const secretFileName = "credentials.json"
const configFileName = "config.json"

// ClientDetails is a simple struct that contains client related information
type ClientDetails struct {
	id     string
	secret string
}

// Config is the config object
type Config struct {
	CalendarName string
}

// GetFileAndPath get's the full filepath and the path for the given file
func GetFileAndPath(fileName string) (string, string) {
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path += "/.wms"
	file := path + "/" + fileName
	return file, path
}

// CheckToken checks if a token exists and can be decoded
func CheckToken() error {
	file, _ := GetFileAndPath(tokenFileName)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return err
}

// CreateOauth2Config creates an OAuth 2 Config object
func CreateOauth2Config(clientID string, clientSecret string) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: 0,
		},
		RedirectURL: "http://localhost",
		Scopes:      []string{calendar.CalendarReadonlyScope},
	}

	file, path := GetFileAndPath(secretFileName)
	err := os.MkdirAll(path, 0755)

	fmt.Printf("Saving credential file to: %s\n", file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache credentials: %v", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(config)

	return config
}

// GetTokenFromWeb Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// SaveToken Saves a token to a file path.
func SaveToken(token *oauth2.Token) {
	file, path := GetFileAndPath(tokenFileName)
	err := os.MkdirAll(path, 0755)

	fmt.Printf("Saving credential file to: %s\n", file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	file, _ := GetFileAndPath(tokenFileName)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func configFromFile() (*oauth2.Config, error) {
	file, _ := GetFileAndPath(secretFileName)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	config := &oauth2.Config{}
	err = json.NewDecoder(f).Decode(config)
	return config, err
}

// GetClient Retrieve a token, saves the token, then returns the generated client.
func GetClient() *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	config, err := configFromFile()
	if err != nil {
		log.Fatal("Unable to load token, have you tried 'wms initialize'?")
	}

	tok, err := tokenFromFile()
	if err != nil {
		log.Fatal("Unable to load config, have you tried 'wms initialize'?")
	}

	return config.Client(context.Background(), tok)
}

// CreateConfig creates a config file that stores configuration options
func CreateConfig() {
	config := &Config{
		CalendarName: "default",
	}

	file, _ := GetFileAndPath(configFileName)

	fmt.Printf("Saving config file to: %s\n", file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to create config: %v", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(config)
}

// LoadConfig loads in a config file and returns a config object
func LoadConfig() (*Config, error) {
	file, _ := GetFileAndPath(configFileName)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := &Config{}
	err = json.NewDecoder(f).Decode(config)

	return config, err
}
