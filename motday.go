package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/dweidenfeld/motday/config"
	"github.com/dweidenfeld/motday/flickr"
)

func main() {
	http.HandleFunc("/", photoHandler)
	http.HandleFunc("/favicon.ico", http.NotFoundHandler().ServeHTTP)
	http.ListenAndServe(":8080", nil)
}

func photoHandler(w http.ResponseWriter, r *http.Request) {
	config, err := config.Load("config.json")
	if nil != err {
		panic(err)
	}
	flickr := flickr.New(config.Flickr.APIKey)

	motive := config.RandomMotive()
	query := strings.Join(motive.Queries, ", ")
	//query := *motive.RandomQuery()

	image, err := flickr.SearchRandom(query)
	if nil != err {
		panic(err)
	}

	t, err := template.ParseFiles("photoTemplate.html")
	if nil != err {
		panic(err)
	}

	t.Execute(w, Data{
		Motive: motive,
		Image:  image,
	})
}

// Data Model
type Data struct {
	Motive *config.MotiveConf
	Image  *flickr.Image
}
