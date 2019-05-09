package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/ahmetb/go-linq"
)

type upcomingitem struct {
	Category           string  `json:"category"`
	Product            string  `json:"product"`
	RootCollection     string  `json:"root_collection"`
	Name               string  `json:"name"`
	IsReprint          bool    `json:"is_reprint"`
	CSSClass           string  `json:"css_class"`
	Price              float64 `json:"price"`
	ExpectedByOverride string  `json:"expected_by_override"`
	Collection         string  `json:"collection"`
	ProductURL         string  `json:"product_url"`
	OrderIndex         int     `json:"order_index"`
	StatusImageURL     string  `json:"status_image_url"`
	LastUpdated        float64 `json:"last_updated"`
	CollectionCrumbs   string  `json:"collection_crumbs"`
	ProductCode        string  `json:"product_code"`
	ProductImageURL    string  `json:"product_image_url"`
	ExpectedBy         string  `json:"expected_by"`
}

func getUpcomingData(s string) string {
	startingString := "upcoming_data ="
	endingString := "expected_by\": \"\"}]"
	startIndex := strings.Index(s, startingString)
	endIndex := strings.Index(s, endingString)
	return s[startIndex+len(startingString) : endIndex+len(endingString)]
}

func getUpcomingArray(s string) []upcomingitem {
	arrays := []byte(s)
	keys := make([]upcomingitem, 0)
	json.Unmarshal(arrays, &keys)
	return keys
}

func getStarWars(item []upcomingitem) []upcomingitem {
	var owners []upcomingitem
	From(item).WhereT(func(c upcomingitem) bool {
		return strings.Contains(c.RootCollection, "Star Wars: Legion")
	}).SelectT(func(c upcomingitem) upcomingitem {
		return c
	}).ToSlice(&owners)
	return owners
}

func getPrettyStatus(item []upcomingitem) string {
	var result string
	for _, element := range item {
		result += element.Product + " " + element.Name + "\r\n"
	}
	return result
}

func main() {
	resp, err := http.Get("https://www.fantasyflightgames.com/en/upcoming/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	upcomingData := getUpcomingData(string(body))
	upcomingDataArray := getUpcomingArray(upcomingData)
	onlyStarWars := getStarWars(upcomingDataArray)
	fmt.Println(getPrettyStatus(onlyStarWars))
}
