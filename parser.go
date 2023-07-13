package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

func main() {
	r, err := http.Get("https://www.octranspo.com/en/plan-your-trip/schedules-maps/")
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		os.Exit(1)
	}

	if r.StatusCode != http.StatusOK {
		fmt.Errorf("Did not recieve 200 OK.\n")
		os.Exit(1)
	}

	body := r.Body
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Errorf("Did not recieve an HTML reply.\n")
		os.Exit(1)
	}

	// this section assumes that the webpage contains an element called selectRoute
	ddList := doc.Find("#selectRoute")
	children := ddList.Children()
	var childrenText [][]string
	childrenText = [][]string{
		{"route_number", "route"},
	}

	// parse the children of the selectRoute element
	children.Each(func(i int, child *goquery.Selection) {
		// separate the route number and the route
		data := strings.Split(child.Text(), " ")
		routeName := strings.Join(data[1:], " ")
		// append to selectRoute
		childrenText = append(childrenText, []string{data[0], routeName})
	})
	fmt.Println(childrenText)

	file, _ := os.Create("oc_transpo_routes.csv")
	writer := csv.NewWriter(file)
	err = writer.WriteAll(childrenText)
	if err != nil {
		fmt.Errorf("%v", err)
	}
}
