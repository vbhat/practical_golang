package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
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

type NewsAggPage struct {
	Title string
	News  map[string]StoryParams
}

// func (l Location) String() string {
// 	return fmt.Sprintf(l.Loc)
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1> Woah! Go is neat! </h1>")
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
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
	p := NewsAggPage{Title: "Some amazing news!", News: newsMap}
	t, _ := template.ParseFiles("basictemplating.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg", newsAggHandler)
	http.ListenAndServe(":8080", nil)
}
