package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SiteMapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Stories []Story `xml:"url"`
}

type Story struct {
	Title    string `xml:"news>title"`
	Keywords string `xml:"news>keywords"`
	Location string `xml:"loc"`
}

type StoryParams struct {
	Keywords string
	Location string
}

// func (l Location) String() string {
// 	return fmt.Sprintf(l.Loc)
// }

func main() {
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var s SiteMapIndex
	xml.Unmarshal(bytes, &s)

	var n News

	newsMap := make(map[string]StoryParams)

	for _, location := range s.Locations {
		location = strings.TrimSpace(location) //Remove trailing whitespace at the end. Not having this causes errors
		resp, _ := http.Get(location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(bytes, &n)

		for _, story := range n.Stories {
			newsMap[story.Title] = StoryParams{story.Keywords, story.Location}
		}
	}

	for title, data := range newsMap {
		fmt.Printf(title, "\n\n")
		fmt.Printf(data.Location, "\n")
		fmt.Printf(data.Keywords, "\n\n ------------------- \n\n")
	}
}
