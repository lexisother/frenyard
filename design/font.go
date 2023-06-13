package design

import (
	"github.com/golang/freetype"
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/integration"
	"github.com/lexisother/frenyard/integration/fonts/roboto"
	"golang.org/x/image/font"
)

// PageTitleFont for a page title
var PageTitleFont font.Face

// GlobalFont for most text
var GlobalFont font.Face

// GlobalItalicFont is GlobalFont but italic
var GlobalItalicFont font.Face

// ButtonTextFont for buttons
var ButtonTextFont font.Face

// ListItemTextFont for main list item text
var ListItemTextFont font.Face

// ListItemSubTextFont for undertext
var ListItemSubTextFont font.Face

func deSetupFonts() {
	font, err := freetype.ParseFont(roboto.RobotoRegular)
	if err != nil {
		panic(err)
	}
	fontB, err := freetype.ParseFont(roboto.RobotoBold)
	if err != nil {
		panic(err)
	}
	fontI, err := freetype.ParseFont(roboto.RobotoItalic)
	if err != nil {
		panic(err)
	}
	PageTitleFont = integration.CreateTTFFont(font, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 24)))
	GlobalFont = integration.CreateTTFFont(font, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 16)))
	GlobalItalicFont = integration.CreateTTFFont(fontI, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 16)))
	ButtonTextFont = integration.CreateTTFFont(fontB, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 14)))

	// 16dp
	ListItemTextFont = GlobalFont
	// 14dp
	ListItemSubTextFont = integration.CreateTTFFont(font, integration.DPIPixels, float64(frenyard.Scale(DesignScale, 14)))
}
