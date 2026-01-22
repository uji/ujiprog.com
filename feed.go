package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/syumai/workers/cloudflare/r2"
)

type ArticlesData struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishedAt string `json:"published_at"`
	Platform    string `json:"platform"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	AtomNS  string   `xml:"xmlns:atom,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	AtomLink      AtomLink `xml:"atom:link"`
	Items         []Item   `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        GUID   `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

type GUID struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

func feedHandler(w http.ResponseWriter, req *http.Request) {
	bucket, err := r2.NewBucket("STATIC_BUCKET")
	if err != nil {
		http.Error(w, "bucket error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := bucket.Get("articles.json")
	if err != nil {
		http.Error(w, "get error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if obj == nil {
		http.Error(w, "articles.json not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(obj.Body)
	if err != nil {
		http.Error(w, "read error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var data ArticlesData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "json error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]Item, 0, len(data.Articles))
	for _, article := range data.Articles {
		pubDate := article.PublishedAt
		if t, err := time.Parse(time.RFC3339, article.PublishedAt); err == nil {
			pubDate = t.Format(time.RFC1123Z)
		}
		items = append(items, Item{
			Title: article.Title,
			Link:  article.URL,
			GUID: GUID{
				Value:       article.URL,
				IsPermaLink: "true",
			},
			PubDate:     pubDate,
			Description: article.Title,
		})
	}

	rss := RSS{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: Channel{
			Title:         "ujiprog.com",
			Link:          "https://ujiprog.com/",
			Description:   "uji のブログ",
			Language:      "ja",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			AtomLink: AtomLink{
				Href: "https://ujiprog.com/feed.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Items: items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(xml.Header))
	xml.NewEncoder(w).Encode(rss)
}
