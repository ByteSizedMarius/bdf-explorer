package bdf

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"image"
	"os"
	"strconv"
	"strings"
)

// FromFile loads the Font from the given file.
func FromFile(path string) (f *Font, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		err = errors.Join(errors.New("error on read"), err)
		return
	}

	f, err = Parse(data)
	if err != nil {
		err = errors.Join(errors.New("error on parse"), err)
		return
	}

	return
}

// Parse parses the BDF font data.
func Parse(data []byte) (*Font, error) {
	reader := bytes.NewReader(data)
	s := bufio.NewScanner(reader)

	f := Font{
		CharMap:     make(map[rune]*Character),
		DefaultChar: 32,
		BPP:         1,
	}

	var err error

	err = parseGlobalsAndProperties(s, &f)
	if err != nil {
		return nil, err
	}

	charMap := findCharmap(f.Encoding)

	char := -1
	row := -1
	inBitmap := false
	for s.Scan() {
		components := strings.Split(s.Text(), " ")

		if !inBitmap {
			switch components[0] {

			case "STARTCHAR":
				char++
				f.Characters[char].Name = components[1]
			case "ENCODING":
				var code int
				code, err = strconv.Atoi(components[1])
				if err != nil {
					return nil, err
				}

				var r rune
				if charMap != nil {
					r = charMap.DecodeByte(byte(code))
				} else {
					r = rune(code)
				}
				f.Characters[char].Encoding = r
				f.CharMap[r] = &f.Characters[char]
			case "DWIDTH":
				f.Characters[char].Advance[0], err = strconv.Atoi(components[1])
				if err != nil {
					return nil, err
				}

				f.Characters[char].Advance[1], err = strconv.Atoi(components[2])
				if err != nil {
					return nil, err
				}
			case "BBX":
				var w, h int
				w, err = strconv.Atoi(components[1])
				if err != nil {
					return nil, err
				}

				h, err = strconv.Atoi(components[2])
				if err != nil {
					return nil, err
				}

				// Lower-left corner
				var lx, ly int
				lx, err = strconv.Atoi(components[3])
				if err != nil {
					return nil, err
				}
				ly, err = strconv.Atoi(components[4])
				if err != nil {
					return nil, err
				}

				f.Characters[char].LowerPoint[0] = lx
				f.Characters[char].LowerPoint[1] = ly

				f.Characters[char].Alpha = &image.Alpha{
					Stride: w,
					Rect: image.Rectangle{
						Max: image.Point{
							X: w,
							Y: h,
						},
					},
					Pix: make([]byte, w*h),
				}
			case "BITMAP":
				inBitmap = true
				row = -1
			}
		} else {
			if components[0] == "ENDCHAR" {
				inBitmap = false
				continue
			}

			row = row + 1
			var b []byte
			b, err = hex.DecodeString(s.Text())
			if err != nil {
				return nil, err
			}

			for i := 0; i < f.Characters[char].Alpha.Stride; i++ {
				val := byte(0x00)
				for j := 0; j < f.BPP; j++ {
					val <<= 1
					val |= bitAt(b, i*f.BPP+j)
				}
				f.Characters[char].Alpha.Pix[row*f.Characters[char].Alpha.Stride+i] = byte(uint32(val) * 0xff / ((1 << f.BPP) - 1))
			}
		}
	}

	return &f, nil
}
