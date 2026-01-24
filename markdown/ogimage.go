package markdown

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// FontConfig holds font configuration for OG image generation
type FontConfig struct {
	ASCIIFontPath    string
	JapaneseFontPath string
	FontSize         float64
}

// OGImageGenerator generates OG images for articles
type OGImageGenerator struct {
	templatePath string
	fontConfig   FontConfig
	asciiFace    font.Face
	japaneseFace font.Face
}

// NewOGImageGenerator creates a new OG image generator with multiple fonts
func NewOGImageGenerator(templatePath string, fontConfig FontConfig) *OGImageGenerator {
	gen := &OGImageGenerator{
		templatePath: templatePath,
		fontConfig:   fontConfig,
	}

	// Load ASCII font
	if fontConfig.ASCIIFontPath != "" {
		if face, err := loadFontFace(fontConfig.ASCIIFontPath, fontConfig.FontSize); err == nil {
			gen.asciiFace = face
		}
	}

	// Load Japanese font
	if fontConfig.JapaneseFontPath != "" {
		if face, err := loadFontFace(fontConfig.JapaneseFontPath, fontConfig.FontSize); err == nil {
			gen.japaneseFace = face
		}
	}

	return gen
}

// loadFontFace loads a font file and returns a font.Face
func loadFontFace(path string, size float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

// textSegment represents a segment of text with its font type
type textSegment struct {
	text       string
	isJapanese bool
}

// isJapanese checks if a rune is a Japanese character
func isJapanese(r rune) bool {
	return (r >= 0x3040 && r <= 0x309F) || // ひらがな
		(r >= 0x30A0 && r <= 0x30FF) || // カタカナ
		(r >= 0x4E00 && r <= 0x9FFF) || // 漢字
		(r >= 0x3000 && r <= 0x303F) || // CJK句読点
		(r >= 0xFF00 && r <= 0xFFEF) // 全角文字
}

// segmentText splits text into segments based on character type (ASCII vs Japanese)
func segmentText(text string) []textSegment {
	var segments []textSegment
	var current strings.Builder
	var currentIsJapanese *bool

	for _, r := range text {
		isJP := isJapanese(r)

		if currentIsJapanese == nil {
			currentIsJapanese = &isJP
			current.WriteRune(r)
		} else if *currentIsJapanese == isJP {
			current.WriteRune(r)
		} else {
			// Character type changed, save current segment
			segments = append(segments, textSegment{
				text:       current.String(),
				isJapanese: *currentIsJapanese,
			})
			current.Reset()
			currentIsJapanese = &isJP
			current.WriteRune(r)
		}
	}

	// Add final segment
	if current.Len() > 0 && currentIsJapanese != nil {
		segments = append(segments, textSegment{
			text:       current.String(),
			isJapanese: *currentIsJapanese,
		})
	}

	return segments
}

// drawTextWithMergedMask draws text using a luminance-based alpha mask
// Renders black text on white background, then uses inverted luminance as alpha
// This avoids transparency accumulation issues at glyph intersections
func (g *OGImageGenerator) drawTextWithMergedMask(dc *gg.Context, line string, centerX, y float64, textColor color.Color) {
	segments := segmentText(line)

	// Calculate total width first
	totalWidth := g.measureLineWidth(dc, line)

	// Create a temporary context with white background for rendering
	bounds := dc.Image().Bounds()
	tmpDc := gg.NewContext(bounds.Dx(), bounds.Dy())

	// Fill with white background
	tmpDc.SetRGB(1, 1, 1)
	tmpDc.Clear()

	// Set black color for text
	tmpDc.SetRGB(0, 0, 0)

	// Start drawing from left side of centered text
	x := centerX - totalWidth/2

	for _, seg := range segments {
		var face font.Face
		if seg.isJapanese && g.japaneseFace != nil {
			face = g.japaneseFace
		} else if g.asciiFace != nil {
			face = g.asciiFace
		} else {
			continue
		}

		tmpDc.SetFontFace(face)
		tmpDc.DrawStringAnchored(seg.text, x, y, 0, 0.5)
		w, _ := tmpDc.MeasureString(seg.text)
		x += w
	}

	// Extract alpha mask from the rendered image's luminance
	tmpImg := tmpDc.Image()
	mask := image.NewAlpha(bounds)

	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			// Get pixel color from temporary image
			c := tmpImg.At(px, py)
			r, _, _, _ := c.RGBA()
			// Convert to 8-bit and invert: white (255) -> 0 alpha, black (0) -> 255 alpha
			alpha := uint8(255 - (r >> 8))
			mask.SetAlpha(px, py, color.Alpha{A: alpha})
		}
	}

	// Draw the mask onto the context with the text color
	dst := dc.Image()
	if rgba, ok := dst.(*image.RGBA); ok {
		src := image.NewUniform(textColor)
		draw.DrawMask(rgba, bounds, src, image.Point{}, mask, image.Point{}, draw.Over)
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

	// Text color (dark gray to match site design)
	textColor := color.RGBA{R: 74, G: 75, B: 74, A: 255}

	// Calculate text area
	textAreaX := float64(width) / 2
	textAreaY := float64(height) / 2
	maxWidth := float64(width) * 0.7

	// Split by explicit newlines first
	explicitLines := strings.Split(title, "\n")

	// Wrap each line if necessary
	var allLines []string
	for _, line := range explicitLines {
		wrapped := g.wrapTextMultiFont(dc, line, maxWidth)
		allLines = append(allLines, wrapped...)
	}

	// Calculate total height of text block
	lineHeight := g.fontConfig.FontSize * 1.25
	totalTextHeight := float64(len(allLines)) * lineHeight

	// Draw each line centered using merged mask method (prevents transparency at glyph intersections)
	startY := textAreaY - totalTextHeight/2 + lineHeight/2
	for i, line := range allLines {
		y := startY + float64(i)*lineHeight
		g.drawTextWithMergedMask(dc, line, textAreaX, y, textColor)
	}

	// Save to file
	return dc.SavePNG(outputPath)
}

// wrapTextMultiFont wraps text considering multiple fonts
func (g *OGImageGenerator) wrapTextMultiFont(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string

	// Use Japanese font for measurement if available (typically wider)
	if g.japaneseFace != nil {
		dc.SetFontFace(g.japaneseFace)
	} else if g.asciiFace != nil {
		dc.SetFontFace(g.asciiFace)
	}

	// Measure and wrap character by character for accuracy with mixed fonts
	runes := []rune(text)
	if len(runes) == 0 {
		return lines
	}

	currentLine := ""
	for _, r := range runes {
		testLine := currentLine + string(r)
		lineWidth := g.measureLineWidth(dc, testLine)

		if lineWidth <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = string(r)
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// measureLineWidth measures the width of a line considering multiple fonts
func (g *OGImageGenerator) measureLineWidth(dc *gg.Context, line string) float64 {
	segments := segmentText(line)
	totalWidth := 0.0

	for _, seg := range segments {
		if seg.isJapanese && g.japaneseFace != nil {
			dc.SetFontFace(g.japaneseFace)
		} else if g.asciiFace != nil {
			dc.SetFontFace(g.asciiFace)
		}
		w, _ := dc.MeasureString(seg.text)
		totalWidth += w
	}

	return totalWidth
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
