package main

//go:generate go run ./design/data-compiler txt main credits

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"encoding/base64"
	"strings"
	"strconv"
	"runtime"
)

var secretDeveloperModeCounter int = 0

func (app *upApplication) ShowCredits(back framework.ButtonBehavior) {
	listSlots := []framework.FlexboxSlot{}
	credits, _ := base64.StdEncoding.DecodeString(creditsB64)
	sections := strings.Split(string(credits), "\n~!~\n")
	for i := 0; i < len(sections); i += 2 {
		sectionHeaderPieces := strings.Split(sections[i], "\n")
		sectionName := sectionHeaderPieces[0]
		sectionSubtext := "PARSER ERROR"
		if len(sectionHeaderPieces) == 2 {
			sectionSubtext = sectionHeaderPieces[1]
		}
		sectionText := "PARSER ERROR"
		if i != len(sections) - 1 {
			sectionText = sections[i + 1]
		}
		listSlots = append(listSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Text: sectionName,
				Subtext: sectionSubtext,
				Click: func () {
					app.GSRightwards()
					app.MessageBox(sectionName, sectionText, func () {
						app.GSLeftwards()
						app.ShowCredits(back)
					})
				},
			}),
		})
	}
	listSlots = append(listSlots, framework.FlexboxSlot{
		Element: design.ListItem(design.ListItemDetails{
			Text: "Build Information",
			Subtext: runtime.GOOS + " " +runtime.GOARCH + " " + runtime.Compiler + " " + runtime.Version(),
			Click: func () {
				text := ""
				secretDeveloperModeCounter++
				if secretDeveloperModeCounter < 4 {
						text = "developer mode will be toggled by going here " + strconv.Itoa(4 - secretDeveloperModeCounter) + " more times\n"
				} else {
					app.cachedPrimaryView = nil
					app.config.DevMode = !app.config.DevMode
					middle.WriteUpdaterConfig(app.config)
				}
				if !app.config.DevMode {
					text += "developer mode is disabled"
				} else {
					text += "developer mode is enabled"
				}
				app.GSRightwards()
				app.MessageBox("Dev", text, func () {
					app.GSLeftwards()
					app.ShowCredits(back)
				})
			},
		}),
	}, framework.FlexboxSlot{
		Grow: 1,
	})
	
	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Credits",
		Back: back,
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: listSlots,
	}), true))
}
