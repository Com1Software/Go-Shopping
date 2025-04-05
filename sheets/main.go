package main

import (
	"context"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func main() {
	// Create the application
	a := app.New()
	w := a.NewWindow("Spreadsheet Reader")

	// Memo field to display spreadsheet data
	dataDisplay := widget.NewMultiLineEntry()
	dataDisplay.SetPlaceHolder("Shopping List will appear here...")

	// Load data from Google Sheets
	loadDataButton := widget.NewButton("Load Shopping List", func() {
		ctx := context.Background()

		// Load the service account key file
		b, err := os.ReadFile("C:/Users/infor/Documents/dependable-glow-836-65f6fa83b621.json")
		if err != nil {
			dataDisplay.SetText(fmt.Sprintf("Error reading service account key file: %v", err))
			return
		}

		// Authenticate and create a Sheets service
		config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
		if err != nil {
			dataDisplay.SetText(fmt.Sprintf("Error parsing service account key file: %v", err))
			return
		}

		client := config.Client(ctx)
		srv, err := sheets.New(client)
		if err != nil {
			dataDisplay.SetText(fmt.Sprintf("Error creating Sheets client: %v", err))
			return
		}

		// Spreadsheet ID and range
		spreadsheetId := "1uZTXl8XP6VaZII2wtG0oMZFLyEaGqRw7nuEVAon3iRQ"
		readRange := "ShoppingList!A1:A99"

		// Retrieve data
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
		if err != nil {
			dataDisplay.SetText(fmt.Sprintf("Error reading spreadsheet data: %v", err))
			return
		}

		if len(resp.Values) == 0 {
			dataDisplay.SetText("No data found.")
		} else {
			var data string
			for _, row := range resp.Values {
				data += fmt.Sprintf("%s\n", row)
			}
			dataDisplay.SetText(data)
		}
	})

	// Set up the layout
	w.SetContent(container.NewVBox(
		dataDisplay,
		loadDataButton,
		widget.NewButton("Exit", func() { os.Exit(0) }),
	))

	// Resize and run the application
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
