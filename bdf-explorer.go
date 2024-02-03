package main

import (
	"flag"
	"fmt"
	"github.com/ByteSizedMarius/bdf-explorer/bdf"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define command-line flags
	fontFile := flag.String("font", "", "Path to the font file")
	exportAll := flag.Bool("export", false, "Export all characters to individual images")
	flag.Parse()

	// Check if font file is provided
	if *fontFile == "" {
		fmt.Println("Please provide a font file")
		os.Exit(1)
	}

	// Load the font
	font, err := bdf.FromFile(*fontFile)
	if err != nil {
		fmt.Println("Error loading font:", err)
		os.Exit(1)
	}
	fontName := filepath.Base(*fontFile)
	fontName = strings.TrimSuffix(fontName, filepath.Ext(fontName))
	fmt.Println("Font loaded:", fontName)

	// Render the file
	imgP := fontName + ".png"
	err = font.RenderFontImage(imgP, 5)
	if err != nil {
		fmt.Println("Error rendering font:", err)
		os.Exit(1)
	}
	fmt.Println("Font image saved to", imgP)

	if !*exportAll {
		return
	}

	// Create a directory for the single imgs
	err = os.Mkdir(fontName, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	// Render each character
	for i, char := range font.Characters {
		err = char.RenderCharacterImage(fontName)
		if err != nil {
			fmt.Println("Error rendering character:", err)
			os.Exit(1)
		}
		fmt.Printf("\rRendering character %d of %d", i+1, len(font.Characters))
	}
	fmt.Println("\nDone")

}
