package design

import (
	"strings"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
)

// SortListItemDetails is a sort.Interface implementation for use on ListItemDetails slices.
type SortListItemDetails []ListItemDetails

func (slid SortListItemDetails) Len() int {
	return len(slid)
}
func (slid SortListItemDetails) Swap(i int, j int) {
	a := slid[i]
	slid[i] = slid[j]
	slid[j] = a
}
func (slid SortListItemDetails) Less(i int, j int) bool {
	res := strings.Compare(strings.ToLower(slid[i].Text), strings.ToLower(slid[j].Text))
	return res < 0
}

// NewUISearchBoxPtr : Given ListItemDetails, implements a box that can search through them.
func NewUISearchBoxPtr(hint string, area []ListItemDetails) framework.UILayoutElement {
	lastRegenSearchTerm := ""
	searchBox := framework.NewUITextboxPtr("", hint, GlobalFont, ThemeText, ThemeTextInputSuggestion, ThemeTextInputHint, 0, frenyard.Alignment2i{X: frenyard.AlignStart})
	searchBoxContainer := framework.NewUIOverlayContainerPtr(searchboxTheme, []framework.UILayoutElement{searchBox})
	regenContent := func () framework.FlexboxContainer {
		lastRegenSearchTerm = searchBox.Text()
		slots := []framework.FlexboxSlot{}
		for _, v := range area {
			if strings.Contains(strings.ToLower(v.Text), strings.ToLower(lastRegenSearchTerm)) || strings.Contains(strings.ToLower(v.Subtext), strings.ToLower(lastRegenSearchTerm)) {
				slots = append(slots, framework.FlexboxSlot{
					Element: ListItem(v),
				})
			}
		}
		slots = append(slots, framework.FlexboxSlot{
			Grow: 1,
		})
		return framework.FlexboxContainer{
			DirVertical: true,
			Slots: slots,
		}
	}
	vboxFlex := framework.NewUIFlexboxContainerPtr(regenContent())
	searchBox.OnConfirm = func () {
		// The reason why we wait for stall is because this reduces the lag.
		if lastRegenSearchTerm != searchBox.Text() {
			vboxFlex.SetContent(regenContent())
		}
	}
	searchBox.OnStall = searchBox.OnConfirm
	return framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: searchBoxContainer,
			},
			framework.FlexboxSlot{
				Element: vboxFlex,
				Grow: 1,
			},
		},
	})
}
