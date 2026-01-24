package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/r2"
)

func main() {
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
	http.HandleFunc("/feed.xml", feedHandler)
	http.HandleFunc("/articles/", articlesHandler)
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

func articlesHandler(w http.ResponseWriter, req *http.Request) {
	bucket, err := r2.NewBucket("STATIC_BUCKET")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract the path after /articles/
	path := strings.TrimPrefix(req.URL.Path, "/articles/")
	if path == "" {
		http.NotFound(w, req)
		return
	}

	// Determine content type and R2 key
	var contentType string
	var r2Key string

	if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
		r2Key = "articles/" + path
	} else {
		// Default to HTML - add .html extension if not present
		contentType = "text/html; charset=utf-8"
		if !strings.HasSuffix(path, ".html") {
			r2Key = "articles/" + path + ".html"
		} else {
			r2Key = "articles/" + path
		}
	}

	obj, err := bucket.Get(r2Key)
	if err != nil || obj == nil {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Content-Type", contentType)
	if contentType == "text/html; charset=utf-8" {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; img-src 'self' https: data: blob:; script-src 'unsafe-inline' https://platform.twitter.com;")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	}

	io.Copy(w, obj.Body)
}
