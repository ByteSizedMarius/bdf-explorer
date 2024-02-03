package bdf

import (
	"image"
	"image/color"
	"math"
	"sort"
)

// drawGridLines draws grid lines around the character based on its position.
func drawGridLines(img *image.RGBA, x *int, y *int, maxWidth, maxHeight, padding int) {
	// Vertical grid line
	for i := *y; i < *y+maxHeight+2*padding; i++ {
		img.Set(*x-1, i, color.RGBA{255, 255, 255, 255}) // Set just before moving x back for the next character
	}

	// Horizontal grid line at the bottom of the current row
	if *y+maxHeight+2*padding < img.Bounds().Dy() {
		for i := 0; i < img.Bounds().Dx(); i++ {
			img.Set(i, *y+maxHeight+2*padding, color.RGBA{255, 255, 255, 255})
		}
	}
}

// calculateImageSize calculates the total image size, considering dynamic character padding.
func calculateImageSize(numRows, numCols, maxWidth, maxHeight, padding int) (int, int) {
	gridLineWidth := 1
	totalWidth := numCols*(maxWidth+2*padding+gridLineWidth) - gridLineWidth
	totalHeight := numRows*(maxHeight+2*padding+gridLineWidth) - gridLineWidth
	return totalWidth, totalHeight
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// calculateLayout calculates the layout of characters in the image.
func calculateLayout(characters []Character) (int, int, int, int) {
	numCharacters := len(characters)
	numRows := int(math.Sqrt(float64(numCharacters)))
	numCols := (numCharacters + numRows - 1) / numRows
	maxWidth, maxHeight := getMaxCharacterSize(characters)
	return numRows, numCols, maxWidth, maxHeight
}

// getMaxCharacterSize finds the maximum width and height among all characters.
func getMaxCharacterSize(characters []Character) (int, int) {
	maxWidth, maxHeight := 0, 0
	for _, char := range characters {
		if alpha := char.Alpha; alpha != nil {
			bounds := alpha.Bounds()
			maxWidth, maxHeight = max(maxWidth, bounds.Dx()), max(maxHeight, bounds.Dy())
		}
	}
	return maxWidth, maxHeight
}

// sortCharacters sorts the characters by their encoding values.
func sortCharacters(characters []Character) {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Encoding < characters[j].Encoding
	})
}
