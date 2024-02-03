package bdf

import (
	"bufio"
	"golang.org/x/text/encoding/charmap"
	"strconv"
	"strings"
)

func parseGlobalsAndProperties(s *bufio.Scanner, f *Font) error {
	var err error

	var registry string
	var encoding string
	var defaultChar int

scan:
	for s.Scan() {
		components := strings.Split(s.Text(), " ")
		switch components[0] {
		case "FONT":
			f.Name = components[1]
		case "SIZE":
			f.Size, err = strconv.Atoi(components[1])
			if err != nil {
				return err
			}

			f.DPI[0], err = strconv.Atoi(components[2])
			if err != nil {
				return err
			}

			f.DPI[1], err = strconv.Atoi(components[3])
			if err != nil {
				return err
			}

			if len(components) > 4 {
				f.BPP, err = strconv.Atoi(components[4])
				if err != nil {
					return err
				}
			}
		case "CHARSET_REGISTRY":
			registry = components[1]
		case "CHARSET_ENCODING":
			encoding = components[1]
		case "PIXEL_SIZE":
			f.PixelSize, err = strconv.Atoi(components[1])
		case "FONT_ASCENT":
			f.Ascent, err = strconv.Atoi(components[1])
			if err != nil {
				return err
			}
		case "FONT_DESCENT":
			f.Descent, err = strconv.Atoi(components[1])
			if err != nil {
				return err
			}
		case "CAP_HEIGHT":
			f.CapHeight, err = strconv.Atoi(components[1])
			if err != nil {
				return err
			}
		case "X_HEIGHT":
			f.XHeight, err = strconv.Atoi(components[1])
			if err != nil {
				return err
			}
		case "DEFAULT_CHAR":
			defaultChar, err = strconv.Atoi(components[1])
			if err != nil {
				return err
			}
		case "CHARS":
			count, err := strconv.Atoi(components[1])
			if err != nil {
				return err
			}
			f.Characters = make([]Character, count)
			break scan
		}
	}

	f.Encoding = registry + "-" + encoding
	f.DefaultChar = charToRune(f.Encoding, defaultChar)

	return nil
}

func charToRune(encoding string, char int) rune {
	if charMap := findCharmap(encoding); charMap != nil {
		return charMap.DecodeByte(byte(char))
	}
	return rune(char)
}

func findCharmap(requested string) *charmap.Charmap {
	trimmed := strings.TrimSpace(strings.ToLower(requested))

	knownMaps := map[string]*charmap.Charmap{
		"iso8859-1":  charmap.ISO8859_1,
		"iso8859-2":  charmap.ISO8859_2,
		"iso8859-9":  charmap.ISO8859_9,
		"iso8859-15": charmap.ISO8859_15,
	}

	charMap := knownMaps[trimmed]
	return charMap
}

func bitAt(xs []byte, i int) byte {
	return (xs[i>>3] >> (7 - (i % 8))) & 1
}
