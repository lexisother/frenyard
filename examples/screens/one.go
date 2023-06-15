package screens

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

var ScreenOne []framework.FlexboxSlot

func SetupOne() {
	ScreenOne = []framework.FlexboxSlot{
		{
			Grow: 1,
		},
		{
			Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("haiiii I am screen one!!", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
		},
		{
			Grow: 1,
		},
	}
}
