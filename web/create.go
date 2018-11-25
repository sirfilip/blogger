package web

import (
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"

	"blogger/posts"
)

type CreatePostForm struct {
	Title  string `valid:"required~Title is blank"`
	Body   string `valid:"required~Body is blank"`
	Author string `valid:"required~Author is blank"`
	Error  error  `valid:"-"`
}

func (self *CreatePostForm) Submit(r *http.Request) error {
	self.Title = r.FormValue("title")
	self.Body = r.FormValue("body")
	self.Author = r.FormValue("author")
	_, self.Error = govalidator.ValidateStruct(self)
	log.Printf("Got error: %v", self.Error)
	return self.Error
}

func (self *Web) Create(w http.ResponseWriter, r *http.Request) {
	ctx := NewViewContext("New-Post")
	form := &CreatePostForm{}
	ctx.Set("Form", form)
	if r.Method == "POST" {
		if err := form.Submit(r); err == nil {
			post, err := self.postRepo.Save(&posts.Post{
				Title:  form.Title,
				Body:   form.Body,
				Author: form.Author,
			})
			if err != nil {
				ResponseError(w, err, http.StatusInternalServerError)
				return
			}
			self.postSearcher.Index(post)
			http.Redirect(w, r, "/", 302)
			return
		}
	}
	self.ExecuteTemplate(w, "create.html", ctx)
}
