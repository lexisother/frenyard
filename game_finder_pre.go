package main

import (
	//"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
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

	suggestSlots := []frenyard.FlexboxSlot{}
	for _, location := range locations {
		suggestSlots = append(suggestSlots, frenyard.FlexboxSlot{
			Element: design.ListItem(design.GameIconID, "CrossCode " + location.Version, location.Location),
			RespectMinimumSize: true,
		})
	}
	// Space-taker to prevent wrongly scaled list items
	suggestSlots = append(suggestSlots, frenyard.FlexboxSlot{
		Grow: 1,
		Shrink: 0,
	})

	foundInstallsScroller := design.ScrollboxV(frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		DirVertical: true,
		WrapMode: frenyard.FlexboxWrapModeNone,
		Slots: suggestSlots,
	}))
	
	content := frenyard.NewUIFlexboxContainerPtr(frenyard.FlexboxContainer{
		DirVertical: true,
		Slots: []frenyard.FlexboxSlot{
			frenyard.FlexboxSlot{
				Element: frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("Welcome to the unofficial CrossCode mod updater UI.\nBefore we begin, you need to indicate which CrossCode installation you want to install mods to.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			frenyard.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			frenyard.FlexboxSlot{
				Element: foundInstallsScroller,
				Grow: 1,
				Shrink: 1,
				RespectMinimumSize: true,
			},
			frenyard.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			frenyard.FlexboxSlot{
				Element: frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("If the installation you'd like to install mods to isn't shown here, you can locate it manually.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			frenyard.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			frenyard.FlexboxSlot{
				Element: design.ButtonBar([]frenyard.UILayoutElement{
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
