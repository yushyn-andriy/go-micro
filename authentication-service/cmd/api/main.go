package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/yushyn-andriy/authentication/data"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service...")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to db...")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Listening on %s\n", webPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		log.Println("DSN:", dsn)
		connection, err := openDB(dsn)
		if err != nil {

			counts += 1
		} else {
			log.Println("Connected to Postgresql!")
			return connection
		}

		if counts >= 10 {
			return nil
		}
		log.Println("DSN:", dsn)
		log.Println("Backing off for two seconds...")
		time.Sleep(time.Second * 2)
	}
}
