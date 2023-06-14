package src

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func (app *UpApplication) ShowPrimaryView(newSlots ...[]framework.FlexboxSlot) {
	var slots []framework.FlexboxSlot
	if newSlots == nil {
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
	}

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "UI Playground",
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: false,
		Slots:       slots,
	}), false))

}
