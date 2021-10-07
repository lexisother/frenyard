package integration

import (
	"encoding/base64"
	"github.com/yellowsink/frenyard"
	"image"
	"image/png"
	"strings"
)

// GoImageToTexture imports an image from Go's "image" library to a texture.
func GoImageToTexture(img image.Image, ct []ColourTransform) frenyard.Texture {
	min := img.Bounds().Min
	sizePreTranslate := img.Bounds().Size()
	pixels := make([]uint32, sizePreTranslate.X*sizePreTranslate.Y)
	size := frenyard.Vec2i{X: int32(sizePreTranslate.X), Y: int32(sizePreTranslate.Y)}
	index := 0
	for y := 0; y < sizePreTranslate.Y; y++ {
		for x := 0; x < sizePreTranslate.X; x++ {
			res := ConvertGoImageColourToUint32(img.At(x+min.X, y+min.Y))
			for _, v := range ct {
				res = v(res)
			}
			pixels[index] = res
			index++
		}
	}
	return frenyard.GlobalBackend.CreateTexture(size, pixels)
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
	destBounds := image.Rect(sourceBounds.Min.X/2, sourceBounds.Min.Y/2, sourceBounds.Max.X/2, sourceBounds.Max.Y/2)
	dest := image.NewNRGBA(destBounds)
	for y := sourceBounds.Min.Y; y < sourceBounds.Max.Y; y += 2 {
		for x := sourceBounds.Min.X; x < sourceBounds.Max.X; x += 2 {
			a := ConvertGoImageColourToUint32(source.At(x, y))
			b := ConvertGoImageColourToUint32(source.At(x, y+1))
			c := ConvertGoImageColourToUint32(source.At(x+1, y))
			d := ConvertGoImageColourToUint32(source.At(x+1, y+1))
			result := ColourMix(ColourMix(a, b, 0.5), ColourMix(c, d, 0.5), 0.5)
			dest.SetNRGBA(x/2, y/2, ConvertUint32ToGoImageColour(result))
		}
	}
	return dest
}

// CreateHardcodedPNGTexture gets a main.Texture from a base64 string.
func CreateHardcodedPNGTexture(pngb64 string, ct []ColourTransform) frenyard.Texture {
	return GoImageToTexture(CreateHardcodedPNGImage(pngb64), ct)
}

// Run runs the colour transform on a pair of images to create a third image.
func (ct ColourTransform2) Run(part1 image.Image, part2 image.Image, shift frenyard.Vec2i) image.Image {
	sourceBounds := part1.Bounds()
	part2Bounds := part2.Bounds()
	dest := image.NewNRGBA(sourceBounds)
	for y := sourceBounds.Min.Y; y < sourceBounds.Max.Y; y++ {
		for x := sourceBounds.Min.X; x < sourceBounds.Max.X; x++ {
			x2 := (x - sourceBounds.Min.X) + part2Bounds.Min.X + int(shift.X)
			y2 := (y - sourceBounds.Min.Y) + part2Bounds.Min.Y + int(shift.Y)
			a := ConvertGoImageColourToUint32(part1.At(x, y))
			b := ConvertGoImageColourToUint32(part2.At(x2, y2))
			dest.SetNRGBA(x, y, ConvertUint32ToGoImageColour(ct(a, b)))
		}
	}
	return dest
}
