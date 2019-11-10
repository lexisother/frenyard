package frenyard

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"golang.org/x/image/math/fixed"
)

// FontRectangleConverter converts a fixed.Rectangle26_6 into pixels (rounding outwards)
func FontRectangleConverter(bounds fixed.Rectangle26_6) Area2i {
	area := bounds.Max.Sub(bounds.Min)
	return Area2i{
		Area1i{int32(bounds.Min.X.Floor()), int32(area.X.Ceil())},
		Area1i{int32(bounds.Min.Y.Floor()), int32(area.Y.Ceil())},
	}
}

// FontInterline gets the amount of pixels between lines of the font.
func FontInterline(font font.Face) int32 {
	return int32(font.Metrics().Height.Ceil())
}

// FontDraw draws a font image into a texture.
func FontDraw(fnt font.Face, text string) Texture {
	drawer := font.Drawer{
		Face: fnt,
	}
	bounds, _ := drawer.BoundString(text)
	area := FontRectangleConverter(bounds)
	img := image.NewNRGBA(image.Rect(0, 0, int(area.X.Size), int(area.Y.Size)))
	drawer.Dot = fixed.Point26_6{}.Sub(bounds.Min)
	drawer.Src = image.White
	drawer.Dst = img
	drawer.DrawString(text)
	return GoImageToTexture(img)
}

// FontSize gets the size of some text in the given font.
func FontSize(fnt font.Face, text string) Vec2i {
	drawer := font.Drawer{
		Face: fnt,
	}
	bounds, _ := drawer.BoundString(text)
	return FontRectangleConverter(bounds).Size()
}

// CreateTTFFont is a wrapper around truetype.NewFace
func CreateTTFFont(ft *truetype.Font, dpi float64, size float64) font.Face {
	return truetype.NewFace(ft, &truetype.Options{
		Size: size,
		DPI: dpi,
		Hinting: font.HintingNone,
	})
}
