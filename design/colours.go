package design
import "github.com/20kdc/CCUpdaterUI/frenyard"

const THEME_TEXT = 0xFFFFFFFF
const THEME_BACKGROUND = 0xFF202020
const THEME_BACKGROUND_TITLE = 0xFF404040

func ButtonContentOkAction(text string) frenyard.UIButtonThemedContent {
	textElm := frenyard.NewUILabelPtr(text, ButtonTextFont, 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return func (hover bool, down bool) (frenyard.NinePatchPackage, frenyard.UILayoutElement) {
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
