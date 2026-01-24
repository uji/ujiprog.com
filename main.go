package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/r2"
)

func main() {
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(`User-agent: *
Allow: /
Allow: /articles/
Allow: /feed.xml
Allow: /avator.jpg

Sitemap: https://ujiprog.com/sitemap.xml`))
	})
	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://ujiprog.com/</loc>
    <lastmod>2026-01-24</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://ujiprog.com/feed.xml</loc>
    <lastmod>2026-01-24</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>`))
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
	http.HandleFunc("/style.css", func(w http.ResponseWriter, req *http.Request) {
		bucket, err := r2.NewBucket("STATIC_BUCKET")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		obj, err := bucket.Get("style.css")
		if err != nil || obj == nil {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		io.Copy(w, obj.Body)
	})
	http.HandleFunc("/article.css", func(w http.ResponseWriter, req *http.Request) {
		bucket, err := r2.NewBucket("STATIC_BUCKET")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		obj, err := bucket.Get("article.css")
		if err != nil || obj == nil {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		io.Copy(w, obj.Body)
	})
	http.HandleFunc("/main.js", func(w http.ResponseWriter, req *http.Request) {
		bucket, err := r2.NewBucket("STATIC_BUCKET")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		obj, err := bucket.Get("main.js")
		if err != nil || obj == nil {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		io.Copy(w, obj.Body)
	})
	http.HandleFunc("/article.js", func(w http.ResponseWriter, req *http.Request) {
		bucket, err := r2.NewBucket("STATIC_BUCKET")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		obj, err := bucket.Get("article.js")
		if err != nil || obj == nil {
			http.NotFound(w, req)
			return
		}

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
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
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' blob: data:; script-src 'self'; object-src 'none'; base-uri 'self'; frame-src https://platform.twitter.com https://syndication.twitter.com; frame-ancestors 'none';")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
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
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; img-src 'self' https: data: blob:; script-src 'self' https://platform.twitter.com; frame-src https://platform.twitter.com https://syndication.twitter.com;")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	}

	io.Copy(w, obj.Body)
}
