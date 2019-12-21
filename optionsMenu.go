package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"time"
)

// ShowOptionsMenu shows the options menu (run game, credits)
func (app *upApplication) ShowOptionsMenu(back framework.ButtonBehavior) {
	backHere := func () {
		app.GSLeftwards()
		app.ShowOptionsMenu(back)
	}
	listSlots := []framework.FlexboxSlot{
		{
			Element: design.ListItem(design.ListItemDetails{
				Text: "Run Game",
				Subtext: "Attempts to run the game",
				Click: func () {
					backupFrameTime := frenyard.TargetFrameTime
					app.GSRightwards()
					app.ShowWaiter("Running...", func (progress func (string)) {
						progress("Trying to run game...")
						time.Sleep(time.Second * 1)
						proc, err := middle.Launch(app.gameInstance.Base())
						if err != nil {
							progress("Unable to launch CrossCode.\nIf on a Unix-like, try adding a 'run' script to the directory containing 'assets'.\nIf on Windows, ensure said directory contains nw.exe or CrossCode.exe for usage by the game.")
							time.Sleep(time.Second * 5)
						} else {
							progress("Game running...")
							app.upQueued <- func () {
								frenyard.TargetFrameTime = 1
							}
							proc.Wait()
							app.upQueued <- func () {
								frenyard.TargetFrameTime = backupFrameTime
							}
							// give the system time to 'calm down'
							time.Sleep(time.Second * 2)
						}
					}, func () {
						// make sure in case of threading shenanigans
						frenyard.TargetFrameTime = backupFrameTime
						app.GSLeftwards()
						app.ShowPrimaryView()
					})
				},
			}),
		},
		{
			Element: design.ListItem(design.ListItemDetails{
				Text: "Credits",
				Subtext: "See the various components that make up CCUpdaterUI",
				Click: func () {
					app.GSRightwards()
					app.ShowCredits(backHere)
				},
			}),
		},
		{
			Element: design.ListItem(design.ListItemDetails{
				Text: "Show Modloader",
				Subtext: "Show the installed modloader",
				Click: func () {
					app.GSRightwards()
					app.ShowPackageView(backHere, "ccloader")
				},
			}),
		},
		{
			Grow: 1,
		},
	}
	
	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Options",
		Back: back,
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: listSlots,
	}), true))
}
