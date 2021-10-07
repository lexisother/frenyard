package integration

import (
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

type fyUnderlineTypeChunk struct {
	interior TypeChunk
	colour   uint32
}

// NewUnderlineTypeChunk underlines a TypeChunk.
func NewUnderlineTypeChunk(text TypeChunk, underlineColour uint32) TypeChunk {
	return fyUnderlineTypeChunk{text, underlineColour}
}

// FyCComponentCount implements TypeChunk.FyCComponentCount
func (ul fyUnderlineTypeChunk) FyCComponentCount() int {
	return ul.interior.FyCComponentCount()
}

// FyCComponentBreakStatus implements TypeChunk.FyCComponentBreakStatus
func (ul fyUnderlineTypeChunk) FyCComponentBreakStatus(index int) TypeChunkComponentBreakStatus {
	return ul.interior.FyCComponentBreakStatus(index)
}

// FyCComponentAdvance implements TypeChunk.FyCComponentAdvance
func (ul fyUnderlineTypeChunk) FyCComponentAdvance(index int, kerning bool) fixed.Int26_6 {
	return ul.interior.FyCComponentAdvance(index, kerning)
}

// FyCSection implements TypeChunk.FyCSection
func (ul fyUnderlineTypeChunk) FyCSection(start int, end int) TypeChunk {
	return fyUnderlineTypeChunk{ul.interior.FyCSection(start, end), ul.colour}
}

// FyCHeight implements TypeChunk.FyCHeight
func (ul fyUnderlineTypeChunk) FyCHeight() int {
	return ul.interior.FyCHeight()
}

// FyCBounds implements TypeChunk.FyCBounds
func (ul fyUnderlineTypeChunk) FyCBounds(dot fixed.Point26_6) (fixed.Point26_6, fixed.Rectangle26_6) {
	return ul.interior.FyCBounds(dot)
}

// FyCDraw implements TypeChunk.FyCDraw
func (ul fyUnderlineTypeChunk) FyCDraw(img draw.Image, dot fixed.Point26_6) fixed.Point26_6 {
	result := ul.interior.FyCDraw(img, dot)
	line := dot.Y.Ceil()
	lineHeight := 1
	draw.Draw(img, image.Rect(dot.X.Floor(), line, result.X.Ceil(), line+lineHeight), image.NewUniform(ConvertUint32ToGoImageColour(ul.colour)), image.Point{}, draw.Over)
	return result
}
