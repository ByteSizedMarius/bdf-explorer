package bdf

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

type Font struct {
	Name        string
	Size        int
	PixelSize   int
	DPI         [2]int
	BPP         int
	Ascent      int
	Descent     int
	CapHeight   int
	XHeight     int
	Characters  []Character
	CharMap     map[rune]*Character
	Encoding    string
	DefaultChar rune
}

func (f *Font) lookup(r rune) *Character {
	c, ok := f.CharMap[r]
	if !ok {
		c, ok = f.CharMap[f.DefaultChar]
		if !ok {
			return nil
		}
	}
	return c
}

// RenderFontImage renders the characters of a font into a single image with grid lines and padding.
func (f *Font) RenderFontImage(output string, padding int) error {
	sortCharacters(f.Characters)
	numRows, numCols, maxWidth, maxHeight := calculateLayout(f.Characters)
	totalWidth, totalHeight := calculateImageSize(numRows, numCols, maxWidth, maxHeight, padding)
	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))
	f.drawCharacters(img, maxWidth, maxHeight, padding)
	return saveImage(output, img)
}

// drawCharacters draws each character on the image.
func (f *Font) drawCharacters(img *image.RGBA, maxWidth, maxHeight, padding int) {
	x, y := 0, 0
	for _, char := range f.Characters {
		f.drawCharacter(img, &char, &x, &y, maxWidth, maxHeight, padding)
	}
}

// drawCharacter draws a single character with configurable padding and updates the current position.
func (f *Font) drawCharacter(img *image.RGBA, char *Character, x *int, y *int, maxWidth, maxHeight, padding int) {
	if alpha := char.Alpha; alpha != nil {
		bounds := alpha.Bounds()
		charWidth := bounds.Dx()
		charHeight := bounds.Dy()

		// Calculate destination rectangle for the character considering dynamic padding
		destRect := image.Rect(*x+padding, *y+padding, *x+charWidth+padding, *y+charHeight+padding)
		draw.Draw(img, destRect, alpha, bounds.Min, draw.Over)

		// Update x-coordinate for the next character, including padding and grid line width
		*x += maxWidth + 2*padding + 1 // Adjust for dynamic padding and 1 pixel for the grid line

		// Check if we need to move to the next row
		if *x >= img.Bounds().Dx()-maxWidth {
			*x = 0
			*y += maxHeight + 2*padding + 1 // Adjust for dynamic padding and 1 pixel for the grid line
		}

		// Draw grid lines around the character
		drawGridLines(img, x, y, maxWidth, maxHeight, padding)
	}
}

// saveImage saves the image to a file.
func saveImage(output string, img *image.RGBA) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
