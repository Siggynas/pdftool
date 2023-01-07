package main

import (

	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

func main() {
	pdf.DebugOn = true
	//content, err := readPdf2("test.pdf")
	content, err := readPdfRow("test.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
	return
}

func readPdf(path string) (string, error){
	file,r,err := pdf.Open(path)
	defer file.Close()
	if err != nil {
		return "",err
	}

	var buf bytes.Buffer
	b,err := r.GetPlainText()
	if err != nil {
		return "",err
	}

	buf.ReadFrom(b)

	return buf.String(),nil
}

//Cette fonction m'intéresse, décryptage
func readPdf2(path string) (string, error) {

	//Ouverture du fichier pdf
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	
	//Récupération du nombre de pages
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
		}
	}
	return "", nil
}
func isSameSentence(text pdf.Text, lastTextStyle pdf.Text) bool {
	return (text.Font == lastTextStyle.Font) && (text.FontSize == lastTextStyle.FontSize) && (text.X == lastTextStyle.X)
}


func readPdfRow(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
		    println(">>>> row: ", row.Position)
		    for _, word := range row.Content {
		        fmt.Println(word.S)
		    }
		}
	}
	return "", nil
}
