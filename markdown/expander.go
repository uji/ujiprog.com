package markdown

import (
	"regexp"
	"strings"
)

// Expander handles expansion of special URLs in markdown content
type Expander struct {
	github  *GitHubCodeExpander
	ogcard  *OGCardFetcher
	twitter *TwitterEmbedder
}

// NewExpander creates a new Expander with all sub-expanders
func NewExpander() *Expander {
	return &Expander{
		github:  NewGitHubCodeExpander(),
		ogcard:  NewOGCardFetcher(),
		twitter: NewTwitterEmbedder(),
	}
}

// urlOnlyLinePattern matches lines that contain only a URL
var urlOnlyLinePattern = regexp.MustCompile(`^\s*(https?://[^\s]+)\s*$`)

// ExpandContent processes markdown content and expands special URLs
// This should be called BEFORE parsing with goldmark
func (e *Expander) ExpandContent(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inCodeBlock := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Track code blocks to avoid processing URLs inside them
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			continue
		}

		if inCodeBlock {
			result = append(result, line)
			continue
		}

		// Check if line is a URL-only line
		if match := urlOnlyLinePattern.FindStringSubmatch(line); match != nil {
			url := match[1]
			expanded := e.expandURL(url)
			result = append(result, expanded)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// expandURL expands a URL based on its type
func (e *Expander) expandURL(url string) string {
	// Check if it's a GitHub URL
	if _, ok := ParseGitHubURL(url); ok {
		expanded, err := e.github.ExpandToCodeBlock(url)
		if err == nil {
			return expanded
		}
		// Fall through to OG card on error
	}

	// Check if it's a Twitter/X URL
	if IsTwitterURL(url) {
		expanded, err := e.twitter.GenerateEmbed(url)
		if err == nil {
			return expanded
		}
		// Fall through to OG card on error
	}

	// Default: generate OG card
	card, err := e.ogcard.GenerateOGCard(url)
	if err == nil {
		return card
	}

	// If all else fails, return as a simple link
	return "[" + url + "](" + url + ")"
}

// imageWidthPattern matches custom image width syntax: ![alt](url){width=500}
var imageWidthPattern = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)\{width=(\d+)\}`)

// ProcessImageWidths converts custom image width syntax to HTML
func ProcessImageWidths(content string) string {
	return imageWidthPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := imageWidthPattern.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}
		alt := parts[1]
		src := parts[2]
		width := parts[3]
		return `<img src="` + escapeHTML(src) + `" alt="` + escapeHTML(alt) + `" width="` + width + `" loading="lazy">`
	})
}
