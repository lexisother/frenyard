package frenyard

// Constants

const SIZE_UNLIMITED int32 = 0x7FFFFFFF
func Vec2iUnlimited() Vec2i { return Vec2i{SIZE_UNLIMITED, SIZE_UNLIMITED} }

// Part I: Basic maths

func Max(a int32, b int32) int32 {
	if (a > b) {
		return a
	} else {
		return b
	}
}
func Min(a int32, b int32) int32 {
	if (a < b) {
		return a
	} else {
		return b
	}
}

// Part II: Point Type

/*
 * The basic 2-dimensional vector type.
 */
type Vec2i struct {
	X, Y int32;
}

func (a Vec2i) Add(b Vec2i) Vec2i {
	return Vec2i{a.X + b.X, a.Y + b.Y}
}

func (a Vec2i) Min(b Vec2i) Vec2i {
	return Vec2i{Min(a.X, b.X), Min(a.Y, b.Y)}
}
func (a Vec2i) Max(b Vec2i) Vec2i {
	return Vec2i{Max(a.X, b.X), Max(a.Y, b.Y)}
}

func (a Vec2i) Eq(b Vec2i) bool {
	return a.X == b.X && a.Y == b.Y
}
// These two are framed in the sense of 'an area of size A could contain an area of size B'.
func (a Vec2i) Gt(b Vec2i) bool {
	return a.X > b.X && a.Y > b.Y
}
func (a Vec2i) Ge(b Vec2i) bool {
	return a.X >= b.X && a.Y >= b.Y
}

func (a Vec2i) Negate() Vec2i {
	return Vec2i{-a.X, -a.Y}
}
// Conditionally swaps X/Y
func (a Vec2i) ConditionalTranspose(yes bool) Vec2i {
	if yes {
		return Vec2i{a.Y, a.X}
	}
	return Vec2i{a.X, a.Y}
}

// Part III: Area Types (AABBs)

/*
 * 1-dimensional area. Useful for intersection.
 */
type Area1i struct {
	Pos int32;
	Size int32;
}

func Area1iOfSize(a int32) Area1i { return Area1i{0, a} }
// Size <= 0
func (a Area1i) Empty() bool {
	return a.Size <= 0
}
// Replaces empty areas with zeroed areas.
func (a Area1i) Normalized() Area1i {
	if (a.Empty()) {
		return Area1i{}
	} else {
		return a
	}
}

// Unions two areas.
func (a Area1i) Union(b Area1i) Area1i {
	pos := Min(a.Pos, b.Pos)
	end := Max(a.Pos + a.Size, b.Pos + b.Size)
	return Area1i{
		pos,
		end - pos,
	}
}
// Intersects two areas. Always returns a normalized area.
func (a Area1i) Intersect(b Area1i) Area1i {
	pos := Max(a.Pos, b.Pos)
	end := Min(a.Pos + a.Size, b.Pos + b.Size)
	return Area1i{
		pos,
		end - pos,
	}.Normalized()
}
// Translates an area by an offset.
func (a Area1i) Translate(i int32) Area1i { return Area1i{a.Pos + i, a.Size} }
// Expands an area by a margin, expressed as the area around as if this had zero size.
func (a Area1i) Expand(n Area1i) Area1i { return Area1i{a.Pos + n.Pos, a.Size + n.Size} }
// Contracts an area by a margin; the reverse of Expand.
func (a Area1i) Contract(n Area1i) Area1i { return Area1i{a.Pos - n.Pos, a.Size - n.Size} }
// Checks if a point is within the area.
func (a Area1i) Contains(i int32) bool { return (i >= a.Pos) && (i < a.Pos + a.Size) }
// Aligns an area within another.
func (a Area1i) Align(content int32, x Alignment1i) Area1i {
	if x == AlignStart {
		return Area1i{
			a.Pos,
			content,
		}
	} else if x == AlignEnd {
		return Area1i{
			a.Size - content,
			content,
		}
	} else {
		return Area1i{
			((a.Size - content) / 2) + a.Pos,
			content,
		}
	}
}

func UnionArea1i(areas []Area1i) Area1i {
	base := Area1i{}
	for _, area := range areas {
		base = base.Union(area)
	}
	return base
}

/*
 * The basic rectangle type.
 */
type Area2i struct {
	X Area1i;
	Y Area1i;
}

func (a Area2i) Pos() Vec2i { return Vec2i{a.X.Pos, a.Y.Pos} }
func (a Area2i) Size() Vec2i { return Vec2i{a.X.Size, a.Y.Size} }

