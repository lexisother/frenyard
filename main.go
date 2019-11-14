package main

import (
	"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

func main() {
	design.Setup(frenyard.ScaleFromFloat64(2))
	fmt.Printf("%v\n", design.DesignScale)
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
	textTitleText         string
	text                  *frenyard.UILabel
	textText              string
	bodyMode              bool
	remainingToAddToLabel string
	counter               float64
}
func newUpTextPanelPtr(text string, ok func ()) *upTextPanel {
	elem := &upTextPanel{
		remainingToAddToLabel: text,
	}
	elem.text = frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignStart})
	elem.textTitle = frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("", design.PageTitleFont), design.ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignMiddle, Y: frenyard.AlignStart})

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
		} else {
			if !dialog.bodyMode {
				dialog.textTitleText += piece
				dialog.textTitle.SetText(frenyard.NewTextTypeChunk(dialog.textTitleText, design.PageTitleFont))
			} else {
				dialog.textText += piece
				dialog.text.SetText(frenyard.NewTextTypeChunk(dialog.textText, design.GlobalFont))
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
	elem = frenyard.NewUISlideTransitionContainerPtr(newUpTextPanelPtr("Test yup YUP a test\tMeep", slideFn))
	wnd, err := frenyard.CreateBoundWindow("CCUpdaterUI Installation Helper", true, design.ThemeBackground, elem)
	fmt.Printf("%v\n", wnd.GetLocalDPI())
	if err != nil {
		panic(err)
	}
}
