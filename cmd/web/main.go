package main

import (
	"chat/internal/driver"
	"chat/internal/models"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	api  string
	db   struct {
		dsn string
	}
	secretkey string
}

type application struct {
	config
	infoLog  log.Logger
	errorLog log.Logger
	Messages []string
	DB       models.DBModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}
	app.infoLog.Printf("Starting the server at port : %d", app.config.port)
	return srv.ListenAndServe()
}

func main() {
	var cfg config
	// const (
	// 	host     = "localhost"
	// 	port     = 5432
	// 	user     = "localuser"
	// 	password = "Pwar_1234"
	// 	dbname   = "chatbox"
	// )
	flag.StringVar(&cfg.api, "api", "http://localhost:8000", "URL to api")
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.secretkey, "secret", "dSgVkYp3s6v9y$B&E)H@McQeThWmZq4t", "secret key")
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// flag.StringVar(&cfg.db.dsn, "dsn", psqlInfo, "DSN")
	flag.StringVar(&cfg.db.dsn, "dsn", "localuser:Pwar_1234@tcp(127.0.0.1:3306)/chatbox?parseTime=true&tls=false", "DSN")

	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO \t ", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR \t ", log.Ltime|log.Ldate)

	//DB conn
	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  *infoLog,
		errorLog: *errorLog,
		Messages: make([]string, 0),
		DB:       models.DBModel{DB: conn},
	}

	go app.ListenToWsChannel()

	err = app.serve()
	if err != nil {
		app.errorLog.Printf("Unable to start the server : %d", app.config.port)
	}

}