func Area2iFromVecs(pos Vec2i, size Vec2i) Area2i {
	return Area2i{
		Area1i{pos.X, size.X},
		Area1i{pos.Y, size.Y},
	}
}

func Area2iOfSize(a Vec2i) Area2i { return Area2i{Area1iOfSize(a.X), Area1iOfSize(a.Y)} }
// Size <= 0
func (r Area2i) Empty() bool {
	return r.X.Empty() || r.Y.Empty()
}
// Replaces empty areas with zeroed areas.
func (a Area2i) Normalized() Area2i {
	if (a.Empty()) {
		return Area2i{}
	} else {
		return a
	}
}

// Unions two areas.
func (a Area2i) Union(b Area2i) Area2i { return Area2i{a.X.Union(b.X), a.Y.Union(b.Y)} }
// Intersects two areas. Always returns a normalized area.
func (a Area2i) Intersect(b Area2i) Area2i { return Area2i{a.X.Intersect(b.X), a.Y.Intersect(b.Y)}.Normalized() }
// Translates an area by an offset.
func (a Area2i) Translate(v Vec2i) Area2i { return Area2i{a.X.Translate(v.X), a.Y.Translate(v.Y)} }
// Expands an area by a margin, expressed as the area around as if this had zero size.
func (a Area2i) Expand(b Area2i) Area2i { return Area2i{a.X.Expand(b.X), a.Y.Expand(b.Y)} }
// Contracts an area by a margin; the reverse of Expand.
func (a Area2i) Contract(b Area2i) Area2i { return Area2i{a.X.Contract(b.X), a.Y.Contract(b.Y)} }
// Checks if a point is within the area.
func (a Area2i) Contains(v Vec2i) bool { return a.X.Contains(v.X) && a.Y.Contains(v.Y) }
// Aligns an area within another.
func (a Area2i) Align(content Vec2i, align Alignment2i) Area2i {
	return Area2i{
		a.X.Align(content.X, align.X),
		a.Y.Align(content.Y, align.Y),
	}
}

func UnionArea2i(areas []Area2i) Area2i {
	base := Area2i{}
	for _, area := range areas {
		base = base.Union(area)
	}
	return base
}

// Part IV: Alignment Structures (see Align methods of areas)

type Alignment1i int8

const AlignStart Alignment1i = -1
const AlignMiddle Alignment1i = 0
const AlignEnd Alignment1i = 1

type Alignment2i struct {
	X Alignment1i
	Y Alignment1i
}

// Part V: Even More Utilities

func ColourFromARGB(a uint8, r uint8, g uint8, b uint8) uint32 { return (uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | (uint32(b) << 0) }

type Area1iGrid3 struct {
	A Area1i
	B Area1i
	C Area1i
}
func SplitArea1iGrid3(outer Area1i, inner Area1i) Area1iGrid3 {
	return Area1iGrid3{
		Area1i{outer.Pos, inner.Pos - outer.Pos},
		inner,
		Area1i{inner.Pos + inner.Size, (outer.Pos + outer.Size) - (inner.Pos + inner.Size)},
	}
}
func (a Area1iGrid3) AsMargin() Area1i {
	return Area1i{-a.A.Size, a.C.Size + a.A.Size}
}

type Area2iGrid3x3 struct {
	A Area2i ; B Area2i ; C Area2i
	D Area2i ; E Area2i ; F Area2i
	G Area2i ; H Area2i ; I Area2i
}
func SplitArea2iGrid3x3(outer Area2i, inner Area2i) Area2iGrid3x3 {
	xSplit := SplitArea1iGrid3(outer.X, inner.X)
	ySplit := SplitArea1iGrid3(outer.Y, inner.Y)
	return Area2iGrid3x3{
		Area2i{xSplit.A, ySplit.A}, Area2i{xSplit.B, ySplit.A}, Area2i{xSplit.C, ySplit.A},
		Area2i{xSplit.A, ySplit.B}, Area2i{xSplit.B, ySplit.B}, Area2i{xSplit.C, ySplit.B},
		Area2i{xSplit.A, ySplit.C}, Area2i{xSplit.B, ySplit.C}, Area2i{xSplit.C, ySplit.C},
	}
}
func (a Area2iGrid3x3) AsMargin() Area2i {
	return Area2iFromVecs(a.A.Size().Negate(), a.I.Size().Add(a.A.Size()))
}
