package markdown

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// ArticleMeta represents the YAML frontmatter metadata
type ArticleMeta struct {
	Title       string
	PublishedAt time.Time
}

// ParsedArticle represents a parsed markdown article
type ParsedArticle struct {
	Meta     ArticleMeta
	Content  string // HTML content
	Filename string // filename without extension (e.g., "2026-01-22_created-my-own-blog")
}

// Parser handles markdown parsing with goldmark
type Parser struct {
	md       goldmark.Markdown
	expander *Expander
}

// NewParser creates a new markdown parser with goldmark configuration
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	return &Parser{
		md:       md,
		expander: NewExpander(),
	}
}

// ParseFile parses a markdown file and returns a ParsedArticle
func (p *Parser) ParseFile(path string) (*ParsedArticle, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return p.Parse(content, filenameWithoutExt(path))
}

// Parse parses markdown content and returns a ParsedArticle
func (p *Parser) Parse(source []byte, filename string) (*ParsedArticle, error) {
	var buf bytes.Buffer
	context := parser.NewContext()

	if err := p.md.Convert(source, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	metaData := meta.Get(context)
	articleMeta := extractMeta(metaData)

	return &ParsedArticle{
		Meta:     articleMeta,
		Content:  buf.String(),
		Filename: filename,
	}, nil
}

// ParseWithExpansion parses markdown content with URL expansion
// This expands GitHub URLs, Twitter embeds, and OG cards before parsing
func (p *Parser) ParseWithExpansion(source []byte, filename string) (*ParsedArticle, error) {
	// Expand special URLs
	expanded := p.expander.ExpandContent(string(source))
	// Process custom image widths
	expanded = ProcessImageWidths(expanded)

	return p.Parse([]byte(expanded), filename)
}

// ParseFileWithExpansion parses a markdown file with URL expansion
func (p *Parser) ParseFileWithExpansion(path string) (*ParsedArticle, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return p.ParseWithExpansion(content, filenameWithoutExt(path))
}

// extractMeta extracts ArticleMeta from the frontmatter map
func extractMeta(metaData map[string]interface{}) ArticleMeta {
	am := ArticleMeta{}

	if title, ok := metaData["title"].(string); ok {
		am.Title = title
	}

	if publishedAt, ok := metaData["published_at"].(string); ok {
		if t, err := time.Parse("2006-01-02", publishedAt); err == nil {
			am.PublishedAt = t
		}
	}

	return am
}

// filenameWithoutExt returns the filename without the directory and extension
func filenameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
