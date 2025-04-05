package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func main() {
	ctx := context.Background()

	// Load the service account key file
	// b, err := os.ReadFile("path/to/your-service-account-key.json")
	b, err := os.ReadFile("C:/Users/infor/Documents/dependable-glow-836-65f6fa83b621.json")
	if err != nil {
		log.Fatalf("Unable to read service account key file: %v", err)
	}

	// Authenticate and create a Sheets service
	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse service account key file to config: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Spreadsheet ID and range
	spreadsheetId := ""
	readRange := "To do!A4:C4"

	// Read data
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			fmt.Println(row)
		}
	}
}
