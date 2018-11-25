package web

import (
	"blogger/posts"
	"log"
	"net/http"
)

func (self *Web) Index(w http.ResponseWriter, r *http.Request) {
	var postsFound []*posts.Post
	var err error
	q := r.URL.Query().Get("q")
	ctx := NewViewContext("Home")
	ctx.Set("Q", q)
	if q == "" {
		postsFound, err = self.postRepo.All(0, 10)
		if err != nil {
			ResponseError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		postIDs, err := self.postSearcher.Search(q, 0, 10)
		if err != nil {
			ResponseError(w, err, http.StatusInternalServerError)
			return
		}
		for _, postID := range postIDs {
			post, err := self.postRepo.Find(postID)
			if err != nil {
				ResponseError(w, err, http.StatusInternalServerError)
				return
			}
			postsFound = append(postsFound, post)
		}
	}
	ctx.Set("Posts", postsFound)
	if err = self.ExecuteTemplate(w, "index.html", ctx); err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}
