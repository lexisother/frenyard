package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"path/filepath"
	"sort"
)

func (app *upApplication) ShowPackedModFinder(back framework.ButtonBehavior, success func (path string), vfsPath string) {
	var vfsList []middle.PackedModLocation
	
	app.ShowWaiter("Reading", func (progress func(string)) {
		progress("Scanning to find all of the contents of in:\n" + vfsPath + "\nIf this includes CD/DVD drives or network partitions, this may take a while.")
		vfsList = middle.PackedModFinderVFSList(vfsPath)
	}, func () {
		items := []design.ListItemDetails{}
		
		for _, v := range vfsList {
			thisLocation := v.Location
			ild := design.ListItemDetails{
				Icon: design.DirectoryIconID,
				Text: filepath.Base(thisLocation),
			}
			ild.Click = func () {
				app.GSRightwards()
				app.ShowPackedModFinder(func () {
					app.GSLeftwards()
					app.ShowPackedModFinder(back, success, vfsPath)
				}, success, thisLocation)
			}
			if v.Valid {
				ild.Click = func () {
					app.GSRightwards()
					success(thisLocation)
				}
				ild.Text = v.Metadata.HumanName()
				ild.Subtext = thisLocation
				ild.Icon = design.ModIconID
			} else if v.Drive != "" {
				ild.Text = v.Drive
				ild.Subtext = v.Location
				ild.Icon = design.DriveIconID
			}
			items = append(items, ild)
		}

		sort.Sort(design.SortListItemDetails(items))
		primary := design.LayoutDocument(design.Header{
			Back: back,
			Title: "Enter CrossCode's location",
		}, design.NewUISearchBoxPtr("Directory name...", items), true)
		app.Teleport(primary)
	})
}
