package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"

	"blogger/posts"
)

type CreatePostForm struct {
	Title  string `valid:"required~Title is blank"`
	Body   string `valid:"required~Body is blank"`
	Author string `valid:"required~Author is blank"`
	err    error  `valid:"-"`
	Csrf   string `valid:"-"`
}

func (self *CreatePostForm) Error() error {
	return self.err
}

func (self *CreatePostForm) Submit(r *http.Request) {
	csrf, err := r.Cookie("csrf")
	if err != nil {
		self.err = errors.New("CSRF failed!")
		return
	}
	if csrf == nil {
		self.err = errors.New("CSRF failed!")
	}
	if csrf.Value != r.FormValue("csrf") {
		self.err = errors.New("CSRF failed!")
		return
	}
	self.Title = r.FormValue("title")
	self.Body = r.FormValue("body")
	self.Author = r.FormValue("author")
	_, self.err = govalidator.ValidateStruct(self)
}

func (self *Web) Create(w http.ResponseWriter, r *http.Request) {
	ctx := NewViewContext("New-Post")
	form := &CreatePostForm{}
	ctx.Set("Form", form)
	if r.Method == "POST" {
		if form.Submit(r); form.Error() == nil {
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
	csrf, _ := uuid.NewV4()
	form.Csrf = csrf.String()
	http.SetCookie(w, &http.Cookie{
		Name:    "csrf",
		Value:   csrf.String(),
		Expires: time.Now().Add(5 * time.Minute),
	})
	self.ExecuteTemplate(w, "create.html", ctx)
}
