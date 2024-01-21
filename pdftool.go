package pdftool

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/signintech/gopdf"
)

func F_writeTextOnPDF(pdf *gopdf.GoPdf, X float64, Y float64, fontName string, fontSize int, text string) {
	err := pdf.SetFont(fontName, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.SetXY(X, Y)
	pdf.Cell(nil, text)
}

func F_writeTextOnPDF_Delimiter(pdf *gopdf.GoPdf, X float64, Y float64, fontName string, fontSize int, text, delimiter string) (float64, float64) {
	var decalageY float64
	decalageY = 0
	multiLineText := strings.Split(text, delimiter)
	for _, Line := range multiLineText {
		F_writeTextOnPDF(pdf, X, Y+decalageY, fontName, fontSize, Line)
		decalageY = decalageY + float64(fontSize)
	}
	return X, Y + decalageY
}

func F_writeTextOnPDF_Underline(pdf *gopdf.GoPdf, X float64, Y float64, fontName string, fontSize int, text string) {
	F_writeTextOnPDF(pdf, X, Y, fontName, fontSize, text)
	pdf.SetLineWidth(1)
	w01, _ := pdf.MeasureTextWidth(text)
	pdf.Line(X, Y+float64(fontSize), X+float64(w01), Y+float64(fontSize))
}

func F_writeTextOnPDF_RightAlign(pdf *gopdf.GoPdf, X float64, Y float64, fontName string, fontSize int, text string) {
	err := pdf.SetFont(fontName, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.SetXY(X, Y)
	pdf.CellWithOption(&gopdf.Rect{
		W: 80,
		H: 10,
	}, text, gopdf.CellOption{Align: gopdf.Right})
}
func F_writeTextOnPDF_RightAlign2(pdf *gopdf.GoPdf, X float64, Y float64, Width float64, fontName string, fontSize int, text string) {
	err := pdf.SetFont(fontName, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.SetXY(X, Y)
	pdf.CellWithOption(&gopdf.Rect{
		W: Width,
		H: 10,
	}, text, gopdf.CellOption{Align: gopdf.Right})
}

func F_writeTextOnPDF_WithBorders(pdf *gopdf.GoPdf, X float64, Y float64, fontName string, fontSize int, text string, padding float64) {

	F_writeTextOnPDF(pdf, X, Y, fontName, fontSize, text)
	w01, _ := pdf.MeasureTextWidth(text)
	pdf.RectFromUpperLeft(X-(padding/2), Y-(padding/4), float64(w01)+(padding), float64(fontSize)+(padding/4))
}

func F_addFontOnPDF(pdf *gopdf.GoPdf, FontSrc, FontName string) {
	err := pdf.AddTTFFont(FontName, FontSrc)
	if err != nil {
		log.Print(err.Error())
		return
	}
}

func F_ajouterPage(pdf *gopdf.GoPdf, compteurPage *int) {
	pdf.AddPage()
	*compteurPage += 1
}

func readPdf(path string) (string, error) {
	file, r, err := pdf.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	buf.ReadFrom(b)

	return buf.String(), nil
}

// Cette fonction m'intéresse, décryptage
func F_readPdf2(path string) (string, error) {

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
