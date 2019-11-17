package design

import "github.com/20kdc/CCUpdaterUI/frenyard"

// ThemeText is the colour for most text.
const ThemeText = 0xFFFFFFFF

// ThemePlaceholder is the colour for placeholders.
const ThemePlaceholder = 0xFFFF0000

// ThemeBackground is the colour for most page content.
const ThemeBackground = 0xFF202020
// For debugging.
//const ThemeBackground = 0xFFFFFFFF

// ThemeBackgroundTitle is the colour for the page title background.
const ThemeBackgroundTitle = 0xFF404040

// ThemeBackgroundUnderlayer is the colour for backgrounds in "underground" lists.
const ThemeBackgroundUnderlayer = 0xFF101010

// Header describes a 'title' header.
type Header struct {
	Title string
}

// ButtonOkAction creates a 'OK' button for some given text (likely 'OK')
func ButtonOkAction(text string, behavior frenyard.ButtonBehavior) *frenyard.UIButton {
	textElm := frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk(text, ButtonTextFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return newDeUIDesignButtonPtr(textElm, textElm, behavior)
}

// LayoutDocument creates a 'document' format element, with a title header & body.
func LayoutDocument(header Header, body frenyard.UILayoutElement, scrollable bool) frenyard.UILayoutElement {
	titleWrapper := frenyard.NewUIOverlayContainerPtr(BorderTitle(ThemeBackgroundTitle), []frenyard.UILayoutElement{
		frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk(header.Title, PageTitleFont), ThemeText, 0, frenyard.Alignment2i{}),
	})
	
	body = frenyard.NewUIMarginContainerPtr(body, MarginBody())
	if scrollable {
		body = frenyard.NewUIScrollboxPtr(ScrollbarThemeV, body, true)
	}
	
	titleAndThenBody := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		DirVertical: true,
		Slots: []frenyard.FlexboxSlot{
			{
				Element: body,
				Grow:    1,
				Shrink:  1,
				Order:   1,
			},
			{
				Element:    titleWrapper,
				Shrink:                1,
				Order:                 0,
				RespectMinimumSize: true,
			},
		},
	})
	return titleAndThenBody
}

// LayoutMsgbox provides a 'message box' body layout.
func LayoutMsgbox(text string, confirm func()) frenyard.UILayoutElement {
	okButton := ButtonOkAction("OK", confirm)

	buttonBar := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		Slots: []frenyard.FlexboxSlot{
			{Grow: 1},
			{
				Element: okButton,
			},
		},
	})

	bodyItself := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		DirVertical: true,
		Slots: []frenyard.FlexboxSlot{
			{
				Element: frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk(text, GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
				Grow:    1,
			},
			{Basis: SizeMarginAroundEverything},
			{
				Element: buttonBar,
			},
		},
	})
	return bodyItself
}

// ScrollboxV sets up a fancy scrollbox with added decoration.
func ScrollboxV(inwards frenyard.UILayoutElement) frenyard.UILayoutElement {
	sbox := frenyard.NewUIScrollboxPtr(ScrollbarThemeV, inwards, true)
	return frenyard.NewUIOverlayContainerPtr(ScrollboxExterior, []frenyard.UILayoutElement{sbox})
}
