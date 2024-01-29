package screens

import (
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
)

var ScreenButtons []framework.FlexboxSlot

func SetupButtons() {
	ScreenButtons = []framework.FlexboxSlot{
		{
			Element: framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
				DirVertical: false,
				Slots: []framework.FlexboxSlot{
					{
						Element: framework.NewUIOverlayContainerPtr(design.ScrollboxExterior, []framework.UILayoutElement{
							framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
								DirVertical: true,
								Slots: []framework.FlexboxSlot{
									{
										Grow: 1,
										Element: design.ButtonBar([]framework.UILayoutElement{
											design.ButtonAction(design.ThemeOkActionButton, "Ok", func() {}),
											design.ButtonAction(design.ThemeRemoveActionButton, "Remove", func() {}),
											design.ButtonAction(design.ThemeUpdateActionButton, "Update", func() {}),
											design.ButtonAction(design.ThemeImpossibleActionButton, "Disabled", func() {}),
											design.ButtonAction(design.ThemePageActionButton, "Details", func() {}),
										}),
									},
								},
							}),
						}),
					},
					{
						Grow: 1,
					},
					{
						Element: framework.NewUIOverlayContainerPtr(design.ScrollboxExterior, []framework.UILayoutElement{
							framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
								DirVertical: true,
								Slots:       []framework.FlexboxSlot{},
							}),
						}),
					},
				},
			}),
		},
		{
			Grow: 1,
		},
	}
}
