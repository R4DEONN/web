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

	router := mux.NewRouter()

	router.HandleFunc("/home", index(client))

	router.HandleFunc("/post/{postID}", post(client))

	router.HandleFunc("/admin", admin(client))
	router.HandleFunc("/api/post", createPost(client)).Methods(http.MethodPost)
	router.HandleFunc("/api/logout", logOut)

	router.HandleFunc("/login", login(client))
	router.HandleFunc("/api/login", auth(client)).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Start server " + host)
	err = http.ListenAndServe(host, router)
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
