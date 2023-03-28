//BIOAFF/backend/cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"time"
)

// version umber 1
const version = "0.1.0"

// configuration settings
type config struct {
	port int
	env  string // development | staging | production
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64 //requests per second
		burst   int     //how many request an the intial momment
		enabled bool    //rate limiting toggle
	}
	smtp struct {
		host     string
		port     int
		username string //from email service probably mail trap
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

// dependency injection
type application struct {
	config config
	logger *jsonlog.logger
	models data.models
	mailer mailer.mailer
	wg     sync.waitGroup
}

func main() {
	var cfg config

	//flags for webserver
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development | staging | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv(), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idel connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	//flags for the rate limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Rate limiter enables")

	//flags for the mailer
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smpt.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", " ", "SMTP username")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", " ", "SMPT sender")

	//cors' flag
	flag.Func("cors-trusted-origin", "Trusted CORS origin (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	//creating the logger instance
	logger := jsonlog.New(os.stdout, jsonlog.Levelinfo)

	//create the connecction pool
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	//loging the successful connection
	logger.PrintInfo("database connection pool established", nil)

	//instance of app struct
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	//call app server() to start the server
	err = app.server()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

// OpenDB() function returns a *sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	//create a context with a 5-second timeout dealine
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
