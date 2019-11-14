package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// SizeMarginAroundEverything is a useful margin around the body, etc.
var SizeMarginAroundEverything int32

// SizeTextNudge is the amount to nudge text vertically downwards to make it seem even.
var SizeTextNudge int32

func deSetupSizes() {
	SizeMarginAroundEverything = DesignScale.Scale(16, frenyard.ScaleRMCeil)
	SizeTextNudge = DesignScale.Scale(4, frenyard.ScaleRMCeil)
}

// MarginBody is the amount to push the page body by.
func MarginBody() frenyard.Area2i {
	return frenyard.Area2i{
		X: frenyard.Area1i{Pos: -SizeMarginAroundEverything, Size: SizeMarginAroundEverything * 2},
		Y: frenyard.Area1i{Pos: -(SizeMarginAroundEverything + SizeTextNudge), Size: (SizeMarginAroundEverything * 2) + SizeTextNudge},
	}
}
