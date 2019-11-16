package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// SizeMarginAroundEverything is a useful margin around the body, etc.
var SizeMarginAroundEverything int32

// SizeTextNudge is the amount to nudge text vertically downwards to make it seem even.
var SizeTextNudge int32

// SizeWindow is the size of the main window.
var SizeWindow frenyard.Vec2i

// sizeScale scales an integer size.
func sizeScale(size int32) int32 {
	return frenyard.Scale(DesignScale, size)
}

func deSetupSizes() {
	SizeMarginAroundEverything = sizeScale(16)
	SizeTextNudge = sizeScale(4)
	SizeWindow = frenyard.ScaleVec2i(DesignScale, frenyard.Vec2i{X: 320, Y: 200})
}

// MarginBody is the amount to push the page body by.
func MarginBody() frenyard.Area2i {
	return frenyard.Area2i{
		X: frenyard.Area1i{Pos: -SizeMarginAroundEverything, Size: SizeMarginAroundEverything * 2},
		Y: frenyard.Area1i{Pos: -(SizeMarginAroundEverything + SizeTextNudge), Size: (SizeMarginAroundEverything * 2) + SizeTextNudge},
	}
}
