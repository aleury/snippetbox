package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"adameury.io/snippetbox/pkg/models/postgres"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP netork address")
	dsn := flag.String("dsn", "postgres://user:@localhost:5432/snippetbox", "Postgres database url")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &postgres.SnippetRepo{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
