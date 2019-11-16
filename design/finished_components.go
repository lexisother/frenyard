package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// ThemeText is the colour for most text.
const ThemeText = 0xFFFFFFFF

// ThemeBackground is the colour for most page content.
const ThemeBackground = 0xFF202020
// For debugging.
//const ThemeBackground = 0xFFFFFFFF

// ThemeBackgroundTitle is the colour for the page title background.
const ThemeBackgroundTitle = 0xFF404040

// ButtonOkAction creates a 'OK' button for some given text (likely 'OK')
func ButtonOkAction(text string, behavior frenyard.ButtonBehavior) *frenyard.UIButton {
	textElm := frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk(text, ButtonTextFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return newDeUIDesignButtonPtr(textElm, textElm, behavior)
}
