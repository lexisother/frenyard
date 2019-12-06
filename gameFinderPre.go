package main

import (
	"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
)

func (app *upApplication) ResetWithGameLocation(save bool, location string) {
	app.gameInstance = nil
	app.config.GamePath = location
	if save {
		middle.WriteUpdaterConfig(app.config)
	}
	// Re-kick
	app.ShowGameFinderPreface()
}

func (app *upApplication) ShowGameFinderPreface() {
	var gameLocations []middle.GameLocation
	app.ShowWaiter("Starting...", func (progress func(string)) {
		progress("Preparing remote packages...")
		middle.GetRemotePackages()
		progress("Scanning local installation...")
		gi := ccmodupdater.NewGameInstance(app.config.GamePath)
		fmt.Printf("Doing preliminary check of %s\n", app.config.GamePath)
		lp, err := local.AllLocalPackagePlugins(gi)
		if err == nil {
			gi.LocalPlugins = lp
			_, hasCC := gi.Packages()["crosscode"]
			if hasCC {
				app.gameInstance = gi
				return
			}
			fmt.Printf("Game not present?\n")
		} else {
			fmt.Printf("Failed check: %s\n", err.Error())
		}
		progress("Not configured ; Autodetecting game locations...")
		gameLocations = middle.AutodetectGameLocations()
	}, func () {
		if app.gameInstance == nil {
			app.ShowGameFinderPrefaceInternal(gameLocations)
		} else {
			app.ShowPrimaryView()
		}
	})
}

func (app *upApplication) ShowGameFinderPrefaceInternal(locations []middle.GameLocation) {

	suggestSlots := []framework.FlexboxSlot{}
	for _, location := range locations {
		suggestSlots = append(suggestSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Icon: design.GameIconID,
				Text: "CrossCode " + location.Version,
				Subtext: location.Location,
				Click: func () {
					app.GSRightwards()
					app.ResetWithGameLocation(true, location.Location)
				},
			}),
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
					design.ButtonAction(design.ThemeOkActionButton, "LOCATE MANUALLY", func () {
						app.GSDownwards()
						app.ShowGameFinder(func () {
							app.GSUpwards()
							app.ShowGameFinderPrefaceInternal(locations)
						}, middle.GameFinderVFSPathDefault)
					}),
				}),
			},
		},
	})
	primary := design.LayoutDocument(design.Header{
		BackIcon: design.WarningIconID,
		Back: func () {
			app.GSLeftwards()
			app.ShowCredits(func () {
				app.GSRightwards()
				app.ShowGameFinderPrefaceInternal(locations)
			})
		},
		Title: "Welcome",
	}, content, true)
	app.Teleport(primary)
}
