package sheets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string = "4/2wEbk1htf3gsMzYupfp0U1wnDrM_QDoKgTAhJN6Soeb8hyXM38Ds3M8"
	/*
		if _, err := fmt.Scan(&authCode); err != nil {
			log.Fatalf("Unable to read authorization code: %v", err)
		}
	*/
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func NewService() *Service {

	newService := &Service{}
	err := newService.init()
	if err != nil {
		log.Fatal(err)
	}
	return newService

}

func (s *Service) init() error {
	//TODO make this configurable
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		errors.Wrap(err, "Unable to parse client secret file to config")
	}

	client := getClient(config)

	s.srv, err = sheets.New(client)
	if err != nil {
		errors.Wrap(err, "Unable to retrieve Sheets client:")
	}
	return nil
}

func (s *Service) create() error {
	if s.currentSheet != "" {
		return nil
	}
	sSheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "go-github-license-check",
		},
	}
	ctx := context.Background()

	ss, err := s.srv.Spreadsheets.Create(sSheet).Context(ctx).Do()
	if err != nil {
		return errors.Wrap(err, "sheets.create() - ")
	}

	s.currentSheet = ss.SpreadsheetId
	return nil
}

// Save the Results
func (s *Service) Save(input license.LicenseCheckResults) error {

	if err := s.create(); err != nil {
		return err
	}

	//var vr sheets.ValueRange = licenseResultstoValueRange(input)

	var vr sheets.ValueRange

	var newSheet map[string]license.LicenseCheckResults = make(map[string]license.LicenseCheckResults)
	var which int = 0
	for _, item := range input {

		if _, ok := newSheet[item.GitHubRepo]; !ok {
			which++
			req := sheets.Request{
				AddSheet: &sheets.AddSheetRequest{
					Properties: &sheets.SheetProperties{
						Title: item.GitHubRepo,
					},
				},
			}

			rbb := &sheets.BatchUpdateSpreadsheetRequest{
				Requests: []*sheets.Request{&req},
			}

			call := s.srv.Spreadsheets.BatchUpdate(s.currentSheet, rbb)
			_, err := call.Do()
			if err != nil {
				return err
			}

		}

		newSheet[item.GitHubRepo] = append(newSheet[item.GitHubRepo], item)

	}

	for sheetName, items := range newSheet {
		myval := []interface{}{"GitHubRepo", "Dependency", "Lang", "DependencyLicense"}
		vr.Values = append(vr.Values, myval)
		for _, item := range items {
			myval = []interface{}{item.GitHubRepo, item.Dependency, item.Lang, item.DependencyLicense}
			vr.Values = append(vr.Values, myval)
		}
		call := s.srv.Spreadsheets.Values.Update(s.currentSheet, sheetName, &vr)
		call.ValueInputOption("RAW").Do()
	}

	return nil

}

func licenseResultstoValueRange(input license.LicenseCheckResults) sheets.ValueRange {
	var vr sheets.ValueRange
	myval := []interface{}{"GitHubRepo", "Dependency", "Lang", "DependencyLicense"}
	vr.Values = append(vr.Values, myval)
	for _, item := range input {
		myval := []interface{}{item.GitHubRepo, item.Dependency, item.Lang, item.DependencyLicense}
		vr.Values = append(vr.Values, myval)
	}

	return vr
}
