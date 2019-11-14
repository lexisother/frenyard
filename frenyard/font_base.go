package frenyard

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

// TypeChunkComponentBreakStatus represents the break status of a component.
type TypeChunkComponentBreakStatus uint8

// TypeChunkComponentBreakStatusNone represents no break at all.
const TypeChunkComponentBreakStatusNone TypeChunkComponentBreakStatus = 0
// TypeChunkComponentBreakStatusSpace is an optional break
const TypeChunkComponentBreakStatusSpace TypeChunkComponentBreakStatus = 1
// TypeChunkComponentBreakStatusNewline is an obligatory break, real component may act erroneously
const TypeChunkComponentBreakStatusNewline TypeChunkComponentBreakStatus = 2

// TypeChunk represents a chunk of abstract text-like stuff.
type TypeChunk interface {
	// Amount of components in the type.
	FyCComponentCount() int
	// Break status of a specific component.
	FyCComponentBreakStatus(index int) TypeChunkComponentBreakStatus
	// Given a specific component, returns the advance. Can account for kerning.
	FyCComponentAdvance(index int, kerning bool) fixed.Int26_6
	// Gets a subsection. 'end' is non-inclusive; thus Section(0, ComponentCount()) returns a chunk equivalent to the TypeChunk itself.
	FyCSection(start int, end int) TypeChunk
	// Line height (should be maximum of total line height in the chunk for evenness). As pixels to keep things nice.
	FyCHeight() int
	// Returns the dot, and then the bounds. These bounds can and should be increased where possible to try and reduce vertical 'bumps'.
	FyCBounds(dot fixed.Point26_6) (fixed.Point26_6, fixed.Rectangle26_6)
	// Draws to an image. Returns the new dot.
	FyCDraw(img draw.Image, dot fixed.Point26_6) fixed.Point26_6
}

type fyTextTypeChunk struct {
	_text []rune
	_face font.Face
}

// NewTextTypeChunk creates a TypeChunk from text and a font. This is the simplest kind of TypeChunk.
func NewTextTypeChunk(text string, face font.Face) TypeChunk {
	runes := []rune{}
	for _, r := range text {
		runes = append(runes, r)
	}
	return fyTextTypeChunk{
		runes,
		face,
	}
}

// FyCComponentCount implements TypeChunk.FyCComponentCount
func (ttc fyTextTypeChunk) FyCComponentCount() int {
	return len(ttc._text)
}

// FyCComponentBreakStatus implements TypeChunk.FyCComponentBreakStatus
func (ttc fyTextTypeChunk) FyCComponentBreakStatus(index int) TypeChunkComponentBreakStatus {
	if ttc._text[index] == 10 {
		return TypeChunkComponentBreakStatusNewline
	} else if ttc._text[index] == ' ' {
		return TypeChunkComponentBreakStatusSpace
	}
	return TypeChunkComponentBreakStatusNone
}

// FyCComponentAdvance implements TypeChunk.FyCComponentAdvance
func (ttc fyTextTypeChunk) FyCComponentAdvance(index int, kerning bool) fixed.Int26_6 {
	adv, _ := ttc._face.GlyphAdvance(ttc._text[index])
	if index > 0 && kerning {
		adv += ttc._face.Kern(ttc._text[index - 1], ttc._text[index])
	}
	return adv
}

// FyCSection implements TypeChunk.FyCSection
func (ttc fyTextTypeChunk) FyCSection(start int, end int) TypeChunk {
	return fyTextTypeChunk{
		ttc._text[start:end],
		ttc._face,
	}
}

// FyCHeight implements TypeChunk.FyCHeight
func (ttc fyTextTypeChunk) FyCHeight() int {
	return ttc._face.Metrics().Height.Ceil()
}

// FyCBounds implements TypeChunk.FyCBounds
func (ttc fyTextTypeChunk) FyCBounds(dot fixed.Point26_6) (fixed.Point26_6, fixed.Rectangle26_6) {
	drawer := font.Drawer{
		Face: ttc._face,
		Dot: dot,
	}
	text := string(ttc._text)
	bounds, advance := drawer.BoundString(text)
	boundSize := bounds.Max.Sub(bounds.Min)
	// Adjust the bounds to try and prevent trouble
	metricHeight := ttc._face.Metrics().Height
	if boundSize.Y < metricHeight {
		boundSize.Y = metricHeight
	}
	// Apply adjustment
	bounds.Max = bounds.Min.Add(boundSize)
	// Apply advance
	dot.X += advance
	return dot, bounds
}

