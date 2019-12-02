package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

func (app *upApplication) ShowWaiter(text string, a func(func(string)), b func()) {
	label := framework.NewUILabelPtr(integration.NewTextTypeChunk("", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{})
	app.slideContainer.TransitionTo(design.LayoutDocument(design.Header{
		Title: text,
	}, label, false), 1.0, false, false)
	go func () {
		a(func (text string) {
			app.upQueued <- func () {
				label.SetText(integration.NewTextTypeChunk(text, design.GlobalFont))
			}
		})
		app.upQueued <- b
	}()
}
