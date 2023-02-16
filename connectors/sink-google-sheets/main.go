// Copyright 2023 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"os"
	"encoding/json"
	"fmt"
	"net/http"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)



func main() {

	
	saveDataToSpreadsheet()

}

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

func saveDataToSpreadsheet() {

		//Create API Context
	ctx := context.Background()
	
	// Set Credentials Path
	const (
    client_secret_path = "./credentials/client_secret.json"
	)

	credBytes, err := os.ReadFile(client_secret_path)
	if err != nil {
		log.Fatalf("Failed to decode google service accounts key %v", err)
	}

	// authenticate and get configuration
	config, err := google.ConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			log.Fatalf("Failed to authenticate google service accounts key %v", err)
			return
		}

	//Create Client
	client := getClient(config)

	//Create Service using Client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
        log.Fatalf("Failed to Create Service Account %v",err)
        return
	}

	//Initialize Sheet ID & Spreadsheet ID
	sheetId := 0
	spreadSheetId := "1tZJPUCOiiR0liRsNtLKhCoQR-Cb8_oPVGMU0kvnRCQw"

	//Get SheetName from SpreadSheetID
	res1, err := srv.Spreadsheets.Get(spreadSheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil {
        log.Fatalf("Failed to get SheetName %v",err)
        return
	}

	sheetName := ""
	for _, v := range res1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			break
		}
	}

	//Append value to Spreadsheet

	row := &sheets.ValueRange{
	Values: [][]interface{}{{"1", "ABC", "abc@gmail.com", "16-02-2023"}},
	}

	response2, err := srv.Spreadsheets.Values.Append(spreadSheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
		if err != nil || response2.HTTPStatusCode != 200 {
		log.Fatalf("Failed to Append Value to Spreadsheet %v",err)
		return
	}
}


