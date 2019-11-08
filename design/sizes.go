package design
import "github.com/20kdc/CCUpdaterUI/frenyard"

const SIZE_TITLE_HEIGHT = 40
const SIZE_MARGIN_AROUND_EVERYTHING = 16
const SIZE_TEXT_NUDGE = 4
func MarginBody() frenyard.Area2i {
	return frenyard.Area2i{
		X: frenyard.Area1i{Pos: -SIZE_MARGIN_AROUND_EVERYTHING, Size: SIZE_MARGIN_AROUND_EVERYTHING * 2},
		Y: frenyard.Area1i{Pos: -(SIZE_MARGIN_AROUND_EVERYTHING + SIZE_TEXT_NUDGE), Size: (SIZE_MARGIN_AROUND_EVERYTHING * 2) + SIZE_TEXT_NUDGE},
	}
}
