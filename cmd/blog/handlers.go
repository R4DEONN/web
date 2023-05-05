package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	Title           string
	SubTitle        string
	FeaturedPosts   []featuredPostData
	MostRecentPosts []mostRecentPostData
}

type featuredPostData struct {
	PostID      string `db:"post_id"`
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	ImgModifier string `db:"image_mod"`
	Author      string `db:"author"`
	AuthorImg   string `db:"author_url"`
	PublishDate string `db:"publish_date"`
}

type mostRecentPostData struct {
	PostID      string `db:"post_id"`
	TopImg      string `db:"image_url"`
	Title       string `db:"title"`
	SubTitle    string `db:"subtitle"`
	AuthorImg   string `db:"author_url"`
	Author      string `db:"author"`
	PublishDate string `db:"publish_date"`
}

type postData struct {
	Title    string `db:"title"`
	SubTitle string `db:"subtitle"`
	Content  string `db:"content"`
	ImageURL string `db:"image_url"`
}

type fullPostData struct {
	Title        string `json:"title"`
	SubTitle     string `json:"subtitle"`
	AuthorName   string `json:"authorName"`
	AuthorAvatar string `json:"authorAvatar"`
	PublishDate  string `json:"publishDate"`
	MainImage    string `json:"mainImage"`
	PreviewImage string `json:"previewImage"`
	Content      string `json:"content"`
}

type userData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func index(client *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("pages/index.html")
		if err != nil {
			http.Error(w, "Inernal Server Error", 500)
			log.Println(err.Error())
			return
		}

		main_post, err := featuredPosts(client)
		if err != nil {
			log.Fatal(err)
		}

		post, err := mostRecentPosts(client)
		if err != nil {
			log.Fatal(err)
		}

		data := indexPage{
			Title:           "Let's do it together.",
			SubTitle:        "We travel the world in search of stories. Come along for the ride.",
			FeaturedPosts:   main_post,
			MostRecentPosts: post,
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		log.Println("Request completed successfully")
	}
}

func post(client *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["postID"]

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusForbidden)
			log.Println(err)
			return
		}

		post, err := postByID(client, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Inernal Server Error", 500)
			log.Println(err.Error())
			return
		}

		err = ts.Execute(w, post)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		log.Println("Request completed successfully")
	}
}

func featuredPosts(client *sqlx.DB) ([]featuredPostData, error) {
	const query = `
		SELECT
		    post_id,
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_mod
		FROM
			post
		WHERE featured = 1
	`

	var posts []featuredPostData

	err := client.Select(&posts, query)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func mostRecentPosts(client *sqlx.DB) ([]mostRecentPostData, error) {
	const query = `
		SELECT
		    post_id,
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url
		FROM
			post
		WHERE featured = 0
	`

	var posts []mostRecentPostData

	err := client.Select(&posts, query)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func postByID(client *sqlx.DB, postID int) (postData, error) {
	const query = `
        SELECT
            title,
            subtitle,
            content,
            image_url
        FROM
            post
        WHERE
            post_id = ?
    `

	var post postData

	err := client.Get(&post, query, postID)
	if err != nil {
		return postData{}, err
	}

	return post, nil
}

func admin(client *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("pages/admin.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}

		var data int

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}
	}
}

func createPost(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const query = `
            INSERT INTO
                post (
                title,
                subtitle,
                author,
                author_url,
                publish_date,
                image_url,
                article_image_url,
                content
            )
            VALUES (
                ?, ?, ?, ?, ?, ?, ?, ?
            );
        `

		decoder := json.NewDecoder(r.Body)
		var post fullPostData
		err := decoder.Decode(&post)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}

		db.MustExec(
			query,
			post.Title,
			post.SubTitle,
			post.AuthorName,
			post.AuthorAvatar,
			post.PublishDate,
			post.PreviewImage,
			post.MainImage,
			post.Content,
		)
	}
}

func login(client *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("pages/auth/logination.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}

		var data int

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}
	}
}

func auth(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const query = `
			SELECT
			  email,
			  password
			FROM
			  user
			WHERE
			  email = ? AND password = ?
		`
		decoder := json.NewDecoder(r.Body)
		var user userData
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}

		_, err = db.Query(query, user.Email, user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
				log.Printf(err.Error())
				return
			}
			http.Error(w, "Internal Server Error", 500)
			log.Printf(err.Error())
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
