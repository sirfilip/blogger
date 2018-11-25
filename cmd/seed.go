package main

import (
	"log"
	"time"

	"blogger/db"
	"blogger/posts"

	"github.com/boltdb/bolt"
	"github.com/icrowley/fake"
)

func main() {
	boltDb, err := bolt.Open("dev.db", 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Println(err)
		return
	}
	repo := db.NewPostRepo(boltDb)

	for i := 0; i < 1000; i++ {
		post := &posts.Post{
			Title:  fake.Title(),
			Body:   fake.ParagraphsN(10),
			Author: fake.FullName(),
		}
		_, err := repo.Save(post)
		if err != nil {
			panic(err)
		}
	}
}
