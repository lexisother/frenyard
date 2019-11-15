package frenyard
import maths "math"

// ScaleRM describes a rounding mode for scaling values.
type ScaleRM uint8
// ScaleRMCeil is ceil(target * scale)
const ScaleRMCeil ScaleRM = 1
// ScaleRMInteger is (target * ceil(scale)) - this is a safe method for pixel art but cannot handle scales under 1.
const ScaleRMInteger ScaleRM = 2
// ScaleRMBinary uses the highest binary fraction. If given images scaled up using the method in imaging.go by another binary fraction, it's the highest-quality method, but it ignores a wide range of input values. ScaleRMBinInt may be more suitable.
const ScaleRMBinary ScaleRM = 3
// ScaleRMBinInt for scale >= 1 is ScaleRMInteger, and otherwise ScaleRMBinary.
const ScaleRMBinInt ScaleRM = 4
// ScaleRMIntBin for scale >= 1 is ScaleRMBinary, and otherwise ScaleRMInteger.
const ScaleRMIntBin ScaleRM = 5

// ScaleRMNinePatch is the scaling method used for NinePatches.
const ScaleRMNinePatch = ScaleRMBinInt

// Scale returns a scaled value as an integer.
func Scale(scale float64, target int32, rm ScaleRM) int32 {
	if rm == ScaleRMBinInt {
		if scale >= 1 {
			rm = ScaleRMInteger
		} else {
			rm = ScaleRMBinary
		}
	}
	if rm == ScaleRMIntBin {
		if scale >= 1 {
			rm = ScaleRMBinary
		} else {
			rm = ScaleRMInteger
		}
	}
	if rm == ScaleRMBinary {
		// 2^math.log(X, 2) == X
		// the inner ceil creates the 'to highest binary fraction' effect
		scale = maths.Pow(2, maths.Ceil(maths.Log2(scale)))
		rm = ScaleRMCeil
	}
	if rm == ScaleRMInteger {
		scale = maths.Ceil(scale)
		rm = ScaleRMCeil
	}
	// ScaleRMCeil (default, used by other types)
	return int32(maths.Ceil(float64(target) * scale))
}

// ScaleMargin1 scales a margin.
func ScaleMargin1(scale float64, target Area1i, rm ScaleRM) Area1i {
	left := target.Pos
	right := target.Pos + target.Size
	left = -Scale(scale, -left, rm)
	right = Scale(scale, right, rm)
	return Area1i{
		Pos: left,
		Size: right - left,
	}
}

// ScaleMargin2 scales a margin.
func ScaleMargin2(scale float64, target Area2i, rm ScaleRM) Area2i {
	return Area2i{
		X: ScaleMargin1(scale, target.X, rm),
		Y: ScaleMargin1(scale, target.Y, rm),
	}
}
