package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var wg sync.WaitGroup

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

func newsRoutine(c chan News, location string) {
	defer wg.Done()
	var n News
	location = strings.TrimSpace(location) //Remove trailing whitespace at the end. Not having this causes errors
	resp, _ := http.Get(location)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()

	c <- n
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1> Woah! Go is neat! </h1>")
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var s SiteMapIndex
	xml.Unmarshal(bytes, &s)

	newsMap := make(map[string]StoryParams)

	queue := make(chan News, 100)

	for _, location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, location)
	}

	wg.Wait()
	close(queue)
	for elem := range queue {
		for _, story := range elem.Stories {
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
	http.ListenAndServe(":8000", nil)
}
