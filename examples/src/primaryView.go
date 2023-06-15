package src

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/examples/screens"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func (app *UpApplication) ShowPrimaryView(rightSide ...framework.FlexboxSlot) {
	// FOR FUTURE REFERENCE! I might be able to achieve some sort of rerendering mechanism by separating the UI
	// components that matter into their own variables... Unsure if that works though.
	// Rambles about it: <https://discord.com/channels/382339402338402315/382339402338402317/1118873796582060162>
	//e := framework.NewUIOverlayContainerPtr(design.ScrollboxExterior, []framework.UILayoutElement{})
	//e.ThisUIPanelDetails.SetContent([]framework.PanelFixedElement{
	//	{
	//		Pos:     frenyard.Vec2i{},
	//		Visible: false,
	//		Locked:  false,
	//		Element: nil,
	//	},
	//})
	//framework.InitUILayoutElementComponent(e)
	//e.ThisUILayoutElementComponentDetails.ContentChanged()

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
								Text:    "helo screen one",
								Subtext: "show me the first screen!!",
								Click: func() {
									app.GSInstant()
									screens.SetupOne()
									app.ShowPrimaryView(screens.ScreenOne...)
								},
							}),
						},
						{
							Element: design.ListItem(design.ListItemDetails{
								Text:    "haiih screen twoo",
								Subtext: "shoowwww!!",
								Click: func() {
									app.GSInstant()
									screens.SetupTwo()
									app.ShowPrimaryView(screens.ScreenTwo...)
								},
							}),
						},
					},
				}),
			}),
		},
	}
	if rightSide == nil {
		slots = append(slots, []framework.FlexboxSlot{
			{
				Grow: 1,
			},
			{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("No UI component selected!", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
			},
			{
				Grow: 1,
			},
		}...)
	} else {
		slots = append(slots, rightSide...)
	}

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "UI Playground",
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: false,
		Slots:       slots,
	}), false))

}
