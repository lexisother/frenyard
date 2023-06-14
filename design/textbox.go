package design

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/framework"
)

func mkTxtBox(hint string, str *string, text []string, newline bool) framework.UILayoutElement {
	// Handling of """default parameters"""
	defaultText := ""
	if len(text) > 0 {
		defaultText = text[0]
	}

	lastInput := ""
	fwTextbox := framework.NewUITextboxPtr(defaultText, hint, GlobalFont, ThemeText, ThemeTextInputSuggestion, ThemeTextInputHint, 0, frenyard.Alignment2i{X: frenyard.AlignStart})
	searchBoxContainer := framework.NewUIOverlayContainerPtr(searchboxTheme, []framework.UILayoutElement{fwTextbox})
	regenContent := func() framework.FlexboxContainer {
		lastInput = fwTextbox.Text()
		*str = lastInput
		slots := []framework.FlexboxSlot{
			{
				Grow: 1,
			},
		}
		return framework.FlexboxContainer{
			DirVertical: true,
			Slots:       slots,
		}
	}

	vboxFlex := framework.NewUIFlexboxContainerPtr(regenContent())
	handleEv := func(confirm bool) {
		if confirm && newline {
			// handle multiline textboxes
			fwTextbox.FyTInput("\n")
		}

		// The reason why we wait for stall is because this reduces the lag.
		if lastInput != fwTextbox.Text() {
			vboxFlex.SetContent(regenContent())
		}
	}

	fwTextbox.OnConfirm = func() { handleEv(true) }
	fwTextbox.OnStall = func() { handleEv(false) }

	return framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			{
				Element: searchBoxContainer,
			},
			{
				Element: vboxFlex,
				Grow:    1,
			},
		},
	})
}

func NewUITextboxPtr(hint string, str *string, init ...string) framework.UILayoutElement {
	return mkTxtBox(hint, str, init, false)
}

// NewUITextareaPtr is NewUITextboxPtr, but with multiline support
func NewUITextareaPtr(hint string, str *string, init ...string) framework.UILayoutElement {
	return mkTxtBox(hint, str, init, true)
}
