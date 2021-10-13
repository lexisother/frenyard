package design

import (
	"github.com/golang/freetype"
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/integration"
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
	PageTitleFont = integration.CreateTTFFont(font, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 24)))
	GlobalFont = integration.CreateTTFFont(font, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 16)))
	ButtonTextFont = integration.CreateTTFFont(fontB, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 14)))

	// 16dp
	ListItemTextFont = GlobalFont
	// 14dp
	ListItemSubTextFont = integration.CreateTTFFont(font, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 14)))
}
