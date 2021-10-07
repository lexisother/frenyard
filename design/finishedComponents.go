package design

import (
	"github.com/yellowsink/frenyard"
	"github.com/yellowsink/frenyard/framework"
	"github.com/yellowsink/frenyard/integration"
)

// ThemeText is the colour for most text.
const ThemeText uint32 = 0xFFFFFFFF

// ThemeTextWarning is the colour for warning text.
const ThemeTextWarning uint32 = 0xFF000000

// ThemeTextInputSuggestion is the colour for text input suggestions.
const ThemeTextInputSuggestion = 0xFF80FF80

// ThemeTextInputHint is the colour for text hints
const ThemeTextInputHint = 0xFF808080

// ThemeSubText is for 'detail' text that doesn't matter that much.
const ThemeSubText uint32 = 0xFFC0C0C0

// ThemePlaceholder is the colour for placeholders.
const ThemePlaceholder uint32 = 0xFFFF0000

// ThemeBackground is the colour for most page content.
const ThemeBackground uint32 = 0xFF202020

// For debugging.
//const ThemeBackground uint32 = 0xFFFFFFFF

// ThemeBackgroundTitle is the colour for the page title background.
const ThemeBackgroundTitle uint32 = 0xFF404040

// ThemeBackgroundUnderlayer is the colour for backgrounds in "underground" lists.
const ThemeBackgroundUnderlayer uint32 = 0xFF101010

// ThemeBackgroundSearch is the colour for searchboxbackground
const ThemeBackgroundSearch uint32 = 0xFF101010

// ThemeBackgroundWarning is the colour for warning backgrounds.
const ThemeBackgroundWarning uint32 = 0xFFFFA000

// ThemeOkActionButton is the colour for OK or Install buttons.
const ThemeOkActionButton uint32 = 0xFF2040FF

// ThemeUpdateActionButton is the colour for Update buttons.
const ThemeUpdateActionButton uint32 = 0xFFD28C44

// ThemeRemoveActionButton is the colour for Remove buttons.
const ThemeRemoveActionButton uint32 = 0xFFB11E1E

// ThemeImpossibleActionButton is the colour for buttons that fail, but will explain why it is impossible.
const ThemeImpossibleActionButton uint32 = 0xFF404040

// ThemePageActionButton is the colour for buttons that show details.
const ThemePageActionButton uint32 = 0xFF205020

// Header describes a 'title' header.
type Header struct {
	Back framework.ButtonBehavior
	// If null, this is changed to BackIconID
	BackIcon IconID
	Title    string
	Forward  framework.ButtonBehavior
	// If null, this is changed to RunIconID
	ForwardIcon IconID
}

// ButtonAction creates a 'OK' button for some given text (likely 'OK')
func ButtonAction(colour uint32, text string, behavior framework.ButtonBehavior) *framework.UIButton {
	textElm := framework.NewUILabelPtr(integration.NewTextTypeChunk(text, ButtonTextFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{})
	return newDeUIDesignButtonPtr(colour, textElm, behavior)
}

// ButtonWarningFixAction creates a 'fix XYZ' button
func ButtonWarningFixAction(text string, behavior framework.ButtonBehavior) *framework.UIButton {
	textElm := framework.NewUILabelPtr(integration.NewTextTypeChunk(text, ButtonTextFont), 0xFF000000, 0, frenyard.Alignment2i{})
	return newDeUIDesignButtonPtr(0xFFEF6C00, textElm, behavior)
}

// ButtonIcon creates an 'icon button'.
func ButtonIcon(icon IconID, dp int32, click framework.ButtonBehavior) *framework.UIButton {
	return newDeUICircleButtonPtr(NewIconPtr(0xFFFFFFFF, icon, dp), click)
}

func headerConstruct(header Header) framework.UILayoutElement {
	if header.BackIcon == NullIconID {
		header.BackIcon = BackIconID
	}
	if header.ForwardIcon == NullIconID {
		header.ForwardIcon = RunIconID
	}
	label := framework.NewUILabelPtr(integration.NewTextTypeChunk(header.Title, PageTitleFont), ThemeText, 0, frenyard.Alignment2i{})

	slots := []framework.FlexboxSlot{}
	if header.Back != nil {
		slots = append(slots,
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
			framework.FlexboxSlot{
				Element:            ButtonIcon(header.BackIcon, 18, header.Back),
				RespectMinimumSize: true,
			},
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
		)
	} else {
		slots = append(slots,
			framework.FlexboxSlot{
				Basis: sizeScale(50),
			},
		)
	}
	slots = append(slots, framework.FlexboxSlot{
		Element:            label,
		Grow:               1,
		Shrink:             1,
		Order:              0,
		RespectMinimumSize: true,
	})
	if header.Forward != nil {
		slots = append(slots,
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
			framework.FlexboxSlot{
				Element:            ButtonIcon(header.ForwardIcon, 18, header.Forward),
				RespectMinimumSize: true,
			},
			framework.FlexboxSlot{
				Basis: sizeScale(16),
			},
		)
	} else {
		slots = append(slots,
			framework.FlexboxSlot{
				Basis: sizeScale(50),
			},
		)
	}
	return framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: false,
		Slots:       slots,
	})
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
				Element:            titleWrapper,
				Order:              0,
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
			Basis:  sizeScale(32),
			Shrink: 1,
		}, framework.FlexboxSlot{
			Element: v,
			Shrink:  1,
		})
	}
	return framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		Slots: slots,
	})
}

