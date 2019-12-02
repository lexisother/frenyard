package main

import (
	"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

type upApplication struct {
	slideContainer *frenyard.UISlideTransitionContainer
	window frenyard.Window
	upQueued chan func()
}

func main() {
	design.Setup(frenyard.InferScale())
	frenyard.TargetFrameTime = 0.016
	slideContainer := frenyard.NewUISlideTransitionContainerPtr(nil)
	slideContainer.FyEResize(design.SizeWindow)
	wnd, err := frenyard.CreateBoundWindow("CCUpdaterUI", true, design.ThemeBackground, slideContainer)
	if err != nil {
		panic(err)
	}
	// Ok, now get it ready.
	app := (&upApplication{
		slideContainer: slideContainer,
		window: wnd,
		upQueued: make(chan func(), 16),
	})
	app.ShowGameFinderPreface()
	// Started!
	frenyard.GlobalBackend.Run(func(frameTime float64) {
		select {
			case fn := <- app.upQueued:
				fn()
			default:
		}
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

	testButtonWrapper := design.ButtonOkAction("OK", ok)

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
				Shrink:   1,
				Order:    0,
			},
		},
	})

	frenyard.InitUILayoutProxy(elem, titleAndThenBody)
	return elem
}
func (dialog *upTextPanel) FyETick(seconds float64) {
	dialog.UILayoutProxy.FyETick(seconds)
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

func upShowFailureToFindGameDialog(app *upApplication) {
	slideNum := 0
	var slideFn func ()
	slideFn = func () {
		app.slideContainer.TransitionTo(newUpTextPanelPtr(fmt.Sprintf("Slide %v\tHello from slide %v\nIt's nice here", slideNum, slideNum), slideFn), 1.0, slideNum & 2 == 0, slideNum & 1 == 0)
		slideNum++
	}
	app.slideContainer.TransitionTo(newUpTextPanelPtr("Test yup YUP a test\tMeep", slideFn), 0, false, false)
}
