package main

import (
	//"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

func (app *upApplication) ShowGameFinderPreface() {
	foundInstallsScroller := design.ScrollboxV(frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("EXAMPLE INSTALLATION 1 [OPEN] [OK]\nEXAMPLE INSTALLATION 2 [OPEN] [OK]\nEXAMPLE INSTALLATION 1 [OPEN] [OK]\nEXAMPLE INSTALLATION 2 [OPEN] [OK]\nEXAMPLE INSTALLATION 1 [OPEN] [OK]\nEXAMPLE INSTALLATION 2 [OPEN] [OK]\n", design.GlobalFont), design.ThemePlaceholder, 0, frenyard.Alignment2i{}))
	
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
				Element: frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("If the installation you'd like to install mods to isn't shown here, you can locate it manually below.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			frenyard.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			frenyard.FlexboxSlot{
				Element: frenyard.NewUILabelPtr(frenyard.NewTextTypeChunk("[LOCATE MANUALLY]", design.GlobalFont), design.ThemePlaceholder, 0, frenyard.Alignment2i{}),
			},
		},
	})
	primary := design.LayoutDocument(design.Header{
		Title: "Welcome",
	}, content, true)
	app.slideContainer.TransitionTo(primary, 1.0, true, true)
	
}
