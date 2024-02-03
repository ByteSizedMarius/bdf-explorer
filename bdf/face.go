package bdf

import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
)

type Face struct {
	Font *Font
}

func (f *Face) Close() error { return nil }

func (f *Face) Metrics() font.Metrics {
	return font.Metrics{
		Ascent:    fixed.I(f.Font.Ascent),
		Descent:   fixed.I(f.Font.Descent),
		CapHeight: fixed.I(f.Font.CapHeight),
		XHeight:   fixed.I(f.Font.XHeight),
		Height:    fixed.I(f.Font.Ascent + f.Font.Descent),
	}
}

func (f *Face) Kern(_, _ rune) fixed.Int26_6 {
	return 0
}

func (f *Face) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	c := f.Font.lookup(r)
	if c == nil {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}

	mask = c.Alpha

	x := int(dot.X)>>6 + c.LowerPoint[0]
	y := int(dot.Y)>>6 - c.LowerPoint[1]
	dr = image.Rectangle{
		Min: image.Point{
			X: x,
			Y: y - c.Alpha.Rect.Max.Y,
		},
		Max: image.Point{
			X: x + c.Alpha.Stride,
			Y: y,
		},
	}

	return dr, mask, image.Point{Y: 0}, fixed.I(c.Advance[0]), true
}

func (f *Face) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	c := f.Font.lookup(r)
	if c == nil {
		return fixed.R(0, -f.Font.Ascent, 0, +f.Font.Descent), 0, false
	}

	return fixed.R(c.LowerPoint[0], -f.Font.Ascent, c.LowerPoint[0]+c.Alpha.Rect.Dx(), f.Font.Descent), fixed.I(c.Advance[0]), true
}

func (f *Face) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	c := f.Font.lookup(r)
	if c == nil {
		return 0, false
	}
	return fixed.I(c.Advance[0]), true
}
