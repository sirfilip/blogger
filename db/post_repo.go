package db

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"

	"blogger/posts"
)

var (
	PostsBucket = []byte("posts")
)

type PostRepo struct {
	db *bolt.DB
}

func (self *PostRepo) All(offset int, limit int) ([]*posts.Post, error) {
	result := make([]*posts.Post, 0)
	err := self.db.View(func(tx *bolt.Tx) error {
		var err error
		bucket := tx.Bucket(PostsBucket)

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if offset > 0 {
				offset = offset - 1
				continue
			}
			if limit > 0 {
				limit = limit - 1
			} else {
				break
			}
			post := &posts.Post{}
			err = json.Unmarshal(v, post)
			if err != nil {
				return err
			}
			result = append(result, post)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (self *PostRepo) Find(id string) (*posts.Post, error) {
	post := &posts.Post{}
	err := self.db.View(func(tx *bolt.Tx) error {
		var err error
		bucket := tx.Bucket(PostsBucket)

		val := bucket.Get([]byte(id))

		if val == nil {
			return errors.New("Post not found")
		}

		err = json.Unmarshal(val, post)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (self *PostRepo) Save(post *posts.Post) (*posts.Post, error) {
	if post.ID == "" {
		uuid, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		post.ID = uuid.String()
	}

	err := self.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PostsBucket)
		postJSON, err := json.Marshal(post)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(post.ID), postJSON)

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (self *PostRepo) Delete(post *posts.Post) error {
	err := self.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PostsBucket)

		err := bucket.Delete([]byte(post.ID))

		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func NewPostRepo(db *bolt.DB) posts.Repo {
	return &PostRepo{db: db}
}
