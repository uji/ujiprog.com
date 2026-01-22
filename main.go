package main

import (
	"io"
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/r2"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})
	http.HandleFunc("/articles.json", func(w http.ResponseWriter, req *http.Request) {
		bucket, err := r2.NewBucket("STATIC_BUCKET")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		obj, err := bucket.Get("articles.json")
		if err != nil || obj == nil {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.Copy(w, obj.Body)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		bucket, err := r2.NewBucket("STATIC_BUCKET")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		obj, err := bucket.Get("index.html")
		if err != nil || obj == nil {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; img-src 'self' blob:; script-src 'unsafe-inline';")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		io.Copy(w, obj.Body)
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
