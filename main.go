package main

import (
	"fmt"
	"os"

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
	w := a.NewWindow("Two Buttons with Memo")
	memo := widget.NewEntry()
	memo.SetPlaceHolder("Enter your memo here...")
	memo.MultiLine = true
	memo.Resize(fyne.NewSize(400, 100))

	helloButton := widget.NewButton("Say Hello", func() {
		dialog.ShowInformation("Hello", "Hello, "+memo.Text, w)
	})
	exitButton := widget.NewButton("Exit", func() {
		os.Exit(0)
	})
	w.SetContent(container.NewVBox(
		memo,
		helloButton,
		exitButton,
	))
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}

func TableCheck() {
	tt := "LIST.DBF"
	if _, err := os.Stat(tt); err == nil {

	} else {

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
