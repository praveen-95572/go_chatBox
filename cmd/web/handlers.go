package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// var templateFS embed.FS

type TemplateData struct {
	Data  []string
	Title string
	API   string
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &TemplateData{API: app.config.api}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
	if err := app.renderTemplate(w, r, "one-user", &TemplateData{API: app.config.api}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) personA(w http.ResponseWriter, r *http.Request) {
	td := &TemplateData{
		Data:  app.Messages,
		Title: "PersonA",
	}
	if err := app.renderTemplate(w, r, "personA", td); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, title string, data *TemplateData) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("%s.partials.gohtml", title)
	//t, err = template.New("baselayout.page.gohtml").ParseFS(templateFS, "templates/baselayout.page.gohtml", templateToRender)
	t, err = template.ParseFiles("templates/baselayout.page.gohtml", fmt.Sprintf("templates/%s", templateToRender))
	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}
