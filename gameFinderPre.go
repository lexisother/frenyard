package main

import (
	//"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
	"github.com/20kdc/CCUpdaterUI/middle"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

func (app *upApplication) ShowGameFinderPreface() {
	var gameLocations []middle.GameLocation
	app.ShowWaiter("Finding CrossCode...", func (progress func(string)) {
		progress("Autodetecting game locations...")
		gameLocations = middle.AutodetectGameLocations()
	}, func () {
		app.ShowGameFinderPrefaceInternal(gameLocations)
	})
}

func (app *upApplication) ShowGameFinderPrefaceInternal(locations []middle.GameLocation) {

	suggestSlots := []framework.FlexboxSlot{}
	for _, location := range locations {
		suggestSlots = append(suggestSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.GameIconID, "CrossCode " + location.Version, location.Location),
			RespectMinimumSize: true,
		})
	}
	// Space-taker to prevent wrongly scaled list items
	suggestSlots = append(suggestSlots, framework.FlexboxSlot{
		Grow: 1,
		Shrink: 0,
	})

	foundInstallsScroller := design.ScrollboxV(framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		WrapMode: framework.FlexboxWrapModeNone,
		Slots: suggestSlots,
	}))
	
	content := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("Welcome to the unofficial CrossCode mod updater UI.\nBefore we begin, you need to indicate which CrossCode installation you want to install mods to.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: foundInstallsScroller,
				Grow: 1,
				Shrink: 1,
				RespectMinimumSize: true,
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("If the installation you'd like to install mods to isn't shown here, you can locate it manually.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: design.ButtonBar([]framework.UILayoutElement{
					design.ButtonOkAction("LOCATE MANUALLY", func () {
						app.ShowGameFinder()
					}),
				}),
			},
		},
	})
	primary := design.LayoutDocument(design.Header{
		Title: "Welcome",
	}, content, true)
	app.slideContainer.TransitionTo(primary, 1.0, false, false)
	
}
