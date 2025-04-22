package main

import (
	"context"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	// Button to add a new item to the shopping list using a dialog box
	addItemButton := widget.NewButton("Add Item to Shopping List", func() {
		// Create a new entry field for the dialog
		newItemEntry := widget.NewEntry()

		// Create the dialog
		dialogBox := dialog.NewCustomConfirm(
			"Add Item",
			"Add",
			"Cancel",
			container.NewVBox(
				widget.NewLabel("Enter the item to add:"),
				newItemEntry,
			),
			func(confirmed bool) {
				if confirmed {
					ctx := context.Background()

					// Load the service account key file
					b, err := os.ReadFile("C:/Users/infor/Documents/dependable-glow-836-65f6fa83b621.json")
					if err != nil {
						dataDisplay.SetText(fmt.Sprintf("Error reading service account key file: %v", err))
						return
					}

					// Authenticate and create a Sheets service
					config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
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

					// Append data to the spreadsheet
					spreadsheetId := "1uZTXl8XP6VaZII2wtG0oMZFLyEaGqRw7nuEVAon3iRQ"
					writeRange := "ShoppingList!A:A"
					valueRange := &sheets.ValueRange{
						Values: [][]interface{}{
							{newItemEntry.Text},
						},
					}
					_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, valueRange).ValueInputOption("RAW").Do()
					if err != nil {
						dataDisplay.SetText(fmt.Sprintf("Error adding item to spreadsheet: %v", err))
						return
					}

					// Confirm item added
					dataDisplay.SetText(fmt.Sprintf("Added '%s' to the shopping list!", newItemEntry.Text))
				}
			},
			w,
		)

		// Show the dialog
		dialogBox.Show()
	})

	// Set up the layout
	w.SetContent(container.NewVBox(
		dataDisplay,
		loadDataButton,
		addItemButton,
		widget.NewButton("Exit", func() { os.Exit(0) }),
	))

	// Resize and run the application
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
