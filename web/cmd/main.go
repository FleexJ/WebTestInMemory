package main

import (
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	tokens *mapTokens
}

func main() {
	app := &application{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		tokens: newMapTokens(),
	}

	srv := &http.Server{
		Addr:     ":4000",
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	app.errorLog.Fatal(srv.ListenAndServe())
}
