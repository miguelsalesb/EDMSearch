package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Print("Choose the field you want to search ? \ncreator (1) \ndescription (2) \ncontributor (3) \ntype (4) \nformat (5) \nextent (6) \nlanguage (7) \nissued (8) " +
		"\naggregatedCHO (9) \ndataProvider (10) \nprovider (11) \nisShownAt (12) \nisShownBy (13) \nobject (14) \nrights (15) \n ")
	var field string
	fmt.Scanln(&field)

	fmt.Print("What word do you want to search ? ")
	var info string
	fmt.Scanln(&info)

	fmt.Print("edm (1), dc (2) ou dcterms (3) ?")
	var fieldF string
	fmt.Scanln(&fieldF)

	fmt.Print("xml:lang (1), rdf:about (2), rdf:resource (3) or 'Enter' for other ?")
	var fieldA string
	fmt.Scanln(&fieldA)

	fmt.Print("Start from record number ? ")
	var firstNumber int
	fmt.Scanln(&firstNumber)

	fmt.Print("End in record number ? ")
	var lastNumber int
	fmt.Scanln(&lastNumber)

	f, err := os.Create("results.csv")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	f.Sync()

	var w = bufio.NewWriter(f)

	w.WriteString("record, data")

	for num := firstNumber; num < lastNumber; num++ {
		urlEdm := fmt.Sprintf("%s%d%s", "REPOSITORY URL", num, "&metadataPrefix=edm")
		number := fmt.Sprintf("%v", num)

		getIsShownBy(number, urlEdm, w, field, info, fieldF, fieldA)
	}
}

func getIsShownBy(number string, urlEdm string, w *bufio.Writer, field string, info string, fieldF string, fieldA string) {

	var format, fieldString, fieldAttr, fieldAttribute, searchAttr, fieldC, fieldContent string
	var infoToSearch = strings.ToLower(info)

	res, err := http.Get(urlEdm)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	if field == "1" {
		fieldString = "creator"
	} else if field == "2" {
		fieldString = "description"
	} else if field == "3" {
		fieldString = "contributor"
	} else if field == "4" {
		fieldString = "type"
	} else if field == "5" {
		fieldString = "format"
	} else if field == "6" {
		fieldString = "extent"
	} else if field == "7" {
		fieldString = "language"
	} else if field == "8" {
		fieldString = "issued"
	} else if field == "9" {
		fieldString = "aggregatedCHO"
	} else if field == "10" {
		fieldString = "dataProvider"
	} else if field == "11" {
		fieldString = "provider"
	} else if field == "12" {
		fieldString = "isShownAt"
	} else if field == "13" {
		fieldString = "isShownBy"
	} else if field == "14" {
		fieldString = "object"
	} else if field == "15" {
		fieldString = "rights"
	}

	if fieldF == "1" {
		format = "edm\\:"
	} else if fieldF == "2" {
		format = "dc\\:"
	} else if fieldF == "3" {
		format = "dcterms\\:"
	}

	searchField := fmt.Sprintf("%s", format+fieldString)

	doc.Find(searchField).Each(func(i int, s *goquery.Selection) {
		fieldC = s.Text()
		fieldContent = strings.ToLower(fieldC)

		if fieldA == "1" {
			searchAttr = "xml:lang"
			fieldAttr, _ = s.Attr("xml:lang")
			fieldAttribute = strings.ToLower(fieldAttr)
		} else if fieldA == "2" {
			searchAttr = "rdf:about"
			fieldAttr, _ = s.Attr("rdf:about")
			fieldAttribute = strings.ToLower(fieldAttr)
		} else if fieldA == "3" {
			searchAttr = "rdf:resource"
			fieldAttr, _ = s.Attr("rdf:resource")
			fieldAttribute = strings.ToLower(fieldAttr)
		}

		if fieldAttribute != "" && strings.Contains(fieldAttribute, infoToSearch) {
			w.WriteString("\n" + number)
			w.WriteString("," + fieldAttr)
			fmt.Printf("\nRecord number %v contains \"%v\" in field \"%v\" atributte \"%v\"", number, info, searchField, searchAttr)
		}
		if strings.Contains(fieldContent, infoToSearch) {
			w.WriteString("\n" + number)
			w.WriteString("," + fieldC)
			fmt.Printf("\nRecord number %v contains \"%v\" in field \"%v\"", number, info, searchField)
		}
		w.Flush()
	})
}
