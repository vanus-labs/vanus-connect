package main

import (
	"fmt"
	"strconv"
	"context"
	"log"
	"os"
	b64 "encoding/base64"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)



func main() {

	// Store Environmental Variable
	storeEnvData()

	//Run Business Logic
	saveDataToSpreadsheet()

}

func storeEnvData() {
	
    data, err := os.ReadFile("credentials.json")
    if err!= nil {
        log.Fatalf("Failed to read file %v", err)
    }
    
    credentials := b64.StdEncoding.EncodeToString([]byte(data))
    os.Setenv("KEY_JSON_BASE64", credentials)
    
    
}


func saveDataToSpreadsheet() {

		//Create API Context
	ctx := context.Background()


	//Decode Auth Key
	credBytes, err := b64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	if err != nil {
		log.Fatalf("Failed to decode google service accounts key %v", err)
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			log.Fatalf("Failed to authenticate google service accounts key %v", err)
			return
		}

	//Create Client
	client := config.Client(ctx)

	//Create Service using Client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
        log.Fatalf("Failed to Create Service Account %v",err)
        return
	}

	//Initialize Sheet ID & Spreadsheet ID
	//spreadSheetUrl := "https://docs.google.com/spreadsheets/d/1tZJPUCOiiR0liRsNtLKhCoQR-Cb8_oPVGMU0kvnRCQw/edit#gid=0"
	fmt.Println("Insert Your Spreadsheet URL")
	var spreadSheetUrl string
	fmt.Scanf("%v \n", &spreadSheetUrl)

	sheetId, err := strconv.Atoi(spreadSheetUrl[93:94])
	if err != nil {
        log.Fatalf("Failed to Convert String %v",err)
        return
	}

	spreadSheetId := spreadSheetUrl[39:83]

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

	fmt.Println("Insert ID")
	var row_id string
    fmt.Scanf("%v \n", &row_id)
	
	fmt.Println("Insert Name")
	var name string
    fmt.Scanf("%v \n", &name)

	fmt.Println("Insert Email")
	var email string
    fmt.Scanf("%v \n", &email)

	fmt.Println("Insert Date - DD/MM/YYYY")
	var date string
    fmt.Scanf("%v \n", &date)

	row := &sheets.ValueRange{
	Values: [][]interface{}{{row_id, name, email, date}},
	}

	response2, err := srv.Spreadsheets.Values.Append(spreadSheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
		if err != nil || response2.HTTPStatusCode != 200 {
		log.Fatalf("Failed to Append Value to Spreadsheet %v",err)
		return
	}
}