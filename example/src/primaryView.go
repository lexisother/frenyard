package src

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/example/screens"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func ChangeScreen(holder *framework.UIFlexboxContainer, screen []framework.FlexboxSlot) {
	// Postmortem: I wasted thirty to fourty minutes of my life trying to figure
	// out why `holder.ThisUIPanelDetails.SetContent` wasn't working. For some
	// reason I didn't realise, despite having looked at the function a million
	// times during debugging, that `UIFlexboxContainer.SetContent` is a
	// function.
	holder.SetContent(framework.FlexboxContainer{
		DirVertical: true,
		Slots:       screen,
	})
}

func (app *UpApplication) ShowPrimaryView() {
	screenHolder := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			{
				Grow: 1,
			},
			{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("No UI component selected!", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
			},
			{
				Grow: 1,
			},
		},
	})

	slots := []framework.FlexboxSlot{
		{
			Grow:  1,
			Basis: 3,
			Element: framework.NewUIOverlayContainerPtr(design.ScrollboxExterior, []framework.UILayoutElement{
				framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
					DirVertical: true,
					Slots: []framework.FlexboxSlot{
						{
							Element: design.ListItem(design.ListItemDetails{
								Text:    "Buttons",
								Subtext: "A simple display of button types",
								Click: func() {
									screens.SetupButtons()
									ChangeScreen(screenHolder, screens.ScreenButtons)
								},
							}),
						},
						{
							Element: design.ListItem(design.ListItemDetails{
								Text:    "haiih screen twoo",
								Subtext: "shoowwww!!",
								Click: func() {
									screens.SetupTwo()
									ChangeScreen(screenHolder, screens.ScreenTwo)
								},
							}),
						},
					},
				}),
			}),
		},
	}

	slots = append(slots, []framework.FlexboxSlot{
		{
			Grow: 1,
		},
		{
			Element: screenHolder,
		},
		{
			Grow: 1,
		},
	}...)

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "UI Playground",
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: false,
		Slots:       slots,
	}), false))

}
