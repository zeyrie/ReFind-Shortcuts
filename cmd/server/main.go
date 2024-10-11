package main

import (
	"html/template"
	"net/http"

	"github.com/zeyrie/ReFind-Shortcuts/internal/repo"
)

var tpl *template.Template
var db = &repo.RepoManager{}

func init() {

	db.InitializeTable()

	tpl = template.Must(template.ParseGlob("web/templates/*"))

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	http.ListenAndServe(":4334", nil)

}
