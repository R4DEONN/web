package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	authCookieName = "uid"
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

type createPostRequest struct {
	Title            string `json:"title"`
	SubTitle         string `json:"subtitle"`
	AuthorName       string `json:"authorName"`
	AuthorAvatarName string `json:"authorAvatarName"`
	AuthorAvatar     string `json:"authorAvatar"`
	PublishDate      string `json:"publishDate"`
	MainImageName    string `json:"mainImageName"`
	MainImage        string `json:"mainImage"`
	PreviewImageName string `json:"previewImageName"`
	PreviewImage     string `json:"previewImage"`
	Content          string `json:"content"`
}

type userData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func index(client *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("pages/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		mainPost, err := featuredPosts(client)
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
			FeaturedPosts:   mainPost,
			MostRecentPosts: post,
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		log.Println("Request completed successfully")
		return
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
				http.Error(w, "Post not found", http.StatusNotFound)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		return
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
		err := authByCookie(client, w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ts, err := template.ParseFiles("pages/admin.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		return
	}
}

func createPost(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := authByCookie(db, w, r)
		if err != nil {
			return
		}

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
		var post createPostRequest
		err = decoder.Decode(&post)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		authorImage, err := base64.StdEncoding.DecodeString(post.AuthorAvatar)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
		mainImage, err := base64.StdEncoding.DecodeString(post.MainImage)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
		previewImage, err := base64.StdEncoding.DecodeString(post.PreviewImage)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		post.AuthorAvatarName = "static/img/" + post.AuthorAvatarName
		post.MainImageName = "static/img/" + post.MainImageName
		post.PreviewImageName = "static/img/" + post.PreviewImageName

		authorImageFile, err := os.Create(post.AuthorAvatarName)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
		_, err = authorImageFile.Write(authorImage)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		mainImageFile, err := os.Create(post.MainImageName)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
		_, err = mainImageFile.Write(mainImage)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		previewImageFile, err := os.Create(post.PreviewImageName)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}
		_, err = previewImageFile.Write(previewImage)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		_, err = db.Exec(
			query,
			post.Title,
			post.SubTitle,
			post.AuthorName,
			post.AuthorAvatarName,
			post.PublishDate,
			post.PreviewImageName,
			post.MainImageName,
			post.Content,
		)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func login(client *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const query = `
        	        SELECT
        	            1
        	        FROM
        	            user
        	        WHERE
        	            user_id = ?
        	    `
		cookie, err := r.Cookie(authCookieName)
		if err == nil {
			_, err := client.Exec(query, cookie.Value)
			if err == nil {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
		}

		ts, err := template.ParseFiles("pages/auth/logination.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		return
	}
}

func auth(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user userData
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		const query = `
                    SELECT
                      user_id
                    FROM
                      user
                    WHERE
                      email = ? AND password = ?
                `

		var id int
		err = db.Get(&id, query, user.Email, user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Incorrect password or email", http.StatusUnauthorized)
				log.Printf(err.Error())
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf(err.Error())
			return
		}

		cookie := http.Cookie{
			Name:    authCookieName,
			Value:   strconv.Itoa(id),
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		}

		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func authByCookie(db *sqlx.DB, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println(err)
			return err
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	userIDStr := cookie.Value

	const query = `
        	        SELECT
        	            1
        	        FROM
        	            user
        	        WHERE
        	            user_id = ?
        	    `
	_, err = db.Exec(query, userIDStr)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Incorrect cookie", http.StatusUnauthorized)
			log.Printf(err.Error())
			return err
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf(err.Error())
		return err
	}

	return nil
}

func logOut(w http.ResponseWriter, _ *http.Request) {
	cookie := http.Cookie{
		Name:    authCookieName,
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, -1),
	}

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusAccepted)
	return
}
