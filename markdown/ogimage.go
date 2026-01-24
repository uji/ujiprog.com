package markdown

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
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

	f, err := sfnt.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
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

// drawTextDirect draws text directly using font.Drawer for accurate glyph coverage
// This method uses FreeType's coverage bitmap directly, avoiding quantization errors
func (g *OGImageGenerator) drawTextDirect(dst draw.Image, line string, centerX, y float64, textColor color.Color) {
	segments := segmentText(line)
	totalWidth := g.measureLineWidth(line)

	// Calculate starting X position for centered text
	x := int(centerX - totalWidth/2)

	// Calculate baseline Y position (y is vertical center, adjust for ascent)
	var ascent fixed.Int26_6
	if g.japaneseFace != nil {
		ascent = g.japaneseFace.Metrics().Ascent
	} else if g.asciiFace != nil {
		ascent = g.asciiFace.Metrics().Ascent
	}
	yBaseline := int(y) + ascent.Round()/2

	src := image.NewUniform(textColor)

	for _, seg := range segments {
		var face font.Face
		if seg.isJapanese && g.japaneseFace != nil {
			face = g.japaneseFace
		} else if g.asciiFace != nil {
			face = g.asciiFace
		} else {
			continue
		}

		d := &font.Drawer{
			Dst:  dst,
			Src:  src,
			Face: face,
			Dot:  fixed.P(x, yBaseline),
		}
		d.DrawString(seg.text)
		x = d.Dot.X.Round()
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

	// Create RGBA image and draw template as background
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, templateImg, bounds.Min, draw.Src)

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
		wrapped := g.wrapTextMultiFont(line, maxWidth)
		allLines = append(allLines, wrapped...)
	}

	// Calculate total height of text block
	lineHeight := g.fontConfig.FontSize * 1.5
	totalTextHeight := float64(len(allLines)) * lineHeight

	// Draw each line centered using direct font.Drawer (accurate glyph coverage)
	startY := textAreaY - totalTextHeight/2 + lineHeight/2
	for i, line := range allLines {
		y := startY + float64(i)*lineHeight
		g.drawTextDirect(dst, line, textAreaX, y, textColor)
	}

	// Save to file
	return savePNG(outputPath, dst)
}

// wrapTextMultiFont wraps text considering multiple fonts
func (g *OGImageGenerator) wrapTextMultiFont(text string, maxWidth float64) []string {
	var lines []string

	// Measure and wrap character by character for accuracy with mixed fonts
	runes := []rune(text)
	if len(runes) == 0 {
		return lines
	}

	currentLine := ""
	for _, r := range runes {
		testLine := currentLine + string(r)
		lineWidth := g.measureLineWidth(testLine)

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
func (g *OGImageGenerator) measureLineWidth(line string) float64 {
	segments := segmentText(line)
	var totalWidth fixed.Int26_6

	for _, seg := range segments {
		var face font.Face
		if seg.isJapanese && g.japaneseFace != nil {
			face = g.japaneseFace
		} else if g.asciiFace != nil {
			face = g.asciiFace
		} else {
			continue
		}
		w := font.MeasureString(face, seg.text)
		totalWidth += w
	}

	return float64(totalWidth) / 64.0 // Convert from fixed.Int26_6 to float64
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

// savePNG saves an image to a PNG file
func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
