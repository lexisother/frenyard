package design

// "go generate ./..." to regenerate bindata!
//go:generate go run ./data-compiler borderButtonX1 borderHeaderX2

// DesignScale is the current scale for the Design.
var DesignScale float64

// Setup sets the sizes, fonts and borders according to the given scale.
func Setup(scale float64) {
	DesignScale = scale
	deSetupSizes()
	deSetupFonts()
}
