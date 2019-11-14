package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// DesignScale is the current scale for the Design.
var DesignScale frenyard.Scale

// Setup sets the sizes, fonts and borders according to the given scale.
func Setup(scale frenyard.Scale) {
	DesignScale = scale
	deSetupSizes()
	deSetupFonts()
}
