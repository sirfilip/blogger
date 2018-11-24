package posts

import "time"

type Status int

const (
	Draft Status = iota
	Published
)

type Post struct {
	ID          string
	Title       string
	Body        string
	Author      string
	Status      Status
	PublishedAt time.Time
}
