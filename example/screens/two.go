package screens

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

var ScreenTwo []framework.FlexboxSlot

func SetupTwo() {
	ScreenTwo = []framework.FlexboxSlot{
		{
			Grow: 1,
		},
		{
			Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("heyyyy I am screen twoo!!!!", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
		},
		{
			Grow: 1,
		},
	}
}
