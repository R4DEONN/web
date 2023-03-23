package main

import (
	"html/template"
	"log"
	"net/http"
)

type indexPage struct {
	Title           string
	SubTitle        string
	FeaturedPosts   []featuredPostData
	MostRecentPosts []mostRecentPostData
}

type featuredPostData struct {
	Title       string
	Subtitle    string
	ImgModifier string
	Author      string
	AuthorImg   string
	PublishDate string
}

type mostRecentPostData struct {
	TopImg      string
	Title       string
	SubTitle    string
	AuthorImg   string
	Author      string
	PublishDate string
}

type postPage struct {
	Title    string
	SubTitle string
}

func index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/index.html")
	if err != nil {
		http.Error(w, "Inernal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := indexPage{
		Title:           "Let's do it together.",
		SubTitle:        "We travel the world in search of stories. Come along for the ride.",
		FeaturedPosts:   featuredPosts(),
		MostRecentPosts: mostRecentPosts(),
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	log.Println("Request completed successfully")
}

func post(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/post.html")
	if err != nil {
		http.Error(w, "Inernal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := postPage{
		Title:    "The Road Ahead",
		SubTitle: "The road ahead might be paved - it might not be.",
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	log.Println("Request completed successfully")
}

func featuredPosts() []featuredPostData {
	return []featuredPostData{
		{
			Title:       "The Road Ahead",
			Subtitle:    "The road ahead might be paved - it might not be.",
			ImgModifier: "main-block__main-post_image_borealis",
			Author:      "Mat Vogels",
			AuthorImg:   "assets/img/mat_vogels.jpg",
			PublishDate: "September 25, 2015",
		},
		{
			Title:       "From Top Down",
			Subtitle:    "Once a year.",
			ImgModifier: "main-block__main-post_image_lantern",
			Author:      "William Wong",
			AuthorImg:   "assets/img/william_wong.jpg",
			PublishDate: "September 25, 2015",
		},
	}
}

func mostRecentPosts() []mostRecentPostData {
	return []mostRecentPostData{
		{
			TopImg:      "assets/img/still_standing_tall.png",
			Title:       "Still Standing Tall",
			SubTitle:    "Life begins at the end of your comfort zone.",
			Author:      "William Wong",
			AuthorImg:   "assets/img/william_wong.jpg",
			PublishDate: "9/25/2015",
		},
		{
			TopImg:      "assets/img/sunny_side_up.png",
			Title:       "Sunny Side Up",
			SubTitle:    "Life No place is ever as bad as they tell you it’s going to be. at the end of your comfort zone.",
			Author:      "Mat Vogels",
			AuthorImg:   "assets/img/mat_vogels.jpg",
			PublishDate: "9/25/2015",
		},
		{
			TopImg:      "assets/img/water_falls.png",
			Title:       "Water Falls",
			SubTitle:    "We travel not to escape life, but for life not to escape us.",
			Author:      "Mat Vogels",
			AuthorImg:   "assets/img/mat_vogels.jpg",
			PublishDate: "9/25/2015",
		},
		{
			TopImg:      "assets/img/through_the_mist.png",
			Title:       "Through the Mist",
			SubTitle:    "Travel makes you see what a tiny place you occupy in the world.",
			Author:      "William Wong",
			AuthorImg:   "assets/img/william_wong.jpg",
			PublishDate: "9/25/2015",
		},
		{
			TopImg:      "assets/img/awaken_early.png",
			Title:       "Awaken Early",
			SubTitle:    "Not all those who wander are lost.",
			Author:      "Mat Vogels",
			AuthorImg:   "assets/img/mat_vogels.jpg",
			PublishDate: "9/25/2015",
		},
		{
			TopImg:      "assets/img/try_it_always.png",
			Title:       "Try it Always",
			SubTitle:    "The world is a book, and those who do not travel read only one page.",
			Author:      "Mat Vogels",
			AuthorImg:   "assets/img/mat_vogels.jpg",
			PublishDate: "9/25/2015",
		},
	}
}
