package db

import (
	"blogger/posts"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

var repo *PostRepo
var db *bolt.DB

func init() {
	var err error
	err = os.Remove("test.db")
	if err != nil {
		log.Printf("Error removing test.db: %v", err)
	}
	db, err = bolt.Open("test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(PostsBucket)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
	repo = &PostRepo{db: db}
}

func TestPostRepoCrud(t *testing.T) {
	post := &posts.Post{Title: "My Post"}
	post, err := repo.Save(post)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	found, err := repo.Find(post.ID)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if found.Title != post.Title {
		t.Errorf("Expected to get: %v but got: %v", post, found)
	}
	err = repo.Delete(post)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

func TestPostRepoAll(t *testing.T) {
	postList := make([]*posts.Post, 0)
	for i := 0; i < 20; i++ {
		post, err := repo.Save(&posts.Post{Title: fmt.Sprintf("%d-Post", i)})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		postList = append(postList, post)
	}

	result, err := repo.All(0, 10)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if len(result) != 10 {
		t.Errorf("Expected 10 results but got: %d", len(result))
	}

	result, err = repo.All(15, 30)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 results got: %v", len(result))
	}

	for _, post := range postList {
		err := repo.Delete(post)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
	}
}
