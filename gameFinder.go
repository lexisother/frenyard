package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"path/filepath"
)

func (app *upApplication) ShowGameFinder(back framework.ButtonBehavior, vfsPath string) {
	var vfsList []middle.GameLocation
	
	app.ShowWaiter("Reading", func (progress func(string)) {
		progress("Scanning to find all of the contents of in:\n" + vfsPath + "\nIf this includes CD/DVD drives or network partitions, this may take a while.")
		vfsList = middle.GameFinderVFSList(vfsPath)
	}, func () {
		slots := []framework.FlexboxSlot{}
		
		for _, v := range vfsList {
			thisLocation := v.Location
			ild := design.ListItemDetails{
				Icon: design.DirectoryIconID,
				Text: filepath.Base(thisLocation),
			}
			ild.Click = func () {
				app.GSRightwards()
				app.ShowGameFinder(func () {
					app.GSLeftwards()
					app.ShowGameFinder(back, vfsPath)
				}, thisLocation)
			}
			if v.Valid {
				ild.Click = func () {
					app.GSRightwards()
					app.ResetWithGameLocation(true, thisLocation)
				}
				ild.Text = "CrossCode " + v.Version
				ild.Subtext = thisLocation
				ild.Icon = design.GameIconID
			} else if v.Drive != "" {
				ild.Text = v.Drive
				ild.Subtext = v.Location
				ild.Icon = design.DriveIconID
			}
			item := design.ListItem(ild)
			slots = append(slots, framework.FlexboxSlot{
				Element: item,
				RespectMinimumSize: true,
			})
		}
		
		slots = append(slots, framework.FlexboxSlot{
			Grow: 1,
		})

		primary := design.LayoutDocument(design.Header{
			Back: back,
			Title: "Enter CrossCode's location",
		}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: true,
			Slots: slots,
		}), true)
		app.Teleport(primary)
	})
}
