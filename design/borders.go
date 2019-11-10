package design
import "github.com/20kdc/CCUpdaterUI/frenyard"

const borderButtonImage = "iVBORw0KGgoAAAANSUhEUgAAABUAAAAVCAYAAACpF6WWAAAABmJLR0QAYgAAAO5dA1AtAAAACXBIWXMAAA+IAAAPiAEWyKWGAAAAB3RJTUUH4wsHFgcJmDU7hQAAABl0RVh0Q29tbWVudABDcmVhdGVkIHdpdGggR0lNUFeBDhcAAAFtSURBVDjLtZW9bsIwFEaPk4ggVFWVOnSreCFm3oIH6ivkJdgRM1uHTkR1F1KgLobboS4yru0wtJasOP45ubrf5xslIsSaUkrR0yRxWIXz18D64GdoAqaCsX9YUnAlIjGg8p4qApVgfAGuEtH5vQgiFuCUAxcZYKm1nlhr5yKyFpG9iKyttXOt9QQog/1JUAFUQK21nkqmaa2nQO32F+EHfGAJDICRMWaRgxpjFsDI7S99cCjEOVIReQFuM07aKKUeAQNYL88XOT3nEhj2AHHrwyC3UaF8u3Q90C5wQFT9C4t0XfecJX6vx3ybjPS0XC6fclC3fopFGlO/Bu6AcdM0s7ZtV8aYrVN827btqmmaGTB2++qc+nhCDYAbJ8YDcO9E+QDegBbYAO/AJ3D0I65ERIK7L8AB2Ll3C7w6q1kH2bl+uObu/9xrnP+OwN6zjbg563nzV16rlFDe00aqlF9QJFmk/7ye/lvl/8t/1BeHUwxkWeldxwAAAABJRU5ErkJggg=="

const borderHeaderImage = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAAICAYAAAA4GpVBAAAC/npUWHRSYXcgcHJvZmlsZSB0eXBlIGV4aWYAAHjazZVbcuUoDIbfWcUsAUkIieVgLlWzg1n+/GB8cpJOdXXX9MOYBLAs6/YhnzD++XuGv3ARcQxJzXPJOeJKJRWu2Hi8r7pnimnP8dw96yd5eL3EEAlWuW89H/mAnKHPR96OnQq5vhkq4zy4Pj+oxxD7cfBEchwJ3Q7iMRzqMSR8PKf7/rrTirm4vafQjv48z3dGflJLYpw1kyXMiaNZLtg7x2SoW1+BDr5joXIy/nIfHlXImYeQRMy+IhSEL0UqZsIcBYaxM+zTlrPwHSlocYj25iX+/PpZ5OEJ/dD7hJTj96hfuzfSYRf0eSBfCOXX+q2c9EMe3pFubm+ec355/iSf6SPYh1x48M3Zfc5xZ1dTRsr5JPWksnfQu1a19lsZw/Cv0Vex1ygYjpZoOEcdh+zCaFSIgXFSok6VJo29NmoIMfFgw8rcWKgFCB0wCreNO61Bkw2AuzgQNxwHgZRfsdB2W7a7Rg7HnTzg9BCM0ToV/2WEX1GaczUUUfS7Trz58io4oljkiEIkqIEIzVNU3QV+xtdrcRUQ1F1mR4I1XssC8F9KH4dLNmiBomK9G5isHwMoESJQBIO2SRQziVKmaMzBiFBIB6CK0FkSX8BCqtwRJCeRDDjoAvjGO0ZblZVvMT6EAKGSA5rQVwsCVkqK82PJcYaqiiZVzWrqWrRmyavDcra8vqjVxJKpZTNzK1aDiydXz27uXrwWLoIvrhb0Y/FSSq1wWmG54u0KhVovvuRKl175ssuvctXGoUlLTVtu1ryVVjt36ejjnrt176XXQQNHaaShIw8bPsqoE0dtykxTZ542fZYw64sanbb9On6DGh1qvEktRXtRg9TsMUHrc6KLGYhxIgC3RYAkMC9m0SklXuQWs1gYXaGMIHXB6bSIgWAaxDrpxe6DnAbJf4ZbAAj+E+TCQvcL5H7k9h21XvcPnWxCqw1XUaOg+6BU2fGHX8zvVxn3NvxE57fW/6chQZVK+Be6sf05X4iU+wAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAAPiAAAD4gBFsilhgAAAAd0SU1FB+MLCBExFsENodYAAAAZdEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBXgQ4XAAAAI0lEQVQI1wXBQQEAIAgEsAkP/tYwkP2THJskgQ8PLgw0nIIFppQEnoQeP6QAAAAASUVORK5CYII="

var borderButtonTexture frenyard.Texture
var borderHeaderTexture frenyard.Texture

func init() {
	borderButtonTexture = frenyard.CreateHardcodedPNGTexture(borderButtonImage)
	borderHeaderTexture = frenyard.CreateHardcodedPNGTexture(borderHeaderImage)
}

// BorderButton creates a border for a button of a given background colour.
func BorderButton(colour uint32) frenyard.NinePatchPackage {
	addedBorderX := int32(12)
	addedBorderY := int32(8)
	return frenyard.NinePatchPackage{
		Under: frenyard.NinePatch{
			Tex: borderButtonTexture,
			Sprite: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 21, Y: 21}),
			Bounds: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 6, Y: 5}, frenyard.Vec2i{X: 9, Y: 9}),
			Centre: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 10, Y: 9}, frenyard.Vec2i{X: 1, Y: 1}),
			ColourMod: colour,
		},
		Over: frenyard.NinePatch{
			Tex: borderButtonTexture,
			Sprite: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 6, Y: 5}, frenyard.Vec2i{X: 9, Y: 9}),
			Bounds: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 6, Y: 5}, frenyard.Vec2i{X: 9, Y: 9}),
			Centre: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 10, Y: 9}, frenyard.Vec2i{X: 1, Y: 1}),
			ColourMod: colour,
		},
		Padding: frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2}),
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
			Tex: borderHeaderTexture,
			Sprite: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 1, Y: 8}),
			Bounds: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 1, Y: 1}),
			Centre: frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 1, Y: 1}),
			ColourMod: colour,
		},
		Padding: frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2}),
		Clipping: true,
	}
}
