package main

import (
	"blogger/db"
	"blogger/search"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func main() {

	boltDb, err := bolt.Open("dev.db", 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Println(err)
		return
	}
	defer boltDb.Close()
	postSearch, err := search.NewPostSearch("posts.bleve")
	repo := db.NewPostRepo(boltDb)

	posts, _ := repo.All(0, 10001)
	for _, post := range posts {
		postSearch.Index(post)
	}
}
