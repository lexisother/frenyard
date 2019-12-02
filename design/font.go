package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

// PageTitleFont for a page title
var PageTitleFont font.Face

// GlobalFont for most text
var GlobalFont font.Face

// ButtonTextFont for buttons
var ButtonTextFont font.Face

// ListItemTextFont for main list item text
var ListItemTextFont font.Face
// ListItemSubTextFont for undertext
var ListItemSubTextFont font.Face

func deSetupFonts() {
	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}
	fontB, err := freetype.ParseFont(gobold.TTF)
	if err != nil {
		panic(err)
	}
	PageTitleFont = frenyard.CreateTTFFont(font, frenyard.DPIPixels, float64(frenyard.Scale(DesignScale, 24)))
	GlobalFont = frenyard.CreateTTFFont(font, frenyard.DPIPixels, float64(frenyard.Scale(DesignScale, 16)))
	ButtonTextFont = frenyard.CreateTTFFont(fontB, frenyard.DPIPixels, float64(frenyard.Scale(DesignScale, 14)))
	
	// 16dp
	ListItemTextFont = GlobalFont
	// 12dp (NON-STANDARD)
	ListItemSubTextFont = frenyard.CreateTTFFont(font, frenyard.DPIPixels, float64(frenyard.Scale(DesignScale, 12)))
}
