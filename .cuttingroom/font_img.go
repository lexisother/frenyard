package frenyard

import "image"
import "fmt"

/*
 * Please note that, should this distinction matter, this attempts to implement
 *  the rules required to display hall-fetica.png of CrossCode.
 * This has been done via inference and a heavy amount of overabstraction to avoid any code copyright concerns.
 * Should concerns still exist, please contact 20kdc.
 */

/*
 * A character in the bitmap font.
 * The default value of this must be a valid 'no-op' char.
 */
type fyBitmapFontCharacter struct {
	sprite Area2i
	move int32
}
func newBitmapFontCharacterNormal(sprite Area2i, postSpacing int32) fyBitmapFontCharacter {
	return fyBitmapFontCharacter{
		sprite: sprite,
		move: sprite.X.Size + postSpacing,
	}
}

type fyBitmapFontPlacement struct {
	position Vec2i
	sprite Area2i
}

/*
 * The bitmap font itself.
 */
type fyBitmapFont struct {
	image image.Image
	interLine int32
	chars map[rune]fyBitmapFontCharacter
	trustAlpha bool
}

func (f *fyBitmapFont) Interline() int32 {
	return f.interLine
}

func (f *fyBitmapFont) layout(text string) []fyBitmapFontPlacement {
	placements := make([]fyBitmapFontPlacement, len(text))
	runes := 0
	x := int32(0)
	for _, char := range text {
		exec := f.chars[char]
		placements[runes] = fyBitmapFontPlacement{Vec2i{x, 0}, exec.sprite}
		x += exec.move
		runes++
	}
	return placements[:runes]
}

func (f *fyBitmapFont) Draw(text string) Texture {
	placements := f.layout(text)
	size := Vec2i{}
	for _, placement := range f.layout(text) {
		size = size.Max(placement.position.Add(placement.sprite.Size()))
	}
	img := image.NewNRGBA(image.Rect(0, 0, int(size.X), int(size.Y)))
	for _, placement := range placements {
		for y := int32(0); y < placement.sprite.Y.Size; y++ {
			for x := int32(0); x < placement.sprite.X.Size; x++ {
				src := f.image.At(int(placement.sprite.X.Pos + x), int(placement.sprite.Y.Pos + y))
				r, _, _, _ := src.RGBA()
				if r != 0 || f.trustAlpha {
					img.Set(int(placement.position.X + x), int(placement.position.Y + y), src)
				}
			}
		}
	}
	return GoImageToTexture(img)
}
func (f *fyBitmapFont) Size(text string) Vec2i {
	size := Vec2i{}
	for _, placement := range f.layout(text) {
		size = size.Max(placement.position.Add(placement.sprite.Size()))
	}
	return size
}

func CreateImpactlikeFont(img image.Image) (Font, error) {
	rect := img.Bounds()
	font := fyBitmapFont{}
	font.image = img
	font.trustAlpha = true
	font.chars = map[rune]fyBitmapFontCharacter{}
	lineHeight := 1
	for {
		_, _, _, a := img.At(rect.Min.X, rect.Min.Y + lineHeight - 1).RGBA()
		if a != 0 {
			break
		}
		lineHeight += 1
		if lineHeight > rect.Dy() {
			return nil, fmt.Errorf("Unable to find starting line of font.")
		}
	}
	currentChar := ' '
	for charBase := 0; charBase + lineHeight <= rect.Dy(); charBase += lineHeight {
		// This is one line of characters. For out-of-bounds values, the image reads return zero, which is fine.
		// This behavior is relied upon to catch the last pixel in a non-awkward way.
		x := 0
		charStart := -1
		for {
			// Handle pixel
			_, _, _, a := img.At(rect.Min.X + x, rect.Min.Y + charBase + lineHeight - 1).RGBA()
			if a != 0 {
				if charStart == -1 {
					charStart = x
				}
			} else {
				if charStart != -1 {
					// Keep in mind charStart is the first on pixel and x is the first off pixel.
					// So no need for a manual +/- 1
					font.chars[currentChar] = newBitmapFontCharacterNormal(Area2iFromVecs(Vec2i{int32(charStart), int32(charBase)}, Vec2i{int32(x - charStart), int32(lineHeight) - 1}), 1)
					currentChar += 1
				}
				charStart = -1
			}
			// Now check if out of range (after a pixel out of range has already been checked; see previous note)
			if x >= rect.Dx() {
				break
			}
			x += 1
		}
	}
	return &font, nil
}

// Useful for basic 'ASCII sheet' or 'codepage 437' fonts.
// Never fails.
func CreateASCIISheetFont(img image.Image, size Vec2i, spacing Vec2i) Font {
	rect := img.Bounds()
	font := fyBitmapFont{}
	font.image = img
	font.interLine = size.Y + spacing.Y
	font.chars = map[rune]fyBitmapFontCharacter{}
	cRune := rune(0)
	for y := 0; y < rect.Dy(); y += int(size.Y) {
		for x := 0; x < rect.Dx(); x += int(size.X) {
			font.chars[cRune] = newBitmapFontCharacterNormal(Area2iFromVecs(Vec2i{int32(x), int32(y)}, size), spacing.X)
			cRune++
		}
	}
	return &font
}
