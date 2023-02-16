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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)



func main() {

	
	saveDataToSpreadsheet()

}

func saveDataToSpreadsheet() {

		//Create API Context
	ctx := context.Background()
	// Get bytes from base64 encoded google service accounts key
	const (
    client_secret_path = "./credentials/credentials.json"
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
	client := oauth2.getClient(ctx, config)

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

func getClient(ctx context.Context, config *oauth2.Config) {
	panic("unimplemented")
}