// FyCDraw implements TypeChunk.FyCDraw
func (ttc fyTextTypeChunk) FyCDraw(img draw.Image, dot fixed.Point26_6) fixed.Point26_6 {
	drawer := font.Drawer{
		Face: ttc._face,
		Dot: dot,
		Src: image.White,
		Dst: img,
	}
	drawer.DrawString(string(ttc._text))
	return drawer.Dot
}

type fyCompoundTypeChunk struct {
	_chunks []TypeChunk
}

// NewCompoundTypeChunk creates a type chunk from multiple sub-type-chunks.
func NewCompoundTypeChunk(content []TypeChunk) TypeChunk {
	return fyCompoundTypeChunk{
		content,
	}
}

// FyCComponentCount implements TypeChunk.FyCComponentCount
func (ctc fyCompoundTypeChunk) FyCComponentCount() int {
	total := 0
	for _, v := range ctc._chunks {
		total += v.FyCComponentCount()
	}
	return total
}

// FyCComponentBreakStatus implements TypeChunk.FyCComponentBreakStatus
func (ctc fyCompoundTypeChunk) FyCComponentBreakStatus(index int) TypeChunkComponentBreakStatus {
	for _, v := range ctc._chunks {
		l := v.FyCComponentCount()
		if index < l {
			return v.FyCComponentBreakStatus(index)
		}
		index -= l
	}
	return TypeChunkComponentBreakStatusNone
}

// FyCComponentAdvance implements TypeChunk.FyCComponentAdvance
func (ctc fyCompoundTypeChunk) FyCComponentAdvance(index int, kerning bool) fixed.Int26_6 {
	for _, v := range ctc._chunks {
		l := v.FyCComponentCount()
		if index < l {
			return v.FyCComponentAdvance(index, kerning)
		}
		index -= l
	}
	return 0
}

// FyCSection implements TypeChunk.FyCSection
func (ctc fyCompoundTypeChunk) FyCSection(start int, end int) TypeChunk {
	sections := []TypeChunk{}
	for _, v := range ctc._chunks {
		if end <= 0 {
			break
		}
		count := v.FyCComponentCount()
		if start < count {
			resultEnd := end
			if resultEnd > count {
				resultEnd = count
			}
			sections = append(sections, v.FyCSection(start, resultEnd))
		}
		start -= count
		if start < 0 {
			start = 0
		}
		end -= count
	}
	return fyCompoundTypeChunk{
		sections,
	}
}

// FyCHeight implements TypeChunk.FyCHeight
func (ctc fyCompoundTypeChunk) FyCHeight() int {
	h := 0
	for _, v := range ctc._chunks {
		vh := v.FyCHeight()
		if vh > h {
			h = vh
		}
	}
	return h
}

// FyCBounds implements TypeChunk.FyCBounds
func (ctc fyCompoundTypeChunk) FyCBounds(dot fixed.Point26_6) (fixed.Point26_6, fixed.Rectangle26_6) {
	totalBounds := fixed.Rectangle26_6{}
	for _, v := range ctc._chunks {
		newDot, bounds := v.FyCBounds(dot)
		if !totalBounds.Empty() {
			totalBounds = totalBounds.Union(bounds)
		} else {
			totalBounds = bounds
		}
		dot = newDot
	}
	return dot, totalBounds
}

// FyCDraw implements TypeChunk.FyCDraw
func (ctc fyCompoundTypeChunk) FyCDraw(img draw.Image, dot fixed.Point26_6) fixed.Point26_6 {
	for _, v := range ctc._chunks {
		dot = v.FyCDraw(img, dot)
	}
	return dot
}



// FontRectangleConverter converts a fixed.Rectangle26_6 into pixels (rounding outwards)
func FontRectangleConverter(bounds fixed.Rectangle26_6) Area2i {
	area := bounds.Max.Sub(bounds.Min)
	return Area2i{
		Area1i{int32(bounds.Min.X.Floor()), int32(area.X.Ceil())},
		Area1i{int32(bounds.Min.Y.Floor()), int32(area.Y.Ceil())},
	}
}

// DPIPixels is the DPI setting where one typographical point is equal to one typographical pixel.
const DPIPixels float64 = 72

// CreateTTFFont is a wrapper around truetype.NewFace
func CreateTTFFont(ft *truetype.Font, dpi float64, size float64) font.Face {
	return truetype.NewFace(ft, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
}
