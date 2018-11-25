package search

import (
	"blogger/posts"
	"log"

	"github.com/blevesearch/bleve"
)

type PostSearchFields struct {
	Title string
	Body  string
}

type PostSearch struct {
	index bleve.Index
}

func (self *PostSearch) Index(post *posts.Post) error {
	return self.index.Index(post.ID, PostSearchFields{
		Title: post.Title,
		Body:  post.Body,
	})
}

func (self *PostSearch) Delete(post *posts.Post) error {
	return self.index.Delete(post.ID)
}

func (self *PostSearch) Search(q string, offset int, limit int) ([]string, error) {
	result := make([]string, 0)
	query := bleve.NewMatchQuery(q)
	search := bleve.NewSearchRequest(query)
	searchResult, err := self.index.Search(search)
	log.Printf("Search result: %#v", searchResult)
	if err != nil {
		return nil, err
	}
	for _, hit := range searchResult.Hits {
		if offset > 0 {
			offset = offset - 1
			continue
		}

		if limit > 0 {
			limit = limit - 1
			result = append(result, hit.ID)
		} else {
			break
		}
	}
	log.Printf("Result ids: %#v", result)
	return result, nil
}

func NewPostSearch(path string) (*PostSearch, error) {
	index, err := bleve.Open(path)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(path, mapping)
		if err != nil {
			return nil, err
		}
	}
	return &PostSearch{index}, nil
}
