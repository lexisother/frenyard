package integration

import (
	"image/color"
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

// ColourTransform transforms colours. It's useful when doing image transforms from alphaless images to alpha/etc.
type ColourTransform func(c uint32) uint32

// ColourTransform2 combines two colours to make a third.
type ColourTransform2 func(a uint32, b uint32) uint32

// ColourTransformBlueToStencil uses the blue channel as an alpha stencil and gives the image a white background.
func ColourTransformBlueToStencil(c uint32) uint32 {
	return ((c & 0xFF) << 24) | 0xFFFFFF
}

// ColourTransformInvert inverts the colour channels.
func ColourTransformInvert(c uint32) uint32 {
	a, r, g, b := ColourToARGB(c)
	return ColourFromARGB(a, 255-r, 255-g, 255-b)
}

// ColourTransform2Blend implements standard blending.
func ColourTransform2Blend(dest uint32, source uint32) uint32 {
	dA := dest >> 24
	sA := source >> 24
	resColour := ColourMix(source, dest, float64(dA)/255) & 0xFFFFFF
	resAlpha := ColourComponentClamp(int32(sA + dA))
	return resColour | uint32(resAlpha)<<24
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
