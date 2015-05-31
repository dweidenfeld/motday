package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/dweidenfeld/motday/config"
	"github.com/dweidenfeld/motday/flickr"
)

var (
	dc chan Data
)

func init() {
	dc = make(chan Data, 10)
}

func main() {
	c, err := config.Load("config.json")
	if nil != err {
		panic(err)
	}
	f := flickr.New(c.Flickr.APIKey)
	go fetchNext(c, f)
	http.HandleFunc("/", photoHandler)
	http.HandleFunc("/favicon.ico", http.NotFoundHandler().ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func fetchNext(c *config.Config, f *flickr.Flickr) {
	for {
		motive := c.RandomMotive()
		query := strings.Join(motive.Queries, ", ")

		image, err := f.SearchRandom(query)
		if nil != err {
			panic(err)
		}

		dc <- Data{
			Motive: motive,
			Image:  image,
		}
	}
}

func photoHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("photoTemplate.html")
	if nil != err {
		panic(err)
	}
	t.Execute(w, <-dc)
}

// Data Model
type Data struct {
	Motive *config.MotiveConf
	Image  *flickr.Image
}
