package main

import (
	"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

func main() {
	frenyard.TargetFrameTime = 0.016
	// Ok, now start...
	upShowFailureToFindGameDialog()
	// Started!
	frenyard.GlobalBackend.Run(func(frameTime float64) {
	})
}

type upTextPanel struct {
	frenyard.UILayoutProxy
	textTitle             *frenyard.UILabel
	text                  *frenyard.UILabel
	bodyMode              bool
	remainingToAddToLabel string
	counter               float64
}
func newUpTextPanelPtr(text string, ok func ()) *upTextPanel {
	elem := &upTextPanel{
		remainingToAddToLabel: text,
	}
	elem.text = frenyard.NewUILabelPtr("", design.GlobalFont, design.ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignStart})
	elem.textTitle = frenyard.NewUILabelPtr("", design.PageTitleFont, design.ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignMiddle, Y: frenyard.AlignStart})

	testButtonWrapper := frenyard.NewUIButtonPtr(design.ButtonContentOkAction("OK"), ok)

	buttonBar := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		Slots: []frenyard.FlexboxSlot{
			{Grow: 1},
			{
				Element: testButtonWrapper,
			},
		},
	})

	bodyItself := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		DirVertical: true,
		Slots: []frenyard.FlexboxSlot{
			{
				Element: elem.text,
				Grow:    1,
			},
			{Basis: design.SizeMarginAroundEverything},
			{
				Element: buttonBar,
			},
		},
	})

	titleWrapper := frenyard.NewUIOverlayContainerPtr(design.BorderTitle(design.ThemeBackgroundTitle), []frenyard.UILayoutElement{
		elem.textTitle,
	})

	titleAndThenBody := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		DirVertical: true,
		Slots: []frenyard.FlexboxSlot{
			{
				Element: frenyard.NewUIMarginContainerPtr(bodyItself, design.MarginBody()),
				Grow:    1,
				Shrink:  1,
				Order:   1,
			},
			{
				Element:  titleWrapper,
				MinBasis: design.SizeTitleHeight,
				Order:    0,
			},
		},
	})

	frenyard.InitUILayoutProxy(elem, titleAndThenBody)
	elem.FyEResize(frenyard.Vec2i{X: 320, Y: 200})
	return elem
}
func (dialog *upTextPanel) FyETick(seconds float64) {
	dialog.counter += seconds
	for dialog.counter > 0.05 {
		dialog.counter -= 0.05
		cutPoint := len(dialog.remainingToAddToLabel)
		if cutPoint == 0 {
			dialog.counter = 0
			return
		}
		for idx := range dialog.remainingToAddToLabel {
			if idx != 0 {
				cutPoint = idx
				break
			}
		}
		piece := dialog.remainingToAddToLabel[:cutPoint]
		dialog.remainingToAddToLabel = dialog.remainingToAddToLabel[cutPoint:]
		if piece == "\t" {
			dialog.bodyMode = !dialog.bodyMode
			if !dialog.bodyMode {
				dialog.textTitle.SetText("")
				dialog.text.SetText("")
			}
		} else {
			if !dialog.bodyMode {
				dialog.textTitle.SetText(dialog.textTitle.Text() + piece)
			} else {
				dialog.text.SetText(dialog.text.Text() + piece)
			}
		}
	}
}

func upShowFailureToFindGameDialog() {
	var elem *frenyard.UISlideTransitionContainer
	slideNum := 0
	var slideFn func ()
	slideFn = func () {
		elem.TransitionTo(newUpTextPanelPtr(fmt.Sprintf("Slide %v\tHello from slide %v\nIt's nice here", slideNum, slideNum), slideFn), 1.0, slideNum & 2 == 0, slideNum & 1 == 0)
		slideNum++
	}
	elem = frenyard.NewUISlideTransitionContainerPtr(newUpTextPanelPtr("Test\tMeep", slideFn))
	_, err := frenyard.CreateBoundWindow("CCUpdaterUI Installation Helper", true, design.ThemeBackground, elem)
	if err != nil {
		panic(err)
	}
}
