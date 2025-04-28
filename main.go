package main

import (
	"log"
	"net/http"
	"shorten/app"
	"shorten/redisConn"
)

func main() {
	rdb := redisConn.RedisDB()
	defer rdb.Close()

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		app.InputURL(w, r, rdb)
	})

	http.HandleFunc("/short", func(w http.ResponseWriter, r *http.Request) {
		app.Redirect(w, r, rdb)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
