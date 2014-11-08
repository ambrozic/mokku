package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"mokku/constants"
	"net/http"
)

const (
	dir_templates = "/go/src/pkg/mokku/templates/"
)

var (
	context = Context{}
)

type Web struct {
	Instance
	storage *Storage
}

type Context struct {
	storage *Storage
}

func NewWeb(settings *Settings, storage *Storage) *Web {
	return &Web{Instance: Instance{settings: settings}, storage: storage}
}

func (self *Web) Start() {
	fmt.Println("starting http server")

	context.storage = self.storage

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/json", handleIndexJson)
	HandleError(http.ListenAndServe(fmt.Sprintf(":%d", constants.HTTP_SERVER_PORT), nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := json.MarshalIndent(context.storage, "", "  ")
	HandleError(err)

	templateParser, err := template.ParseFiles(fmt.Sprintf("%s%s", dir_templates, "index.html"))
	HandleError(err)

	w.Header().Set("Content-Type", "text/html")
	templateParser.Execute(w, &struct{ Title, Body string }{Title: "title", Body: string(content)})
}

func handleIndexJson(w http.ResponseWriter, r *http.Request) {
	content, err := json.MarshalIndent(context.storage, "", "  ")
	HandleError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}
