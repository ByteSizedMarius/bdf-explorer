package bdf

import (
	"fmt"
	"image"
	"image/draw"
	"path/filepath"
)

type Character struct {
	Name       string
	Encoding   rune
	Advance    [2]int
	Alpha      *image.Alpha
	LowerPoint [2]int
}

func (c *Character) RenderCharacterImage(path string) error {
	if alpha := c.Alpha; alpha != nil {
		bounds := alpha.Bounds()
		charWidth := bounds.Dx()
		charHeight := bounds.Dy()

		// Create a new image with the size of the character
		img := image.NewRGBA(image.Rect(0, 0, charWidth, charHeight))

		// Draw the character onto the image
		draw.Draw(img, bounds, alpha, bounds.Min, draw.Over)

		// Save the image to the output file
		return saveImage(filepath.Join(path, fmt.Sprintf("%X.png", c.Encoding)), img)
	}

	return nil
}
