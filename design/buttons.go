package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// ButtonContentOkAction creates a 'OK' button theme for some given text (likely 'OK')
func ButtonContentOkAction(text string) frenyard.UIButtonThemedContent {
	textElm := frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk(text, ButtonTextFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return func(hover bool, down bool) (frenyard.NinePatchPackage, frenyard.UILayoutElement) {
		textElm.SetColour(0xFFFFFFFF)
		if down {
			textElm.SetColour(0xFF404040)
			return BorderButton(0xFF102080, true), textElm
		}
		if hover {
			return BorderButton(0xFF4060FF, true), textElm
		}
		return BorderButton(0xFF2040FF, false), textElm
	}
}
