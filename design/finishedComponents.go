package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

// ThemeText is the colour for most text.
const ThemeText = 0xFFFFFFFF

// ThemeSubText is for 'detail' text that doesn't matter that much.
const ThemeSubText = 0xFFC0C0C0

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
	Back framework.ButtonBehavior
	Title string
}

// ButtonOkAction creates a 'OK' button for some given text (likely 'OK')
func ButtonOkAction(text string, behavior framework.ButtonBehavior) *framework.UIButton {
	textElm := framework.NewUILabelPtr(integration.NewTextTypeChunk(text, ButtonTextFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return newDeUIDesignButtonPtr(0xFF2040FF, textElm, behavior)
}

// ButtonIcon creates an 'icon button'. This is an incomplete API.
func ButtonIcon(icon IconID, dp int32, click framework.ButtonBehavior) *framework.UIButton {
	return framework.NewUIButtonPtr(NewIconPtr(0xFFFFFFFF, icon, dp), click)
}

func headerConstruct(header Header) framework.UILayoutElement {
	label := framework.NewUILabelPtr(integration.NewTextTypeChunk(header.Title, PageTitleFont), ThemeText, 0, frenyard.Alignment2i{})
	if header.Back != nil {
		return framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: false,
			Slots: []framework.FlexboxSlot{
				{
					Basis: 16,
				},
				{
					Element: ButtonIcon(BackIconID, 18, header.Back),
					RespectMinimumSize: true,
				},
				{
					Basis: 16,
				},
				{
					Element: label,
					Grow: 1,
					Shrink: 1,
					Order: 0,
					RespectMinimumSize: true,
				},
				{
					Basis: 50,
				},
			},
		})
	}
	return label
}

// LayoutDocument creates a 'document' format element, with a title header & body.
func LayoutDocument(header Header, body framework.UILayoutElement, scrollable bool) framework.UILayoutElement {
	titleContent := headerConstruct(header)
	titleWrapper := framework.NewUIOverlayContainerPtr(BorderTitle(ThemeBackgroundTitle), []framework.UILayoutElement{
		titleContent,
	})
	
	body = framework.NewUIMarginContainerPtr(body, MarginBody())
	if scrollable {
		body = framework.NewUIScrollboxPtr(ScrollbarThemeV, body, true)
	}
	
	titleAndThenBody := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			{
				Element: body,
				Grow:    1,
				Shrink:  1,
				Order:   1,
			},
			{
				Element: titleWrapper,
				Order: 0,
				RespectMinimumSize: true,
			},
		},
	})
	return titleAndThenBody
}

// ButtonBar provides a 'button-bar' to put at the bottom of dialogs.
func ButtonBar(buttons []framework.UILayoutElement) framework.UILayoutElement {
	slots := []framework.FlexboxSlot{{Grow: 1}}
	for _, v := range buttons {
		slots = append(slots, framework.FlexboxSlot{
			Basis: sizeScale(32),
			Shrink: 1,
		}, framework.FlexboxSlot{
			Element: v,
			Shrink: 1,
		})
	}
	return framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		Slots: slots,
	})
}

// LayoutMsgbox provides a 'message box' body layout.
func LayoutMsgbox(text string, confirm func()) framework.UILayoutElement {
	okButton := ButtonOkAction("OK", confirm)

	buttonBar := ButtonBar([]framework.UILayoutElement{okButton})

	bodyItself := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(text, GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
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
func ScrollboxV(inwards framework.UILayoutElement) framework.UILayoutElement {
	sbox := framework.NewUIScrollboxPtr(ScrollbarThemeV, inwards, true)
	return framework.NewUIOverlayContainerPtr(ScrollboxExterior, []framework.UILayoutElement{sbox})
}

// ListItemDetails contains details for a list item.
type ListItemDetails struct {
	Icon IconID
	Text string
	// If you want this to occupy space (changing format) despite being empty, make it a space.
	Subtext string
	Click framework.ButtonBehavior
}

// ListItem sets up a list item.
// https://material.io/components/lists/#specs
func ListItem(details ListItemDetails) framework.UILayoutElement {
	var labelVertical framework.UILayoutElement
	noIconSize := int32(48)
	withIconSize := int32(56)
	if details.Subtext != "" {
		noIconSize = 64
		withIconSize = 72
		labelVertical = framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: true,
			Slots: []framework.FlexboxSlot{
				framework.FlexboxSlot{
					Grow: 1,
				},
				framework.FlexboxSlot{
					Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Text, ListItemTextFont), ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignEnd}),
				},
				framework.FlexboxSlot{
					Basis: 4,
				},
				framework.FlexboxSlot{
					Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Subtext, ListItemSubTextFont), ThemeSubText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignEnd}),
				},
				framework.FlexboxSlot{
					Grow: 1,
				},
			},
		})
	} else {
		labelVertical = framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Text, ListItemTextFont), ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignMiddle})
	}
	resSize := noIconSize
	horizontalLayoutSlots := []framework.FlexboxSlot{
		framework.FlexboxSlot{
			Basis: sizeScale(16),
		},
		framework.FlexboxSlot{
			Element: labelVertical,
			Shrink: 1,
		},
		framework.FlexboxSlot{
			Basis: sizeScale(16),
		},
	}

	if details.Icon != NullIconID {
		resSize = withIconSize
		horizontalLayoutSlots = []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
			framework.FlexboxSlot{
				Element: NewIconPtr(0xFFFFFFFF, details.Icon, 36),
			},
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
			framework.FlexboxSlot{
				Element: labelVertical,
				Shrink: 1,
			},
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
		}
	}
	horizontalLayout := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		Slots: horizontalLayoutSlots,
	})

	assembledItem := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: horizontalLayout,
				MinBasis: sizeScale(resSize),
				RespectMinimumSize: true,
			},
		},
	})
	if details.Click != nil {
		return deRippleInstall(assembledItem, borderGenSquareMaskX4Raw, borderEffectiveScale, details.Click)
	}
	return assembledItem
}
