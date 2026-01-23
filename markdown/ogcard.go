package markdown

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// OGData contains Open Graph metadata
type OGData struct {
	Title       string
	Description string
	Image       string
	URL         string
	SiteName    string
}

// OGCardFetcher handles fetching OG metadata from URLs
type OGCardFetcher struct {
	client *http.Client
}

// NewOGCardFetcher creates a new OGCardFetcher
func NewOGCardFetcher() *OGCardFetcher {
	return &OGCardFetcher{
		client: &http.Client{},
	}
}

// FetchOGData fetches Open Graph metadata from a URL
func (o *OGCardFetcher) FetchOGData(url string) (*OGData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; OGCardBot/1.0)")

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: status %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	ogData := &OGData{URL: url}
	extractOGData(doc, ogData)

	return ogData, nil
}

// extractOGData traverses the HTML and extracts OG metadata
func extractOGData(n *html.Node, data *OGData) {
	if n.Type == html.ElementNode {
		if n.Data == "meta" {
			var property, content string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "property", "name":
					property = attr.Val
				case "content":
					content = attr.Val
				}
			}
			switch property {
			case "og:title":
				data.Title = content
			case "og:description":
				data.Description = content
			case "og:image":
				data.Image = content
			case "og:site_name":
				data.SiteName = content
			}
		}
		if n.Data == "title" && data.Title == "" && n.FirstChild != nil {
			data.Title = n.FirstChild.Data
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractOGData(c, data)
	}
}

// GenerateOGCard generates HTML for an OG card component
func (o *OGCardFetcher) GenerateOGCard(url string) (string, error) {
	data, err := o.FetchOGData(url)
	if err != nil {
		return "", err
	}

	return RenderOGCard(data), nil
}

// RenderOGCard renders OGData as an HTML card component
func RenderOGCard(data *OGData) string {
	description := data.Description
	if len(description) > 120 {
		description = description[:117] + "..."
	}

	siteName := data.SiteName
	if siteName == "" {
		siteName = extractDomain(data.URL)
	}

	imageHTML := ""
	if data.Image != "" {
		imageHTML = fmt.Sprintf(`<div class="og-card-image">
      <img src="%s" alt="%s" loading="lazy">
    </div>`, escapeHTML(data.Image), escapeHTML(data.Title))
	}

	return fmt.Sprintf(`<a href="%s" class="og-card" target="_blank" rel="noopener noreferrer">
  <div class="og-card-content">
    <div class="og-card-title">%s</div>
    <div class="og-card-description">%s</div>
    <div class="og-card-site">%s</div>
  </div>
  %s
</a>`,
		escapeHTML(data.URL),
		escapeHTML(data.Title),
		escapeHTML(description),
		escapeHTML(siteName),
		imageHTML,
	)
}

// extractDomain extracts the domain from a URL
func extractDomain(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}
	return url
}

// escapeHTML escapes special HTML characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
