package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// ThemeText is the colour for most text.
const ThemeText = 0xFFFFFFFF

// ThemeBackground is the colour for most page content.
const ThemeBackground = 0xFF202020
//const ThemeBackground = 0xFFFFFFFF

// ThemeBackgroundTitle is the colour for the page title background.
const ThemeBackgroundTitle = 0xFF404040

// ButtonContentOkAction creates a 'OK' button theme for some given text (likely 'OK')
func ButtonContentOkAction(text string) frenyard.UIButtonThemedContent {
	textElm := frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk(text, ButtonTextFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return func(hover bool, down bool) (frenyard.NinePatchPackage, frenyard.UILayoutElement) {
		textElm.SetColour(0xFFFFFFFF)
		if down {
			textElm.SetColour(0xFF404040)
			return BorderButton(0xFF102080), textElm
		}
		if hover {
			return BorderButton(0xFF4060FF), textElm
		}
		return BorderButton(0xFF2040FF), textElm
	}
}
