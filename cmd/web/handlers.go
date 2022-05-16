package main

import (
	"fmt"
	"golangify.com/snippetbox/config"
	"html/template"
	"net/http"
	"strconv"
)

// Обработчик главной страницы.
func home(app *config.Application) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		if ts == nil {
			app.InfoLog.Println("Ошибка рендеринга home.page")
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}

// Обработчик для отображения содержимого заметки.
func showSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.ErrLog.Printf("error while getting id")
			http.NotFound(w, r)
			return
		}
		_, err = fmt.Fprintf(w, "Отображение заметки № %d", id)
		if err != nil {
			app.InfoLog.Println("Ошибка рендеринга snippet")
		}
	}
}

// Обработчик для создания новой заметки.
func createSnippet(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			app.InfoLog.Println("ET-method запрещен")
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "GET-method запрещен", 405)
			return
		}
		_, err := w.Write([]byte("Форма для создания новой заметки..."))
		if err != nil {
			app.InfoLog.Println("Ошибка рендеринга snippet/create")
		}
	}
}