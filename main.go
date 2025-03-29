package main

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	table, err := dbase.OpenTable(&dbase.Config{
		Filename:   "LIST.DBF",
		TrimSpaces: true,
	})
	if err != nil {
		panic(err)
	}
	defer table.Close()
	recno := 0
	for !table.EOF() {
		row, err := table.Next()
		if err != nil {
			panic(err)
		}
		field := row.Field(0)
		if field == nil {
			panic("Field not found")
		}
		s := fmt.Sprintf("%v", field.GetValue())
		fmt.Println(s)
		recno++

	}
	memo := widget.NewEntry()
	memo.SetPlaceHolder("Enter your memo here...")
	memo.MultiLine = true
	memo.Resize(fyne.NewSize(400, 100))

	helloButton := widget.NewButton("Add Item", func() {

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
