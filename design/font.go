package design
import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gobold"
)

// PageTitleFont for a page title
var PageTitleFont font.Face
// GlobalFont for most text
var GlobalFont font.Face
// ButtonTextFont for buttons
var ButtonTextFont font.Face

func init() {
	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}
	fontB, err := freetype.ParseFont(gobold.TTF)
	if err != nil {
		panic(err)
	}
	PageTitleFont = frenyard.CreateTTFFont(font, 72, 24)
	GlobalFont = frenyard.CreateTTFFont(font, 72, 16)
	ButtonTextFont = frenyard.CreateTTFFont(fontB, 72, 14)
}