// LayoutMsgbox provides a 'message box' body layout.
func LayoutMsgbox(text string, confirm func()) framework.UILayoutElement {
	okButton := ButtonAction(ThemeOkActionButton, "OK", confirm)

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
	Click   framework.ButtonBehavior
}

// ListItem sets up a list item.
// https://material.io/components/lists/#specs
func ListItem(details ListItemDetails) framework.UILayoutElement {
	var labelVertical framework.UILayoutElement
	noIconSizeDP := int32(48)
	withIconSizeDP := int32(56)
	if details.Subtext != "" {
		noIconSizeDP = 64
		withIconSizeDP = 72
		labelVertical = framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: true,
			Slots: []framework.FlexboxSlot{
				framework.FlexboxSlot{
					Basis:              sizeScale(8),
					RespectMinimumSize: true,
					Grow:               1,
				},
				framework.FlexboxSlot{
					Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Text, ListItemTextFont), ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignEnd}),
				},
				framework.FlexboxSlot{
					Basis: sizeScale(4),
				},
				framework.FlexboxSlot{
					Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Subtext, ListItemSubTextFont), ThemeSubText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignEnd}),
				},
				framework.FlexboxSlot{
					Basis:              sizeScale(8),
					RespectMinimumSize: true,
					Grow:               1,
				},
			},
		})
	} else {
		labelVertical = framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Text, ListItemTextFont), ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignMiddle})
	}
	resSizeDP := noIconSizeDP
	horizontalLayoutSlots := []framework.FlexboxSlot{
		framework.FlexboxSlot{
			Basis: sizeScale(16),
		},
		framework.FlexboxSlot{
			Element: labelVertical,
			Shrink:  1,
		},
		framework.FlexboxSlot{
			Basis: sizeScale(16),
		},
	}

	if details.Icon != NullIconID {
		resSizeDP = withIconSizeDP
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
				Shrink:  1,
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
				Element:            horizontalLayout,
				MinBasis:           sizeScale(resSizeDP),
				RespectMinimumSize: true,
			},
		},
	})
	if details.Click != nil {
		return deRippleInstall(assembledItem, borderGenSquareMaskX4Raw, borderEffectiveScale, details.Click)
	}
	return assembledItem
}

// InformationPanelDetails tags a block of content with an icon
type InformationPanelDetails struct {
	Text       string
	ActionText string
	Action     framework.ButtonBehavior
}

// InformationPanel applies InformationPanelDetails
func InformationPanel(details InformationPanelDetails) framework.UILayoutElement {
	icon := WarningIconID
	margin := sizeScale(8)
	hMargin := sizeScale(4)

	var body framework.UILayoutElement
	body = framework.NewUILabelPtr(integration.NewTextTypeChunk(details.Text, GlobalFont), ThemeTextWarning, 0, frenyard.Alignment2i{X: frenyard.AlignStart})

	primaryHorizontalSlots := []framework.FlexboxSlot{
		framework.FlexboxSlot{
			Element: framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
				DirVertical: true,
				Slots: []framework.FlexboxSlot{
					framework.FlexboxSlot{
						Element:            NewIconPtr(ThemeTextWarning, icon, 18),
						RespectMinimumSize: true,
					},
					framework.FlexboxSlot{
						Grow: 1,
					},
				},
			}),
			RespectMinimumSize: true,
		},
		framework.FlexboxSlot{
			Basis: margin,
		},
		framework.FlexboxSlot{
			Element:            body,
			RespectMinimumSize: true,
			Grow:               1,
			Shrink:             1,
		},
	}

	if details.Action != nil {
		primaryHorizontalSlots = append(primaryHorizontalSlots,
			framework.FlexboxSlot{
				Basis: margin,
			},
			framework.FlexboxSlot{
				Element:            ButtonWarningFixAction(details.ActionText, details.Action),
				RespectMinimumSize: true,
			},
		)
	}
	return framework.NewUIOverlayContainerPtr(framework.NinePatchFrame{
		Layers: []framework.NinePatchFrameLayer{
			{
				Pass:      framework.FramePassOverBefore,
				ColourMod: ThemeBackgroundWarning,
				NinePatch: borderGenRounded4dpMaskX4Mask.Inset(frenyard.Area2iMargin(0, 0, 0, hMargin*2)),
				Scale:     borderEffectiveScale,
			},
		},
		Clipping: true,
		Padding:  frenyard.Area2iMargin(margin, margin, margin, margin+(hMargin*2)),
	}, []framework.UILayoutElement{framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{DirVertical: false, Slots: primaryHorizontalSlots})})
}
