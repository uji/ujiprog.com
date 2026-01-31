package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/uji/ujiprog.com/markdown"
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

// OGMeta represents OG image metadata for an article
type OGMeta struct {
	Title string `json:"title"` // OG画像用タイトル（改行含む）
}

// OGMetaData maps article slug to OG metadata
type OGMetaData map[string]OGMeta

func main() {
	articlesDir := flag.String("articles", "articles", "Directory containing markdown articles")
	outputDir := flag.String("output", "build/articles", "Directory to output generated HTML and images")
	templatePath := flag.String("template", "templates/article.html", "Path to article HTML template")
	articlesJSONPath := flag.String("articles-json", "public/articles.json", "Path to articles.json for merging")
	ogMetaPath := flag.String("og-meta", "", "Path to output og-meta.json (optional)")
	flag.Parse()

	// Ensure output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Find all markdown files
	mdFiles, err := filepath.Glob(filepath.Join(*articlesDir, "*.md"))
	if err != nil {
		log.Fatalf("Failed to find markdown files: %v", err)
	}

	if len(mdFiles) == 0 {
		log.Println("No markdown files found")
		return
	}

	// Create parser and renderer
	parser := markdown.NewParser()
	renderer, err := markdown.NewRenderer(*templatePath)
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}

	// Process each markdown file
	var localArticles []Article
	ogMetaData := make(OGMetaData)
	for _, mdFile := range mdFiles {
		log.Printf("Processing: %s", mdFile)

		// Parse markdown with URL expansion
		article, err := parser.ParseFileWithExpansion(mdFile)
		if err != nil {
			log.Printf("Failed to parse %s: %v", mdFile, err)
			continue
		}

		// Render HTML
		if err := renderer.RenderToFile(article, *outputDir); err != nil {
			log.Printf("Failed to render %s: %v", mdFile, err)
			continue
		}
		log.Printf("Generated: %s/%s.html", *outputDir, article.Filename)

		// Collect OG metadata
		ogMetaData[article.Filename] = OGMeta{
			Title: article.Meta.OGTitle(),
		}

		// Add to local articles list
		localArticles = append(localArticles, Article{
			Title:       article.Meta.Title,
			URL:         "/articles/" + article.Filename,
			PublishedAt: article.Meta.PublishedAt.Format(time.RFC3339),
			Platform:    "blog",
		})
	}

	// Merge with existing articles.json
	if err := mergeArticlesJSON(*articlesJSONPath, localArticles); err != nil {
		log.Printf("Warning: Failed to merge articles.json: %v", err)
	} else {
		log.Printf("Updated: %s", *articlesJSONPath)
	}

	// Save OG metadata if path is specified
	if *ogMetaPath != "" {
		if err := saveOGMeta(*ogMetaPath, ogMetaData); err != nil {
			log.Printf("Warning: Failed to save og-meta.json: %v", err)
		} else {
			log.Printf("Generated: %s", *ogMetaPath)
		}
	}

	log.Println("Generation complete!")
}

// saveOGMeta saves OG metadata to a JSON file
func saveOGMeta(path string, data OGMetaData) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal OG metadata: %w", err)
	}

	if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write og-meta.json: %w", err)
	}

	return nil
}

// mergeArticlesJSON merges local articles with existing articles.json
func mergeArticlesJSON(path string, localArticles []Article) error {
	var existingData ArticlesData

	// Read existing articles.json if it exists
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &existingData); err != nil {
			return fmt.Errorf("failed to parse existing articles.json: %w", err)
		}
	}

	// Create a map of existing articles (excluding blog platform to allow updates)
	articleMap := make(map[string]Article)
	for _, a := range existingData.Articles {
		if a.Platform != "blog" {
			articleMap[a.URL] = a
		}
	}

	// Add local articles (blog platform)
	for _, a := range localArticles {
		articleMap[a.URL] = a
	}

	// Convert map back to slice
	var allArticles []Article
	for _, a := range articleMap {
		allArticles = append(allArticles, a)
	}

	// Sort by published_at descending
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].PublishedAt > allArticles[j].PublishedAt
	})

	// Write updated articles.json
	newData := ArticlesData{Articles: allArticles}
	jsonBytes, err := json.MarshalIndent(newData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal articles: %w", err)
	}

	if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write articles.json: %w", err)
	}

	return nil
}
