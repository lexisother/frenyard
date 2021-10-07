package frenyard

// Constants

// SizeUnlimited is equal to the maximum value an int32 can have, represents an infinite value. Note that 0x80000000 is reserved; -0x7FFFFFFF is considered the definitive 'negative unlimited' value for simplicity's sake.
const SizeUnlimited int32 = 0x7FFFFFFF

// Vec2iUnlimited returns a Vec2i of unlimited size.
func Vec2iUnlimited() Vec2i { return Vec2i{SizeUnlimited, SizeUnlimited} }

// AddCU performs an addition with consideration for unlimited values.
func AddCU(a int32, b int32) int32 {
	if a == SizeUnlimited {
		if b == -SizeUnlimited {
			return 0
		}
		return SizeUnlimited
	} else if b == SizeUnlimited {
		if a == -SizeUnlimited {
			return 0
		}
		return SizeUnlimited
	} else if a == -SizeUnlimited {
		return -SizeUnlimited
	} else if b == -SizeUnlimited {
		return -SizeUnlimited
	}
	return a + b
}

// Part I: Basic maths

// Max returns the higher of two int32 values.
func Max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller of two int32 values.
func Min(a int32, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Part II: Point Type

// Vec2i is the basic 2-dimensional vector type.
type Vec2i struct {
	X, Y int32
}

// Add adds two Vec2is.
func (a Vec2i) Add(b Vec2i) Vec2i {
	return Vec2i{AddCU(a.X, b.X), AddCU(a.Y, b.Y)}
}

// Min returns a Vec2i with each axis value being the minimum of the two input Vec2i values on that axis.
func (a Vec2i) Min(b Vec2i) Vec2i {
	return Vec2i{Min(a.X, b.X), Min(a.Y, b.Y)}
}

// Max is similar to Min, but uses maximum, not minimum.
func (a Vec2i) Max(b Vec2i) Vec2i {
	return Vec2i{Max(a.X, b.X), Max(a.Y, b.Y)}
}

// Eq compares two Vec2is and returns true for equal.
func (a Vec2i) Eq(b Vec2i) bool {
	return a.X == b.X && a.Y == b.Y
}

// These two are framed in the sense of 'an area of size A could contain an area of size B'.

// Gt checks if an area of size A could hold an area of size B with room to spare.
func (a Vec2i) Gt(b Vec2i) bool {
	return a.X > b.X && a.Y > b.Y
}

// Ge checks if an area of size A could hold an area of size B (with or without spare room).
func (a Vec2i) Ge(b Vec2i) bool {
	return a.X >= b.X && a.Y >= b.Y
}

// Negate returns the vector multiplied by -1.
func (a Vec2i) Negate() Vec2i {
	return Vec2i{-a.X, -a.Y}
}

// ConditionalTranspose conditionally swaps X/Y, which is useful when the coordinate system is variable.
func (a Vec2i) ConditionalTranspose(yes bool) Vec2i {
	if yes {
		return Vec2i{a.Y, a.X}
	}
	return Vec2i{a.X, a.Y}
}

// Part III: Area Types (AABBs)

// Area1i is a 1-dimensional axis-aligned area, which is a useful primitive for N-dimensional areas. Please note that Area1i, and by extension Area2i, may malfunction if given infinite values.
type Area1i struct {
	Pos  int32
	Size int32
}

// Area1iOfSize returns an Area1i of the given size (covering range inclusive 0 to exclusive size)
func Area1iOfSize(a int32) Area1i { return Area1i{0, a} }

// Area1iMargin is a quick idiom for a margin of a given size.
func Area1iMargin(l int32, r int32) Area1i { return Area1i{-l, l + r} }

// Empty returns Size <= 0
func (a Area1i) Empty() bool {
	return a.Size <= 0
}

// Normalized replaces empty areas with zeroed areas.
func (a Area1i) Normalized() Area1i {
	if a.Empty() {
		return Area1i{}
	}
	return a
}

// Union unions two areas.
func (a Area1i) Union(b Area1i) Area1i {
	pos := Min(a.Pos, b.Pos)
	end := Max(a.Pos+a.Size, b.Pos+b.Size)
	return Area1i{
		pos,
		end - pos,
	}
}

// Intersect intersects two areas. Always returns a normalized area.
func (a Area1i) Intersect(b Area1i) Area1i {
	pos := Max(a.Pos, b.Pos)
	end := Min(a.Pos+a.Size, b.Pos+b.Size)
	return Area1i{
		pos,
		end - pos,
	}.Normalized()
}

// Translate translates an area by an offset.
func (a Area1i) Translate(i int32) Area1i { return Area1i{a.Pos + i, a.Size} }

// Expand expands an area by a margin, expressed as the area around as if this had zero size.
func (a Area1i) Expand(n Area1i) Area1i { return Area1i{a.Pos + n.Pos, a.Size + n.Size} }

// Contract contracts an area by a margin; the reverse of Expand.
func (a Area1i) Contract(n Area1i) Area1i { return Area1i{a.Pos - n.Pos, a.Size - n.Size} }

// Contains checks if a point is within the area.
func (a Area1i) Contains(i int32) bool { return (i >= a.Pos) && (i < a.Pos+a.Size) }

// Align aligns an area within another.
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

// UnionArea1i gets an area containing all areas in the given slice
func UnionArea1i(areas []Area1i) Area1i {
	base := Area1i{}
	for _, area := range areas {
		base = base.Union(area)
	}
	return base
}

// Area2i is the basic rectangle type. Please note that Area2i may malfunction if given infinite values.
type Area2i struct {
	X Area1i
	Y Area1i
}

// Pos returns a 'position' vector for the area. (see Area2iFromVecs)
func (a Area2i) Pos() Vec2i { return Vec2i{a.X.Pos, a.Y.Pos} }

// Size returns a 'size' vector for the area (see Area2iFromVecs)
func (a Area2i) Size() Vec2i { return Vec2i{a.X.Size, a.Y.Size} }

// Area2iFromVecs returns an Area2i made from a position/size vector pair.
func Area2iFromVecs(pos Vec2i, size Vec2i) Area2i {
	return Area2i{
		Area1i{pos.X, size.X},
		Area1i{pos.Y, size.Y},
	}
}

// Area2iOfSize returns an Area2i of the given size.
func Area2iOfSize(a Vec2i) Area2i { return Area2i{Area1iOfSize(a.X), Area1iOfSize(a.Y)} }

// Area2iMargin is a quick idiom for a margin of a given size.
func Area2iMargin(l int32, u int32, r int32, d int32) Area2i {
	return Area2i{Area1iMargin(l, r), Area1iMargin(u, d)}
}

// Empty returns Size <= 0
func (a Area2i) Empty() bool {
	return a.X.Empty() || a.Y.Empty()
}

// Normalized replaces empty areas with zeroed areas.
func (a Area2i) Normalized() Area2i {
	if a.Empty() {
		return Area2i{}
	}
	return a
}

// Union unions two areas.
func (a Area2i) Union(b Area2i) Area2i { return Area2i{a.X.Union(b.X), a.Y.Union(b.Y)} }

// Intersect intersects two areas. Always returns a normalized area.
func (a Area2i) Intersect(b Area2i) Area2i {
	return Area2i{a.X.Intersect(b.X), a.Y.Intersect(b.Y)}.Normalized()
}

// Translate translates an area by an offset.
func (a Area2i) Translate(v Vec2i) Area2i { return Area2i{a.X.Translate(v.X), a.Y.Translate(v.Y)} }

// Expand expands an area by a margin, expressed as the area around as if this had zero size.
func (a Area2i) Expand(b Area2i) Area2i { return Area2i{a.X.Expand(b.X), a.Y.Expand(b.Y)} }

// Contract contracts an area by a margin; the reverse of Expand.
func (a Area2i) Contract(b Area2i) Area2i { return Area2i{a.X.Contract(b.X), a.Y.Contract(b.Y)} }

// Contains checks if a point is within the area.
func (a Area2i) Contains(v Vec2i) bool { return a.X.Contains(v.X) && a.Y.Contains(v.Y) }

// Align aligns an area within another.
func (a Area2i) Align(content Vec2i, align Alignment2i) Area2i {
	return Area2i{
		a.X.Align(content.X, align.X),
		a.Y.Align(content.Y, align.Y),
	}
}

// UnionArea2i gets an area containing all areas in the given slice
func UnionArea2i(areas []Area2i) Area2i {
	base := Area2i{}
	for _, area := range areas {
		base = base.Union(area)
	}
	return base
}

// Part IV: Alignment Structures (see Align methods of areas)

// Alignment1i specifies an alignment preference along an axis.
type Alignment1i int8

// AlignStart aligns the element at the start (left/top).
const AlignStart Alignment1i = -1

// AlignMiddle aligns the element at the centre.
const AlignMiddle Alignment1i = 0

// AlignEnd aligns the element at the end (right/bottom).
const AlignEnd Alignment1i = 1

// Alignment2i contains two Alignment1is (one for each axis).
type Alignment2i struct {
	X Alignment1i
	Y Alignment1i
}

// Part V: Even More Utilities

// Area1iGrid3 is an Area1i split into a left area, the original 'inner' area, and the right area.
type Area1iGrid3 struct {
	A Area1i
	B Area1i
	C Area1i
}

// SplitArea1iGrid3 splits an Area1i into 3 sections using the bounds of an inner Area1i.
func SplitArea1iGrid3(outer Area1i, inner Area1i) Area1iGrid3 {
	return Area1iGrid3{
		Area1i{outer.Pos, inner.Pos - outer.Pos},
		inner,
		Area1i{inner.Pos + inner.Size, (outer.Pos + outer.Size) - (inner.Pos + inner.Size)},
	}
}

// AsMargin returns a margin around the centre Area1i for Expand.
func (a Area1iGrid3) AsMargin() Area1i {
	return Area1i{-a.A.Size, a.C.Size + a.A.Size}
}

// Area2iGrid3x3 is Area1iGrid3 in two dimensions.
type Area2iGrid3x3 struct {
	A Area2i
	B Area2i
	C Area2i
	D Area2i
	E Area2i
	F Area2i
	G Area2i
	H Area2i
	I Area2i
}

// SplitArea2iGrid3x3 splits an Area2i into 9 sections using the bounds of an inner Area2i.
func SplitArea2iGrid3x3(outer Area2i, inner Area2i) Area2iGrid3x3 {
	xSplit := SplitArea1iGrid3(outer.X, inner.X)
	ySplit := SplitArea1iGrid3(outer.Y, inner.Y)
	return Area2iGrid3x3{
		Area2i{xSplit.A, ySplit.A}, Area2i{xSplit.B, ySplit.A}, Area2i{xSplit.C, ySplit.A},
		Area2i{xSplit.A, ySplit.B}, Area2i{xSplit.B, ySplit.B}, Area2i{xSplit.C, ySplit.B},
		Area2i{xSplit.A, ySplit.C}, Area2i{xSplit.B, ySplit.C}, Area2i{xSplit.C, ySplit.C},
	}
}

// AsMargin returns a margin around the centre Area2i for Expand.
func (a Area2iGrid3x3) AsMargin() Area2i {
	return Area2iFromVecs(a.A.Size().Negate(), a.I.Size().Add(a.A.Size()))
}

// Part VI: Easings

// EasingQuadraticIn is a quadratic ease-in function, which in practice means it just squares the input value.
func EasingQuadraticIn(point float64) float64 {
	return point * point
}

// EasingInOut currys a function. Given an ease-in function, returns an ease in-out function.
func EasingInOut(easeIn func(float64) float64) func(float64) float64 {
	return func(point float64) float64 {
		if point < 0.5 {
			return easeIn(point*2.0) / 2.0
		}
		return 1.0 - (easeIn(2.0-(point*2.0)) / 2.0)
	}
}
