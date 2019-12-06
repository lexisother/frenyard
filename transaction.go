package main

import (
	"os"
	"time"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
)

// PerformTransaction performs a transaction, showing UI for it as well.
func (app *upApplication) PerformTransaction(back framework.ButtonBehavior, tx ccmodupdater.PackageTX) {
	ctx := ccmodupdater.PackageTXContext{
		LocalPackages: app.gameInstance.Packages(),
		RemotePackages: middle.GetRemotePackages(),
	}
	solutions, err := ctx.Solve(tx)
	if err != nil {
		app.MessageBox("Error", err.Error(), back)
		return
	}
	if len(solutions) != 1 {
		app.MessageBox("Issue", "Multiple solutions were given to a dependency problem.\nThis shouldn't occur in the present iteration.\nPlease install dependencies yourself to constrain the system's choices.", back)
		return
	}
	// It begins...
	log := "-- Log started at " + time.Now().Format(time.RFC1123) + " (ccmodupdater.log) --"
	app.ShowWaiter("Working...", func (progress func(string)) {
		err = ctx.Perform(app.gameInstance, solutions[0], func (pkg string, pre bool, rm bool, add bool) {
			if !pre {
				return
			}
			if rm && add {
				log += "\nUpgrading " + pkg + "..."
				progress(log)
			} else if rm {
				log += "\nRemoving " + pkg + "..."
				progress(log)
			} else if add {
				log += "\nInstalling " + pkg + "..."
				progress(log)
			}
		})
		if err != nil {
			log += "\n-- Error --\n" + err.Error()
		} else {
			log += "\n-- Complete --"
		}
	}, func () {
		cfgFile, err := os.OpenFile("ccmodupdater.log", os.O_WRONLY | os.O_CREATE, os.ModePerm)
		if err == nil {
			// Oh well
			cfgFile.WriteString(log + "\n")
		}
		cfgFile.Close()
		app.GSInstant()
		app.MessageBox("Report", log, back)
	});
}
