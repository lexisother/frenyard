package frenyard

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
)

// ColourFromARGB creates a colour from the separate A/R/G/B quantities.
func ColourFromARGB(a uint8, r uint8, g uint8, b uint8) uint32 {
	return (uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | (uint32(b) << 0)
}

// ColourToARGB splits a colour apart.
func ColourToARGB(col uint32) (a uint8, r uint8, g uint8, b uint8) {
	a = uint8((col & 0xFF000000) >> 24)
	r = uint8((col & 0xFF0000) >> 16)
	g = uint8((col & 0xFF00) >> 8)
	b = uint8((col & 0xFF) >> 0)
	return a, r, g, b
}

// ColourComponentClamp clamps an integer into the 0..255 range (acting as a safe converter from int32 to uint8)
func ColourComponentClamp(i int32) uint8 {
	if i < 0 {
		return 0
	} else if i > 255 {
		return 255
	}
	return uint8(i)
}

// ColourMix linearly interpolates between two colours in the standard space (no "linear light" business)
func ColourMix(cola uint32, colb uint32, amount float64) uint32 {
	invAmount := 1 - amount
	aA, aR, aG, aB := ColourToARGB(cola)
	bA, bR, bG, bB := ColourToARGB(colb)
	cA := ColourComponentClamp(int32((float64(aA) * invAmount) + (float64(bA) * amount)))
	cR := ColourComponentClamp(int32((float64(aR) * invAmount) + (float64(bR) * amount)))
	cG := ColourComponentClamp(int32((float64(aG) * invAmount) + (float64(bG) * amount)))
	cB := ColourComponentClamp(int32((float64(aB) * invAmount) + (float64(bB) * amount)))
	return ColourFromARGB(cA, cR, cG, cB)
}

// A ColourTransform transforms colours. It's useful when doing image transforms from alphaless images to alpha/etc.
type ColourTransform func(c uint32) uint32

// ColourTransformBlueToStencil uses the blue channel as an alpha stencil and gives the image a white background.
func ColourTransformBlueToStencil(c uint32) uint32 {
	return ((c & 0xFF) << 24) | 0xFFFFFF
}

// ColourTransformInvert inverts the colour channels.
func ColourTransformInvert(c uint32) uint32 {
	a, r, g, b := ColourToARGB(c)
	return ColourFromARGB(a, 255 - r, 255 - g, 255 - b)
}

// ConvertGoImageColourToUint32 converts a Go colour into an ARGB colour uint32 as losslessly as possible.
func ConvertGoImageColourToUint32(col color.Color) uint32 {
	// Colour type appears to be per-pixel. Great. More problems.
	switch colConv := col.(type) {
	case color.NRGBA:
		return ColourFromARGB(colConv.A, colConv.R, colConv.G, colConv.B)
	case color.RGBA:
		// This is still premultiplied and thus NOT GOOD, but let's roll with it
		a := colConv.A
		r := colConv.R
		g := colConv.G
		b := colConv.B
		if (a != 0) && (a != 255) {
			// Alpha not 0 (would /0) or a NOP
			a16 := uint16(a)
			r = uint8((uint16(colConv.R) * 255) / a16)
			g = uint8((uint16(colConv.G) * 255) / a16)
			b = uint8((uint16(colConv.B) * 255) / a16)
		}
		return ColourFromARGB(a, r, g, b)
	}
	// Give up and work backwards from premultiplied (NOT GOOD!!!)
	r, g, b, a := col.RGBA()
	if (a != 0) && (a != 0xFFFF) {
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

// ConvertUint32ToGoImageColour converts a uint32 colour to a colour.NRGBA (which is essentially in the same format)
func ConvertUint32ToGoImageColour(col uint32) color.NRGBA {
	a, r, g, b := ColourToARGB(col)
	return color.NRGBA{
		A: a,
		R: r,
		G: g,
		B: b,
	}
}

// GoImageToTexture imports an image from Go's "image" library to a texture.
func GoImageToTexture(img image.Image, ct []ColourTransform) Texture {
	sizePreTranslate := img.Bounds().Size()
	pixels := make([]uint32, sizePreTranslate.X*sizePreTranslate.Y)
	size := Vec2i{int32(sizePreTranslate.X), int32(sizePreTranslate.Y)}
	index := 0
	for y := 0; y < sizePreTranslate.Y; y++ {
		for x := 0; x < sizePreTranslate.X; x++ {
			res := ConvertGoImageColourToUint32(img.At(x, y))
			for _, v := range ct {
				res = v(res)
			}
			pixels[index] = res
			index++
		}
	}
	return GlobalBackend.CreateTexture(size, pixels)
}

// CreateHardcodedPNGImage gets an image.Image from a base64 string.
func CreateHardcodedPNGImage(pngb64 string) image.Image {
	bytes, err := base64.StdEncoding.DecodeString(pngb64)
	if err != nil {
		// Hard-coded so should always work
		panic(err)
	}
	decoded, err := png.Decode(strings.NewReader(string(bytes)))
	if err != nil {
		// Hard-coded so should always work
		panic(err)
	}
	return decoded
}

// ScaleImageToHalfSize scales an image to half-size using a trivial "average covered pixels" algorithm that handles a lot of situations well assuming content that's aligned to the implicit "2x2 grid".
func ScaleImageToHalfSize(source image.Image) image.Image {
	sourceBounds := source.Bounds()
	destBounds := image.Rect(sourceBounds.Min.X / 2, sourceBounds.Min.Y / 2, sourceBounds.Max.X / 2, sourceBounds.Max.Y / 2)
	dest := image.NewNRGBA(destBounds)
	for y := destBounds.Min.Y; y < destBounds.Max.Y; y++ {
		for x := destBounds.Min.X; x < destBounds.Max.X; x++ {
			a := ConvertGoImageColourToUint32(source.At(x * 2, y * 2))
			b := ConvertGoImageColourToUint32(source.At(x * 2, (y * 2) + 1))
			c := ConvertGoImageColourToUint32(source.At((x * 2) + 1, y * 2))
			d := ConvertGoImageColourToUint32(source.At((x * 2) + 1, (y * 2) + 1))
			result := ColourMix(ColourMix(a, b, 0.5), ColourMix(c, d, 0.5), 0.5)
			dest.SetNRGBA(x, y, ConvertUint32ToGoImageColour(result))
		}
	}
	return dest
}

// CreateHardcodedPNGTexture gets a frenyard.Texture from a base64 string.
func CreateHardcodedPNGTexture(pngb64 string, ct []ColourTransform) Texture {
	return GoImageToTexture(CreateHardcodedPNGImage(pngb64), ct)
}
