package main

//go:generate go run ./design/data-compiler txt main credits

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"encoding/base64"
	"strings"
	"strconv"
)

var secretDeveloperModeCounter int = 0

func (app *upApplication) ShowCredits(back framework.ButtonBehavior) {
	listSlots := []framework.FlexboxSlot{}
	credits, _ := base64.StdEncoding.DecodeString(creditsB64)
	sections := strings.Split(string(credits), "~!~")
	for i := 0; i < len(sections); i += 2 {
		sectionIdx := i
		sectionName := sections[i]
		sectionText := sections[i + 1]
		listSlots = append(listSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Text: sections[i],
				Click: func () {
					text := ""
					if sectionIdx == 0 {
						secretDeveloperModeCounter++
						if secretDeveloperModeCounter < 4 {
							text = "developer mode will be toggled by going here " + strconv.Itoa(4 - secretDeveloperModeCounter) + " more times\n"
						} else {
							app.config.DevMode = !app.config.DevMode
							middle.WriteUpdaterConfig(app.config)
						}
						if !app.config.DevMode {
							text += "developer mode is disabled"
						} else {
							text += "developer mode is enabled"
						}
					}
					app.GSLeftwards()
					app.MessageBox(sectionName, sectionText + text, func () {
						app.GSRightwards()
						app.ShowCredits(back)
					})
				},
			}),
		})
	}
	listSlots = append(listSlots, framework.FlexboxSlot{
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
