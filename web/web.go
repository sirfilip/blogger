package web

import (
	"blogger/posts"
	"html/template"
	"log"
	"net/http"
)

type Web struct {
	*template.Template
	postRepo     posts.Repo
	postSearcher posts.Searcher
}

type ViewContext struct {
	Title string
	Data  map[string]interface{}
}

func (self *ViewContext) Set(prop string, value interface{}) {
	self.Data[prop] = value
}

func NewViewContext(title string) *ViewContext {
	return &ViewContext{Title: title, Data: make(map[string]interface{})}
}

func New(t *template.Template, postRepo posts.Repo, postSearcher posts.Searcher) *Web {
	return &Web{Template: t, postRepo: postRepo, postSearcher: postSearcher}
}

func ResponseError(w http.ResponseWriter, err error, status int) {
	log.Printf("Got error: %v", err)
	http.Error(w, "Error", status)
}
