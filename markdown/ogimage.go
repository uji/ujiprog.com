package markdown

import (
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

// OGImageGenerator generates OG images for articles
type OGImageGenerator struct {
	templatePath string
	fontPath     string
}

// NewOGImageGenerator creates a new OG image generator
func NewOGImageGenerator(templatePath, fontPath string) *OGImageGenerator {
	return &OGImageGenerator{
		templatePath: templatePath,
		fontPath:     fontPath,
	}
}

// Generate creates an OG image with the given title and saves it to outputPath
func (g *OGImageGenerator) Generate(title, outputPath string) error {
	// Load template image
	templateImg, err := loadPNG(g.templatePath)
	if err != nil {
		return err
	}

	bounds := templateImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create drawing context
	dc := gg.NewContext(width, height)

	// Draw template as background
	dc.DrawImage(templateImg, 0, 0)

	// Load font
	if err := dc.LoadFontFace(g.fontPath, 48); err != nil {
		return err
	}

	// Set text color (dark gray to match site design)
	dc.SetRGB255(74, 75, 74)

	// Calculate text area (white region in the center)
	// Assuming the white area is roughly in the center
	textAreaX := float64(width) / 2
	textAreaY := float64(height) / 2
	maxWidth := float64(width) * 0.7

	// Wrap text if necessary
	lines := wrapText(dc, title, maxWidth)

	// Calculate total height of text block
	lineHeight := 60.0
	totalTextHeight := float64(len(lines)) * lineHeight

	// Draw each line centered
	startY := textAreaY - totalTextHeight/2 + lineHeight/2
	for i, line := range lines {
		y := startY + float64(i)*lineHeight
		dc.DrawStringAnchored(line, textAreaX, y, 0.5, 0.5)
	}

	// Save to file
	return dc.SavePNG(outputPath)
}

// loadPNG loads a PNG image from file
func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return png.Decode(f)
}

// wrapText wraps text to fit within maxWidth
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		testLine := currentLine + " " + word
		w, _ := dc.MeasureString(testLine)
		if w <= maxWidth {
			currentLine = testLine
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	// Also handle very long words by breaking them
	var result []string
	for _, line := range lines {
		w, _ := dc.MeasureString(line)
		if w <= maxWidth {
			result = append(result, line)
		} else {
			// Break long line character by character
			runes := []rune(line)
			currentPart := ""
			for _, r := range runes {
				testPart := currentPart + string(r)
				pw, _ := dc.MeasureString(testPart)
				if pw <= maxWidth {
					currentPart = testPart
				} else {
					if currentPart != "" {
						result = append(result, currentPart)
					}
					currentPart = string(r)
				}
			}
			if currentPart != "" {
				result = append(result, currentPart)
			}
		}
	}

	return result
}
