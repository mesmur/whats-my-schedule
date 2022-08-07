package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

const tokenFileName = "token.json"
const secretFileName = "credentials.json"

// ClientDetails is a simple struct that contains client related information
type ClientDetails struct {
	id     string
	secret string
}

// Config is the config object
type Config struct {
	CalendarName string `mapstructure:"calendar_name"`
}

// GetFileAndPath get's the full filepath and the path for the given file
func GetFileAndPath(fileName string) string {
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	file := path + "/.wms/" + fileName
	return file
}

// CheckToken checks if a token exists and can be decoded
func CheckToken() error {
	file := GetFileAndPath(tokenFileName)
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
		RedirectURL: "http://localhost:8081",
		Scopes:      []string{calendar.CalendarReadonlyScope},
	}

	file := GetFileAndPath(secretFileName)
	fmt.Printf("Saving credential to: %s\n", file)

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
	fmt.Printf(
		"Go to the following link in your browser and authorize the app! \n%v\n",
		authURL,
	)

	browser.OpenURL(authURL)

	var authCode string

	srv := &http.Server{Addr: ":8081"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode = r.URL.Query().Get("code")
		srv.Close()
	})

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token
}

// SaveToken Saves a token to a file path.
func SaveToken(token *oauth2.Token) {
	file := GetFileAndPath(tokenFileName)
	fmt.Printf("Saving token to: %s\n", file)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Retrieves a token from a local file.
func tokenFromFile() (*oauth2.Token, error) {
	file := GetFileAndPath(tokenFileName)

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
	file := GetFileAndPath(secretFileName)
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

// Timer returns a function that prints the name argument and
// the elapsed time between the call to timer and the call to
// the returned function. The returned function is intended to
// be used in a defer statement:
//
//   defer timer("sum")()
func Timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
