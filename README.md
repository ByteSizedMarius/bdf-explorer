## bdf-explorer

Quick BDF font utility, because I needed one. Forked from [zachomedia/go-bdf](https://github.com/zachomedia/go-bdf?tab=readme-ov-file) - thanks!

### commandline usage

```
bdf-explorer -font "<path>" 
```

Generates a png with a grid of all characters.

```
bdf-explorer -font "<path>"  -export
```

Generates a png with a grid of all and exports each character as a separate png file into a subdirectory.

### library usage

```
go get github.com/ByteSizedMarius/bdf-explorer
```

```go
package main

import "github.com/ByteSizedMarius/bdf-explorer/bdf"

func main() {
	// load a file
	font, err := bdf.FromFile("<path>")

	// etc.
}
```

```go
type Face struct {
    Font *Font
}

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

type Character struct {
    Name       string
    Encoding   rune
    Advance    [2]int
    Alpha      *image.Alpha
    LowerPoint [2]int
}
```