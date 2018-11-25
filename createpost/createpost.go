package createpost

type PostCreatorService struct {
	repo     posts.Repo
	searcher posts.Searcher
}

func (self *PostCreatorService) Execute(title string, body string, author string) (*posts.Post, error) {
	post := &posts.Post{
		Title:  title,
		Body:   body,
		Author: author,
	}
	post, err := self.repo.Save(post)
	if err != nil {
		return nil, err
	}
	err = self.searcher.Index(post)
	if err != nil {
		return post, err
	}
	return post, nil
}

func NewPostCreatorService(repo posts.Repository, searcher posts.Searcher) *PostCreatorService {
	return &PostCreatorService{repo: repo, searcher: searcher}
}
