package views

import (
	"bytes"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)
import "io/fs"

type Template struct {
	htmlTemplate *template.Template
	logger       *slog.Logger
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	tpl := t.htmlTemplate
	tpl = tpl.Funcs(template.FuncMap{})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		t.logger.Error("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func Must(t Template, err error) Template {
	if err != nil {
		// This should be only called during startup,
		//because if something is wrong application should start
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, logger *slog.Logger, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, err
	}
	return Template{htmlTemplate: tpl}, nil
}
