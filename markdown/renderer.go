package markdown

import (
	"bytes"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
)

// TemplateData represents the data passed to the article template
type TemplateData struct {
	Title           string
	PublishedAt     string
	Content         template.HTML
	OGImageURL      string
	ArticleURL      string
	TwitterShareURL template.URL
	HatenaShareURL  template.URL
}

// Renderer handles HTML template rendering for articles
type Renderer struct {
	tmpl *template.Template
}

// NewRenderer creates a new renderer with the specified template file
func NewRenderer(templatePath string) (*Renderer, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return &Renderer{tmpl: tmpl}, nil
}

// Render renders a ParsedArticle to HTML using the template
func (r *Renderer) Render(article *ParsedArticle) (string, error) {
	articleURL := "https://ujiprog.com/articles/" + article.Filename
	encodedURL := url.QueryEscape(articleURL)
	encodedTitle := url.QueryEscape(article.Meta.Title)

	twitterShareURL := "https://twitter.com/intent/tweet?url=" + encodedURL + "&text=" + encodedTitle
	hatenaShareURL := "https://b.hatena.ne.jp/add?mode=confirm&url=" + encodedURL + "&title=" + encodedTitle

	data := TemplateData{
		Title:           article.Meta.Title,
		PublishedAt:     article.Meta.PublishedAt.Format("2006-01-02"),
		Content:         template.HTML(article.Content),
		OGImageURL:      "https://ujiprog.com/articles/" + article.Filename + ".png",
		ArticleURL:      articleURL,
		TwitterShareURL: template.URL(twitterShareURL),
		HatenaShareURL:  template.URL(hatenaShareURL),
	}

	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// RenderToFile renders a ParsedArticle and writes it to the specified output directory
func (r *Renderer) RenderToFile(article *ParsedArticle, outputDir string) error {
	html, err := r.Render(article)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, article.Filename+".html")
	return os.WriteFile(outputPath, []byte(html), 0644)
}
