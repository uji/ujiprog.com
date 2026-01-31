package ogimage

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Generator generates OG images for articles (WASM compatible)
type Generator struct {
	templateImg  image.Image
	asciiFace    font.Face
	japaneseFace font.Face
	fontSize     float64
}

// NewGenerator creates a new OG image generator from byte slices
func NewGenerator(templateData, asciiFontData, japaneseFontData []byte, fontSize float64) (*Generator, error) {
	// Decode template image
	templateImg, err := png.Decode(bytes.NewReader(templateData))
	if err != nil {
		return nil, err
	}

	gen := &Generator{
		templateImg: templateImg,
		fontSize:    fontSize,
	}

	// Load ASCII font
	if len(asciiFontData) > 0 {
		face, err := loadFontFace(asciiFontData, fontSize)
		if err == nil {
			gen.asciiFace = face
		}
	}

	// Load Japanese font
	if len(japaneseFontData) > 0 {
		face, err := loadFontFace(japaneseFontData, fontSize)
		if err == nil {
			gen.japaneseFace = face
		}
	}

	return gen, nil
}

// loadFontFace loads a font from byte slice and returns a font.Face
func loadFontFace(data []byte, size float64) (font.Face, error) {
	f, err := sfnt.Parse(data)
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
func (g *Generator) drawTextDirect(dst draw.Image, line string, centerX, y float64, textColor color.Color) {
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

// Generate creates an OG image with the given title and writes it to the writer
func (g *Generator) Generate(title string, w io.Writer) error {
	bounds := g.templateImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create RGBA image and draw template as background
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, g.templateImg, bounds.Min, draw.Src)

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
	lineHeight := g.fontSize * 1.5
	totalTextHeight := float64(len(allLines)) * lineHeight

	// Draw each line centered using direct font.Drawer (accurate glyph coverage)
	startY := textAreaY - totalTextHeight/2 + lineHeight/2
	for i, line := range allLines {
		y := startY + float64(i)*lineHeight
		g.drawTextDirect(dst, line, textAreaX, y, textColor)
	}

	// Encode to PNG
	return png.Encode(w, dst)
}

// wrapTextMultiFont wraps text considering multiple fonts
func (g *Generator) wrapTextMultiFont(text string, maxWidth float64) []string {
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
func (g *Generator) measureLineWidth(line string) float64 {
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
