package frenyard
import (
	"image"
	"image/color"
	"image/png"
	"encoding/base64"
	"strings"
)

/*
 * Converts a Go colour into an ARGB colour uint32 as losslessly as possible.
 */
func fyConvertGoImageColour(col color.Color) uint32 {
	// Actual color appears to be per-pixel. Great. More problems.
	switch colConv := col.(type) {
		case color.NRGBA:
			return ColourFromARGB(colConv.A, colConv.R, colConv.G, colConv.B)
		case color.RGBA:
			// This is still premultiplied and thus NOT GOOD, but let's roll with it
			a := colConv.A
			r := colConv.R
			g := colConv.G
			b := colConv.B
			if ((a != 0) && (a != 255)) {
				// Alpha not 0 (would /0) or a NOP
				a16 := uint16(a)
				r = uint8((uint16(colConv.R) * 255) / a16)
				g = uint8((uint16(colConv.G) * 255) / a16)
				b = uint8((uint16(colConv.B) * 255) / a16)
			}
			return ColourFromARGB(a, r, g, b)
		default: {
			// Give up and work backwards from premultiplied (NOT GOOD!!!)
			r, g, b, a := col.RGBA()
			if ((a != 0) && (a != 0xFFFF)) {
				// Scaling, also implicitly does the 8-bit conversion
				r *= 255
				g *= 255
				b *= 255
				r /= a
				g /= a
				b /= a
			} else {
				// Convert to 8-bit
				r >>= 8
				g >>= 8
				b >>= 8
			}
			// Convert alpha to 8-bit
			a >>= 8
			return ColourFromARGB(uint8(a), uint8(r), uint8(g), uint8(b))
		}
	}
}

/*
 * Imports an image from Go's "image" library to a texture.
 */
func GoImageToTexture(img image.Image) Texture {
	sizePreTranslate := img.Bounds().Size()
	pixels := make([]uint32, sizePreTranslate.X * sizePreTranslate.Y)
	size := Vec2i{int32(sizePreTranslate.X), int32(sizePreTranslate.Y)}
	index := 0
	for y := 0; y < sizePreTranslate.Y; y++ {
		for x := 0; x < sizePreTranslate.X; x++ {
			pixels[index] = fyConvertGoImageColour(img.At(x, y))
			index++
		}
	}
	return GlobalBackend.CreateTexture(size, pixels)
}

func CreateHardcodedPNGImage(pngb64 string) image.Image {
	bytes, err := base64.StdEncoding.DecodeString(pngb64)
	decoded, err := png.Decode(strings.NewReader(string(bytes)))
	if err != nil {
		// Hard-coded so should always work
		panic(err)
	}
	return decoded
}
func CreateHardcodedPNGTexture(pngb64 string) Texture {
	return GoImageToTexture(CreateHardcodedPNGImage(pngb64))
}
