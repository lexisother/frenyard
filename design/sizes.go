package design
import "github.com/20kdc/CCUpdaterUI/frenyard"

// SizeTitleHeight is the size of the page title bar.
const SizeTitleHeight = 40
// SizeMarginAroundEverything is a useful margin around the body, etc.
const SizeMarginAroundEverything = 16
// SizeTextNudge is the amount to nudge text vertically downwards to make it seem even.
const SizeTextNudge = 4

// MarginBody is the amount to push the page body by.
func MarginBody() frenyard.Area2i {
	return frenyard.Area2i{
		X: frenyard.Area1i{Pos: -SizeMarginAroundEverything, Size: SizeMarginAroundEverything * 2},
		Y: frenyard.Area1i{Pos: -(SizeMarginAroundEverything + SizeTextNudge), Size: (SizeMarginAroundEverything * 2) + SizeTextNudge},
	}
}
