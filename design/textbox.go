package design

import (
	"github.com/uwu/frenyard"
	"github.com/uwu/frenyard/framework"
)

func NewUITextboxPtr(hint string, str *string) framework.UILayoutElement {
	lastInput := ""
	searchBox := framework.NewUITextboxPtr("", hint, GlobalFont, ThemeText, ThemeTextInputSuggestion, ThemeTextInputHint, 0, frenyard.Alignment2i{X: frenyard.AlignStart})
	searchBoxContainer := framework.NewUIOverlayContainerPtr(searchboxTheme, []framework.UILayoutElement{searchBox})
	regenContent := func() framework.FlexboxContainer {
		lastInput = searchBox.Text()
		*str = lastInput
		slots := []framework.FlexboxSlot{}
		slots = append(slots, framework.FlexboxSlot{
			Grow: 1,
		})
		return framework.FlexboxContainer{
			DirVertical: true,
			Slots:       slots,
		}
	}
	vboxFlex := framework.NewUIFlexboxContainerPtr(regenContent())
	searchBox.OnConfirm = func() {
		// The reason why we wait for stall is because this reduces the lag.
		if lastInput != searchBox.Text() {
			vboxFlex.SetContent(regenContent())
		}
	}
	searchBox.OnStall = searchBox.OnConfirm
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
