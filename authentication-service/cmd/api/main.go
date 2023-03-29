package main

import (
	"authentication-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the authentication service")

	//	connect to db

	conn, err := connectToDB()
	if err != nil {
		log.Println("Cannot connect to DB")
	}

	// set up routes
	app := &Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("Starting the server on port", webPort)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting the server", err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to DB", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to DB")
		return connection, nil

	}
}
