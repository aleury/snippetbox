package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"adameury.io/snippetbox/pkg/models/postgres"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *postgres.SnippetRepo
	templateCache map[string]*template.Template
}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) defaultData(r *http.Request) defaultData {
	return defaultData{CurrentYear: time.Now().Year()}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
