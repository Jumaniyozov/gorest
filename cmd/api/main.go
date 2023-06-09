package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/jumaniyozov/gorest/models"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development | production)")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres:123456qwe@localhost:5433/govies?sslmode=disable", "postgres connection string")
	flag.StringVar(&cfg.jwt.secret, "jwtsecret", "AIWgBPK+YCPQAFzRUMbIkyd/hQkZC2hYwkQf4SyUPc5wQNFbl3ifwLvd2iErJ8L9yAtPZMh+zQ5mx4/AI5Ao8A==", "JWT secret")
	flag.Parse()

	logger := log.New(os.Stdout, "-> ", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		config: cfg,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	app.logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
