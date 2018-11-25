package posts

type Repo interface {
	All(int, int) ([]*Post, error)
	Find(string) (*Post, error)
	Save(*Post) (*Post, error)
	Delete(*Post) error
}
