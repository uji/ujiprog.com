package markdown

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// GitHubCodeExpander handles expansion of GitHub blob URLs to code blocks
type GitHubCodeExpander struct {
	client *http.Client
}

// NewGitHubCodeExpander creates a new GitHubCodeExpander
func NewGitHubCodeExpander() *GitHubCodeExpander {
	return &GitHubCodeExpander{
		client: &http.Client{},
	}
}

// githubURLPattern matches GitHub blob URLs with optional line ranges
// Format: https://github.com/{owner}/{repo}/blob/{ref}/{path}#L{start}-L{end}
var githubURLPattern = regexp.MustCompile(
	`^https://github\.com/([^/]+)/([^/]+)/blob/([^/]+)/(.+?)(?:#L(\d+)(?:-L(\d+))?)?$`,
)

// GitHubURLInfo contains parsed information from a GitHub URL
type GitHubURLInfo struct {
	Owner     string
	Repo      string
	Ref       string
	Path      string
	StartLine int
	EndLine   int
}

// ParseGitHubURL parses a GitHub blob URL and extracts its components
func ParseGitHubURL(url string) (*GitHubURLInfo, bool) {
	matches := githubURLPattern.FindStringSubmatch(url)
	if matches == nil {
		return nil, false
	}

	info := &GitHubURLInfo{
		Owner:     matches[1],
		Repo:      matches[2],
		Ref:       matches[3],
		Path:      matches[4],
		StartLine: 0,
		EndLine:   0,
	}

	if matches[5] != "" {
		info.StartLine, _ = strconv.Atoi(matches[5])
	}
	if matches[6] != "" {
		info.EndLine, _ = strconv.Atoi(matches[6])
	} else if info.StartLine > 0 {
		info.EndLine = info.StartLine
	}

	return info, true
}

// FetchCode fetches code from a GitHub raw URL and extracts the specified lines
func (g *GitHubCodeExpander) FetchCode(info *GitHubURLInfo) (string, error) {
	rawURL := fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/%s/%s/%s",
		info.Owner, info.Repo, info.Ref, info.Path,
	)

	resp, err := g.client.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch code: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	content := string(body)

	if info.StartLine > 0 {
		lines := strings.Split(content, "\n")
		start := info.StartLine - 1
		end := info.EndLine
		if start < 0 {
			start = 0
		}
		if end > len(lines) {
			end = len(lines)
		}
		if start >= len(lines) {
			return "", fmt.Errorf("start line %d exceeds file length", info.StartLine)
		}
		content = strings.Join(lines[start:end], "\n")
	}

	return content, nil
}

// GetLanguage returns the language based on file extension
func GetLanguage(path string) string {
	ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
	langMap := map[string]string{
		"go":    "go",
		"js":    "javascript",
		"ts":    "typescript",
		"py":    "python",
		"rb":    "ruby",
		"rs":    "rust",
		"java":  "java",
		"c":     "c",
		"cpp":   "cpp",
		"h":     "c",
		"hpp":   "cpp",
		"html":  "html",
		"css":   "css",
		"json":  "json",
		"yaml":  "yaml",
		"yml":   "yaml",
		"md":    "markdown",
		"sh":    "bash",
		"bash":  "bash",
		"sql":   "sql",
		"xml":   "xml",
		"toml":  "toml",
		"swift": "swift",
		"kt":    "kotlin",
	}
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ""
}

// ExpandToCodeBlock fetches code and returns a markdown code block
func (g *GitHubCodeExpander) ExpandToCodeBlock(url string) (string, error) {
	info, ok := ParseGitHubURL(url)
	if !ok {
		return "", fmt.Errorf("invalid GitHub URL: %s", url)
	}

	code, err := g.FetchCode(info)
	if err != nil {
		return "", err
	}

	lang := GetLanguage(info.Path)
	lineInfo := ""
	if info.StartLine > 0 {
		if info.EndLine > info.StartLine {
			lineInfo = fmt.Sprintf(" (L%d-L%d)", info.StartLine, info.EndLine)
		} else {
			lineInfo = fmt.Sprintf(" (L%d)", info.StartLine)
		}
	}

	return fmt.Sprintf("```%s\n// %s/%s%s\n%s\n```",
		lang, info.Repo, info.Path, lineInfo, code), nil
}
