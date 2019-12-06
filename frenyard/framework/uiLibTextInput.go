package framework

import (
	"golang.org/x/image/font"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

// UITextbox is a textbox. NOT YET IMPLEMENTED, OF COURSE...
type UITextbox struct {
	UILayoutProxy
	_text string
	_face font.Face
	_primaryColour uint32
	_suggestionColour uint32
	_label *UILabel
}

// NewUITextboxPtr creates a UITextbox.
func NewUITextboxPtr(text string, face font.Face, primaryColour uint32, suggestionColour uint32, align frenyard.Alignment2i) *UITextbox {
	textbox := &UITextbox{
		_text: text,
		_face: face,
		_primaryColour: primaryColour,
		_suggestionColour: suggestionColour,
		_label: NewUILabelPtr(integration.NewCompoundTypeChunk([]integration.TypeChunk{}), 0xFFFFFFFF, 0, align),
	}
	InitUILayoutProxy(textbox, textbox._label)
	textbox.rebuild()
	return textbox
}

func (tb *UITextbox) rebuild() {
	tb._label.SetText(integration.NewColouredTextTypeChunk(tb._text, tb._face, tb._primaryColour))
}

// FyENormalEvent overrides UILayoutProxy.FyENormalEvent
func (tb *UITextbox) FyENormalEvent(ne frenyard.NormalEvent) {
	switch ev := ne.(type) {
		case FocusEvent:
			if ev.Focused {
			}
	}
	tb.UILayoutProxy.FyENormalEvent(ne)
}
