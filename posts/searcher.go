package posts

type Searcher interface {
	Search(query string, offset int, limit int) ([]string, error)
	Index(*Post) error
	Delete(*Post) error
}
