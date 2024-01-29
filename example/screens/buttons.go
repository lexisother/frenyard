package screens

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

var ScreenButtons []framework.FlexboxSlot

func SetupButtons() {
	ScreenButtons = []framework.FlexboxSlot{
		{
			Grow: 1,
		},
		{
			Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("Simple buttons with themed colours.", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
		},
		{
			Basis: design.SizeMarginAroundEverything,
		},
		{
			Element: framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
				DirVertical: false,
				Slots: []framework.FlexboxSlot{
					{
						Grow: 1,
					},
					{
						Element: framework.NewUIOverlayContainerPtr(design.ScrollboxExterior, []framework.UILayoutElement{
							framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
								DirVertical: true,
								Slots: []framework.FlexboxSlot{
									{
										Element: design.ButtonAction(design.ThemeOkActionButton, "Ok", func() {}),
									},
									{
										Grow: 1,
									},
									{
										Element: design.ButtonAction(design.ThemeRemoveActionButton, "Remove", func() {}),
									},
									{
										Grow: 1,
									},
									{
										Element: design.ButtonAction(design.ThemeUpdateActionButton, "Update", func() {}),
									},
									{
										Grow: 1,
									},
									{
										Element: design.ButtonAction(design.ThemeImpossibleActionButton, "Disabled", func() {}),
									},
									{
										Grow: 1,
									},
									{
										Element: design.ButtonAction(design.ThemePageActionButton, "Details", func() {}),
									},
								},
							}),
						}),
					},
					{
						Grow: 1,
					},
				},
			}),
		},
		{
			Grow: 1,
		},
	}
}
