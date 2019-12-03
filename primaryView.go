package main

import (
	"time"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"sort"
)

// ShowPrimaryView shows the "Primary View" (the mod list right now)
func (app *upApplication) ShowPrimaryView() {
	slots := []framework.FlexboxSlot{}
	
	for _, v := range middle.FindWarnings(app.gameInstance) {
		fixAction := framework.ButtonBehavior(nil)
		if v.Action == middle.InstallOrUpdatePackageWarningID {
			pkgID := v.Parameter
			fixAction = func () {
				app.GSRightwards()
				app.ShowPackageView(func () {
					app.GSLeftwards()
					app.ShowPrimaryView()
				}, pkgID)
			}
		}
		slots = append(slots, framework.FlexboxSlot{
			Element: design.InformationPanel(design.InformationPanelDetails{
				Text: v.Text,
				ActionText: "FIX",
				Action: fixAction,
			}),
		})
	}
	
	// Ok, let's get all the packages in a nice row
	localPackages := app.gameInstance.Packages()
	remotePackages := middle.GetRemotePackages()
	packageSet := make(map[string]bool)
	packageList := []string{}
	for k := range localPackages {
		packageSet[k] = true
	}
	for k := range remotePackages {
		packageSet[k] = true
	}
	for k := range packageSet {
		packageList = append(packageList, k)
	}
	sort.Sort(sort.StringSlice(packageList))
	// Actually build the UI now!
	for _, pkgID := range packageList {
		local := localPackages[pkgID]
		remote := remotePackages[pkgID]
		latest := middle.GetLatestOf(local, remote)
		var typeCheck ccmodupdater.Package = local
		if typeCheck == nil {
			typeCheck = remote
		}
		
		if typeCheck.Metadata().Type() != ccmodupdater.PackageTypeMod {
			continue
		}
		
		status := "unable to comprehend status"
		if local != nil && remote != nil {
			lmv := local.Metadata().Version()
			rmv := remote.Metadata().Version()
			if local.Metadata().Version().GreaterThan(remote.Metadata().Version()) {
				status = lmv.Original() + " installed (local development build, " + rmv.Original() + " remote)"
			} else if remote.Metadata().Version().GreaterThan(local.Metadata().Version()) {
				status = lmv.Original() + " installed (out of date, " + rmv.Original() + " available)"
			} else {
				status = lmv.Original() + " (up to date)"
			}
		} else if local != nil {
			status = latest.Metadata().Version().Original() + " installed (no remote copy available)"
		} else if remote != nil {
			status = "not installed: " + latest.Metadata().Version().Original() + " available"
		}
		pkgIDLocal := pkgID
		slots = append(slots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Icon: middle.PackageIcon(latest),
				Text: latest.Metadata().HumanName(),
				Subtext: status,
				Click: func () {
					app.GSRightwards()
					app.ShowPackageView(func () {
						app.GSLeftwards()
						app.ShowPrimaryView()
					}, pkgIDLocal)
				},
			}),
		})
	}

	slots = append(slots, framework.FlexboxSlot{
		Grow: 1,
	})
	
	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Mods",
		Back: func () {
			app.GSLeftwards()
			app.ResetWithGameLocation(false, middle.GameFinderVFSPathDefault)
		},
		BackIcon: design.GameIconID,
		Forward: func () {
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
					frenyard.TargetFrameTime = 1
					proc.Wait()
					frenyard.TargetFrameTime = backupFrameTime
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
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: slots,
	}), true))
}
