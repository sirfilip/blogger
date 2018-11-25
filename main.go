package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/go-chi/chi"

	"blogger/db"
	"blogger/search"
	"blogger/web"
)

func main() {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	boltDb, err := bolt.Open("dev.db", 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.Println(err)
		return
	}
	err = boltDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(db.PostsBucket)
		return err
	})
	defer boltDb.Close()
	postSearch, err := search.NewPostSearch("posts.bleve")
	if err != nil {
		log.Println(err)
		return
	}
	web := web.New(templates, db.NewPostRepo(boltDb), postSearch)

	router := chi.NewRouter()

	router.HandleFunc("/", web.Index)
	router.HandleFunc("/create", web.Create)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
