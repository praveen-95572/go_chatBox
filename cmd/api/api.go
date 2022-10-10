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
	db struct {
		dsn string
	}
	port int
}

type application struct {
	config
	infoLog  log.Logger
	errorLog log.Logger
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

	infoLog := log.New(os.Stdout, "INFO \t ", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR \t ", log.Ltime|log.Ldate)
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "Server port to listen on")
	flag.StringVar(&cfg.db.dsn, "dsn", "localuser:Pwar_1234@tcp(127.0.0.1:3306)/chatbox?parseTime=true&tls=false", "DSN")
	flag.Parse()

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		config:   cfg,
		infoLog:  *infoLog,
		errorLog: *errorLog,
		DB:       models.DBModel{DB: conn},
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Printf("Unable to start the server : %d", app.config.port)
	}

}
