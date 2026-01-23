package markdown

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

// TwitterEmbedder handles Twitter/X tweet embedding
type TwitterEmbedder struct {
	client *http.Client
}

// NewTwitterEmbedder creates a new TwitterEmbedder
func NewTwitterEmbedder() *TwitterEmbedder {
	return &TwitterEmbedder{
		client: &http.Client{},
	}
}

// twitterURLPattern matches Twitter/X status URLs
var twitterURLPattern = regexp.MustCompile(
	`^https://(twitter\.com|x\.com)/([^/]+)/status/(\d+)`,
)

// IsTwitterURL checks if the URL is a Twitter/X status URL
func IsTwitterURL(urlStr string) bool {
	return twitterURLPattern.MatchString(urlStr)
}

// OEmbedResponse represents the Twitter oEmbed API response
type OEmbedResponse struct {
	HTML string `json:"html"`
}

// GetEmbedHTML fetches the embed HTML for a tweet using Twitter oEmbed API
func (t *TwitterEmbedder) GetEmbedHTML(tweetURL string) (string, error) {
	oembedURL := fmt.Sprintf(
		"https://publish.twitter.com/oembed?url=%s&omit_script=true",
		url.QueryEscape(tweetURL),
	)

	resp, err := t.client.Get(oembedURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oEmbed API returned status %d", resp.StatusCode)
	}

	var oembed OEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&oembed); err != nil {
		return "", err
	}

	return oembed.HTML, nil
}

// GenerateEmbed generates the full embed HTML including wrapper
func (t *TwitterEmbedder) GenerateEmbed(tweetURL string) (string, error) {
	embedHTML, err := t.GetEmbedHTML(tweetURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`<div class="twitter-embed">%s</div>`, embedHTML), nil
}
