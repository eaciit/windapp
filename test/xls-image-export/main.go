package main

import (
	"fmt"
	"image"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/tealeg/xlsx"
)

func main() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}

	reader, err := os.Open("img.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	im, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("%#v \n", err)
	}

	rowImage := sheet.AddRow()
	cellImage := rowImage.AddCell()
	cellImage.SetValue(im)

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "I am a cell!"

	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
