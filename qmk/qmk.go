package qmk

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func ImgToBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X
	lines := height / 8

	data := make([]byte, 0, width*lines)

	for line := 0; line < lines; line++ {
		for x := 0; x < width; x++ {
			var v byte = 0
			for i := 0; i < 8; i++ {
				y := line*8 + i
				b := img.At(x, y)
				p := colorToPixel(b)
				v |= (p << i)
			}
			data = append(data, byte(v))
		}
	}

	return data, nil
}

func ParseImage(path string, name string) error {
	data, err := ImgToBytes(path)
	if err != nil {
		return err
	}

	if name == "" {
		// get filename without path or extension
		name = filepath.Base(path)
		ext := filepath.Ext(name)
		name = name[:len(name)-len(ext)]
	}

	printCode(name, data)

	return nil
}

func colorToPixel(c color.Color) byte {
	r, g, b, _ := c.RGBA()
	v := r + g + b
	if v >= 127 {
		return 1
	}
	return 0
}

func printCode(name string, data []byte) {
	fmt.Printf("static const char PROGMEM %s[] = {\n", name)
	for i, b := range data {
		if i%16 == 0 {
			fmt.Print("\t")
		}
		fmt.Printf("0x%02x, ", b)
		if i%16 == 15 {
			fmt.Println()
		}
	}
	fmt.Println("};")
}
