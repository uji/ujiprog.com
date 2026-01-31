package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/r2"
	"github.com/uji/ujiprog.com/ogimage"
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
		// Redirect unknown paths to root
		if req.URL.Path != "/" {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}

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

// OGMeta represents OG image metadata for an article
type OGMeta struct {
	Title string `json:"title"`
}

// OGMetaData maps article slug to OG metadata
type OGMetaData map[string]OGMeta

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

	// Handle OG image requests dynamically
	if strings.HasSuffix(path, ".png") {
		handleOGImage(w, req, bucket, path)
		return
	}

	// Default to HTML - add .html extension if not present
	var r2Key string
	if !strings.HasSuffix(path, ".html") {
		r2Key = "articles/" + path + ".html"
	} else {
		r2Key = "articles/" + path
	}

	obj, err := bucket.Get(r2Key)
	if err != nil || obj == nil {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; img-src 'self' https: data: blob:; script-src 'self' https://platform.twitter.com; frame-src https://platform.twitter.com https://syndication.twitter.com;")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	io.Copy(w, obj.Body)
}

// handleOGImage generates OG images dynamically
func handleOGImage(w http.ResponseWriter, req *http.Request, bucket *r2.Bucket, path string) {
	// Extract article slug from path (e.g., "my-article.png" -> "my-article")
	slug := strings.TrimSuffix(path, ".png")

	// Load OG metadata
	ogMetaObj, err := bucket.Get("og-meta.json")
	if err != nil || ogMetaObj == nil {
		http.Error(w, "OG metadata not found", http.StatusInternalServerError)
		return
	}
	ogMetaData, err := io.ReadAll(ogMetaObj.Body)
	if err != nil {
		http.Error(w, "Failed to read OG metadata", http.StatusInternalServerError)
		return
	}

	var ogMeta OGMetaData
	if err := json.Unmarshal(ogMetaData, &ogMeta); err != nil {
		http.Error(w, "Failed to parse OG metadata", http.StatusInternalServerError)
		return
	}

	// Find the title for this article
	meta, ok := ogMeta[slug]
	if !ok {
		http.NotFound(w, req)
		return
	}

	// Load template image
	templateObj, err := bucket.Get("templates/blog-ogp-tmpl.png")
	if err != nil || templateObj == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	templateData, err := io.ReadAll(templateObj.Body)
	if err != nil {
		http.Error(w, "Failed to read template", http.StatusInternalServerError)
		return
	}

	// Load fonts
	asciiFontObj, err := bucket.Get("fonts/DMSans-Bold.ttf")
	if err != nil || asciiFontObj == nil {
		http.Error(w, "ASCII font not found", http.StatusInternalServerError)
		return
	}
	asciiFontData, err := io.ReadAll(asciiFontObj.Body)
	if err != nil {
		http.Error(w, "Failed to read ASCII font", http.StatusInternalServerError)
		return
	}

	japaneseFontObj, err := bucket.Get("fonts/NotoSansJP-Bold.ttf")
	if err != nil || japaneseFontObj == nil {
		http.Error(w, "Japanese font not found", http.StatusInternalServerError)
		return
	}
	japaneseFontData, err := io.ReadAll(japaneseFontObj.Body)
	if err != nil {
		http.Error(w, "Failed to read Japanese font", http.StatusInternalServerError)
		return
	}

	// Create OG image generator
	generator, err := ogimage.NewGenerator(templateData, asciiFontData, japaneseFontData, 56)
	if err != nil {
		http.Error(w, "Failed to create OG generator: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate OG image
	var buf bytes.Buffer
	if err := generator.Generate(meta.Title, &buf); err != nil {
		http.Error(w, "Failed to generate OG image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(buf.Bytes())
}
