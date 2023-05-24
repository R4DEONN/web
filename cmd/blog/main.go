package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	host         = "localhost:3030"
	dbDriverName = "mysql"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}
}

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
	mux.HandleFunc("/api/post", createPost(client)).Methods(http.MethodPost)
	mux.HandleFunc("/api/logout", logOut)

	mux.HandleFunc("/login", login(client))
	mux.HandleFunc("/api/login", auth(client)).Methods("POST")

	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Start server " + host)
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	appDatabaseDSN, exists := os.LookupEnv("APP_DATABASE_DSN")
	if !exists {
		return nil, fmt.Errorf("APP_DATABASE_DSN environment variable not found")
	}
	return sql.Open(dbDriverName, appDatabaseDSN)
}
