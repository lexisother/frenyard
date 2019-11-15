package frenyard
import maths "math"
import "os"
import "strconv"

// ModifyScaleBinInt "snaps" the scale to the nearest power-of-two division (if < 1) or the nearest integer (if >= 1)
func ModifyScaleBinInt(scale float64) float64 {
	if scale < 1 {
		// 2^math.log(X, 2) == X
		// the inner round creates the 'to nearest binary fraction' effect
		scale = maths.Pow(2, maths.Round(maths.Log2(scale)))
	} else {
		scale = maths.Round(scale)
	}
	return scale
}

// Scale returns a scaled value as an integer.
func Scale(scale float64, target int32) int32 {
	return int32(maths.Ceil(float64(target) * scale))
}

// ScaleMargin1 scales a margin.
func ScaleMargin1(scale float64, target Area1i) Area1i {
	left := target.Pos
	right := target.Pos + target.Size
	left = -Scale(scale, -left)
	right = Scale(scale, right)
	return Area1i{
		Pos: left,
		Size: right - left,
	}
}

// ScaleMargin2 scales a margin.
func ScaleMargin2(scale float64, target Area2i) Area2i {
	return Area2i{
		X: ScaleMargin1(scale, target.X),
		Y: ScaleMargin1(scale, target.Y),
	}
}

// ScaleVec2i scales a Vec2i.
func ScaleVec2i(scale float64, target Vec2i) Vec2i {
	return Vec2i{
		X: Scale(scale, target.X),
		Y: Scale(scale, target.Y),
	}
}

// InferScale infers a 'reasonable scale' from environment variables or by creating a window.
func InferScale() float64 {
	env := os.Getenv("FRENYARD_SCALE")
	envScale, err := strconv.ParseFloat(env, 64)
	if err == nil {
		return envScale
	}
	
	microTestElement := NewColouredRectPtr(0xFF000000, Vec2i{128, 128})
	microTestWindow, err := CreateBoundWindow("DPI Test Window", false, 0xFF000000, microTestElement)
	if err != nil {
		return 1.0
	}
	val := microTestWindow.GetLocalDPI() / 96.0
	microTestWindow.Destroy()
	// Prevent weirdness for anything that looks even remotely like a sensible DPI value
	return ModifyScaleBinInt(val)
}
