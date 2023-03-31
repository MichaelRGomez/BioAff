// Filename: BIOAFF/backend/cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"se_group1.net/internal/jsonlog"
)

// version number
const version = "0.1"

// configuration settings
type config struct {
	port int
	env  string //development, staging, production
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdelTime  string
	}
	limiter struct {
		rps     float64 //request per second
		burst   int     //how many request at the initial momment
		enabled bool    //rate limiting toggle
	}
	/*
		smtp struct {
			host string
			port int
			username string //from mailing system - mail trap username
			password string
			sender string
		}*/

	cors struct {
		trustedOrigins []string
	}
}

// dependency injection
type application struct {
	config config
	logger *jsonlog.Logger
	/*	models data.Models
		mailer mailer.Mailer */
	wg sync.WaitGroup
}

func main() {
	var cfg config

	//db and webserver flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("BIOAFF_DB_DSN"), "PostgresSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQl max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdelTime, "db-max-idle-time", "15m", "PostgreSQl max connection idle time")

	//rate limiter flags
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Rate limiter enabled")

	//using the flag function to parse for our trusted origins
	flag.Func("cors-trusted-origin", "Trusted CORS origin (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	//parsing all the flags we are using
	flag.Parse()

	//creating the logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	//creating the db connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	//ensuring the db connection will be closed
	defer db.Close()

	//logging the successful connection pool
	logger.PrintInfo("database connection pool established", nil)

	//creating an instance of the app struct
	app := &application{
		config: cfg,
		logger: logger,
		//models: data.NewModels(db),
		//no mailer for now
	}

	//calling app.server() to start the server
	err = app.serve()
	if err != nil {
		logger.PrintError(err, nil)
	}
}

// OpenDB() - returns a *sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdelTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	//creating context with a 5 sec timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
