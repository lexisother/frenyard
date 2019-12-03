package main

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/middle"
)

// The Primary View
func (app *upApplication) ShowPrimaryView() {
	temp := framework.NewUILabelPtr(integration.NewTextTypeChunk("[AN ACTUAL UI HERE]", design.GlobalFont), design.ThemePlaceholder, 0, frenyard.Alignment2i{})
	app.slideContainer.TransitionTo(framework.SlideTransition{
		Element: design.LayoutDocument(design.Header{
			Title: "CCUpdaterUI",
			Back: func () {
				app.ResetWithGameLocation(middle.GameFinderVFSPathDefault)
			},
			BackIcon: design.GameIconID,
		}, temp, true),
		Length: 1.0,
	})
}
