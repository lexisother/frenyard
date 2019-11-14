package frenyard

// Scale is a fixed-point value measured by /16 used for scaling, with 0 being 'no change'. It is done this way to make +/- 1 "make sense".
type Scale int32

// ScaleRM describes a rounding mode for scaling values.
type ScaleRM uint8
// ScaleRMCeil is ceil(target * scale)
const ScaleRMCeil ScaleRM = 1
// ScaleRMInteger is (target * ceil(scale))
const ScaleRMInteger ScaleRM = 2
// ScaleRMBinary for Scale >= 0 is ScaleRMInteger, and otherwise accepts the lowest of -8, -12, -14, and -15.
const ScaleRMBinary ScaleRM = 3

// ScaleRMNinePatch is the scaling method used for NinePatches.
const ScaleRMNinePatch = ScaleRMBinary

// ScaleFromFloat64 approximately converts a floating-point number to a scale.
func ScaleFromFloat64(v float64) Scale {
	return Scale(v * 16) - 16
}

// Scale returns a scaled value as an integer.
func (s Scale) Scale(target int32, rm ScaleRM) int32 {
	ints := int32(s) + 16
	if rm == ScaleRMBinary {
		if ints < 16 {
			if ints <= 1 {
				ints = 1
			} else if ints <= 2 {
				ints = 2
			} else if ints <= 4 {
				ints = 4
			} else if ints <= 8 {
				ints = 8
			} else {
				ints = 16
			}
		} else {
			rm = ScaleRMInteger
		}
	}
	if rm == ScaleRMInteger {
		ints = (ints + 15) & 0x7FFFFFF0
	}
	// ScaleRMCeil
	return ((ints * target) + 15) >> 4
}

// Margin1 scales a margin.
func (s Scale) Margin1(target Area1i, rm ScaleRM) Area1i {
	left := target.Pos
	right := target.Pos + target.Size
	left = -s.Scale(-left, rm)
	right = s.Scale(right, rm)
	return Area1i{
		Pos: left,
		Size: right - left,
	}
}

// Margin2 scales a margin.
func (s Scale) Margin2(target Area2i, rm ScaleRM) Area2i {
	return Area2i{
		X: s.Margin1(target.X, rm),
		Y: s.Margin1(target.Y, rm),
	}
}
