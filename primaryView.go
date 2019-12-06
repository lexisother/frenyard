package main

import (
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
	
	warnings := middle.FindWarnings(app.gameInstance)
	if app.config.DevMode {
		warnings = append(warnings, middle.Warning{
			Text: "You are in developer mode! Go to the Credits (top-right button, 'Credits') to deactivate it.",
			Action: middle.NullActionWarningID,
		})
	}
	for _, v := range warnings {
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
	slots = append(slots, framework.FlexboxSlot{
		Element: framework.NewUITextboxPtr("", "search (NYI)", design.GlobalFont, 0xFFFFFFFF, 0xC0FFFFFF, 0x80FFFFFF, 0xFF000000, frenyard.Alignment2i{X: frenyard.AlignStart}),
	})
	
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
		
		if (!app.config.DevMode) && (typeCheck.Metadata().Type() != ccmodupdater.PackageTypeMod) {
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
		ForwardIcon: design.MenuIconID,
		Forward: func () {
			app.GSRightwards()
			app.ShowOptionsMenu(func () {
				app.GSLeftwards()
				app.ShowPrimaryView()
			})
		},
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: slots,
	}), true))
}
