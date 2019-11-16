package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

var borderButtonTexture frenyard.Texture
var borderHeaderTexture frenyard.Texture

func init() {
	borderButtonTexture = frenyard.CreateHardcodedPNGTexture(borderButtonX1B64)
	borderHeaderTexture = frenyard.CreateHardcodedPNGTexture(borderHeaderX2B64)
}

// BorderButton creates a border for a button of a given background colour.
func BorderButton(colour uint32) frenyard.NinePatchPackage {
	addedBorderX := int32(12)
	addedBorderY := int32(8)
	return frenyard.NinePatchPackage{
		Under: frenyard.NinePatch{
			Tex:       borderButtonTexture,
			Sprite:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 21, Y: 21}),
			Bounds:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 6, Y: 5}, frenyard.Vec2i{X: 9, Y: 9}),
			Centre:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 10, Y: 9}, frenyard.Vec2i{X: 1, Y: 1}),
			ColourMod: colour,
		},
		Over: frenyard.NinePatch{
			Tex:       borderButtonTexture,
			Sprite:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 6, Y: 5}, frenyard.Vec2i{X: 9, Y: 9}),
			Bounds:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 6, Y: 5}, frenyard.Vec2i{X: 9, Y: 9}),
			Centre:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 10, Y: 9}, frenyard.Vec2i{X: 1, Y: 1}),
			ColourMod: colour,
		},
		Padding:  frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2}),
		Scale: DesignScale,
		Clipping: true,
	}
}

// This border deliberately "hangs over" on the OVER layer.
// With correct manipulation of Z (ensure title is LAST in flexbox & use Order to correct)
// this produces a shadow effect that even hangs over other UI.

// BorderTitle produces a border for the shadowing effect under a title.
func BorderTitle(colour uint32) frenyard.NinePatchPackage {
	addedBorderX := int32(8)
	addedBorderY := int32(8)
	return frenyard.NinePatchPackage{
		Over: frenyard.NinePatch{
			Tex:       borderHeaderTexture,
			Sprite:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 1, Y: 16}),
			Bounds:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 1, Y: 2}),
			Centre:    frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 1, Y: 2}),
			ColourMod: colour,
		},
		Padding:  frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2}),
		Scale: DesignScale / 2,
		Clipping: true,
	}
}
