package main

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Com1Software/go-dbase/dbase"
	"golang.org/x/text/encoding/charmap"
)

type List struct {
	Item string `dbase:"ITEM"`
}

func main() {
	TableCheck()
	a := app.New()
	w := a.NewWindow("Shopping List")

	// Open the table
	table, err := dbase.OpenTable(&dbase.Config{
		Filename:   "LIST.DBF",
		TrimSpaces: true,
	})
	if err != nil {
		panic(err)
	}
	defer table.Close()

	// Create a slice to hold the items
	var items []string
	loadItems := func() {
		items = []string{}   // Clear existing items
		err := table.GoTo(0) // Start from the first record
		if err != nil {
			fmt.Println("Error navigating to the first record:", err)
			return
		}
		for {
			row, err := table.Next()
			if err != nil {
				fmt.Println("Error reading row:", err)
				break
			}
			field := row.FieldByName("ITEM")
			if field != nil {
				items = append(items, fmt.Sprintf("%v", field.GetValue()))
			}
			if table.EOF() {
				break // Exit when reaching the end of the file
			}
		}
	}

	// Load initial items
	loadItems()

	// Create a list widget
	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(items[id])
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selectedItem := items[id]

		// Create a dialog to edit the item
		entry := widget.NewEntry()
		entry.SetText(selectedItem)

		dialog.ShowCustomConfirm(
			"Edit Item",
			"Save",
			"Cancel",
			container.NewVBox(entry),
			func(confirm bool) {
				if confirm {
					newValue := entry.Text

					// Update the DBF file
					err := table.GoTo(uint32(id)) // Navigate to the selected record
					if err != nil {
						fmt.Println("Error navigating to record:", err)
						return
					}
					row, err := table.Row()
					if err != nil {
						fmt.Println("Error retrieving row:", err)
						return
					}
					err = row.FieldByName("ITEM").SetValue(newValue)
					if err != nil {
						fmt.Println("Error updating value:", err)
					}
					err = row.Write()
					if err != nil {
						fmt.Println("Error saving changes:", err)
					}

					// Refresh the items list
					loadItems()
					list.Refresh()
				}
			},
			w,
		)
	}

	scrollableList := container.NewScroll(list)
	scrollableList.SetMinSize(fyne.NewSize(600, 400)) // Set the minimum size

	// Entry for new item
	memo := widget.NewEntry()
	memo.SetPlaceHolder("Enter your memo here...")

	// Button to add a new item
	addButton := widget.NewButton("Add Item", func() {
		table, err := dbase.OpenTable(&dbase.Config{
			Filename:   "LIST.DBF",
			TrimSpaces: true,
		})
		if err != nil {
			panic(err)
		}
		defer table.Close()
		recno := "0"
		rn, _ := strconv.Atoi(recno)
		err = table.GoTo(uint32(rn))
		if err != nil {
			panic(err)
		}
		row, err := table.Row()
		if err != nil {
			panic(err)
		}
		err = row.FieldByName("ITEM").SetValue(memo.Text)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = row.Write()
		if err != nil {
			fmt.Println(err.Error())
		}
		memo.SetText("") // Clear entry field
		loadItems()      // Refresh items
		list.Resize(fyne.NewSize(600, 400))
		list.Refresh() // Update list view
	})

	// Exit button
	exitButton := widget.NewButton("Exit", func() {
		os.Exit(0)
	})

	// Layout
	content := container.NewVBox(
		scrollableList,
		memo,
		addButton,
		exitButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 400))
	w.ShowAndRun()
}

func TableCheck() {
	tt := "LIST.DBF"
	if _, err := os.Stat(tt); err == nil {
		// File exists, do nothing
	} else {
		// Create a new DBF file
		file, err := dbase.NewTable(
			dbase.FoxProAutoincrement,
			&dbase.Config{
				Filename:   tt,
				Converter:  dbase.NewDefaultConverter(charmap.Windows1250),
				TrimSpaces: true,
			},
			icolumns(),
			64,
			nil,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		row, err := file.RowFromStruct(&List{
			Item: "ITEM",
		})
		if err != nil {
			panic(err)
		}

		err = row.Add()
		if err != nil {
			panic(err)
		}
		fmt.Printf(
			"Last modified: %v Columns count: %v Record count: %v File size: %v \n",
			file.Header().Modified(0),
			file.Header().ColumnsCount(),
			file.Header().RecordsCount(),
			file.Header().FileSize(),
		)
	}
}

func icolumns() []*dbase.Column {
	itemCol, err := dbase.NewColumn("ITEM", dbase.Varchar, 80, 0, false)
	if err != nil {
		panic(err)
	}
	return []*dbase.Column{
		itemCol,
	}
}
