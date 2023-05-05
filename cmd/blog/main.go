package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/gorilla/mux"
)

const (
	host         = "localhost:3030"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	client := sqlx.NewDb(db, dbDriverName)

	mux := mux.NewRouter()

	mux.HandleFunc("/home", index(client))

	mux.HandleFunc("/post/{postID}", post(client))

    mux.HandleFunc("/admin", admin(client))
    mux.HandleFunc("/createPost"), createPost(client))

    mux.HandleFunc("/login", login(client))

	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Start server " + host)
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "root:Vuzohe67@tcp(localhost:3306)/blog?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
