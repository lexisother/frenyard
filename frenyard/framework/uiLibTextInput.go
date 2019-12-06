package framework

import (
	"golang.org/x/image/font"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

// UITextbox is a textbox.
type UITextbox struct {
	UILayoutProxy
	
	// Called on rebuild (this includes, though isn't strictly limited to, changes)
	OnRebuild func()
	// Called after some time without anything happening
	OnStall func()
	// Called on enter
	OnConfirm func()
	
	_open bool
	_caretBlinker float64
	_stallTimer float64
	
	_hint string

	_textPre string
	_suggesting bool
	_textSuggestion1 string
	_textSuggestion2 string
	_textSuggestion3 string
	_textPost string

	_face font.Face
	_primaryColour uint32
	_suggestionColour uint32
	_hintColour uint32
	_label *UILabel
	_translation frenyard.Vec2i
	_window frenyard.Window
}

// NewUITextboxPtr creates a UITextbox.
func NewUITextboxPtr(text string, hint string, face font.Face, primaryColour uint32, suggestionColour uint32, hintColour uint32, backgroundColour uint32, align frenyard.Alignment2i) *UITextbox {
	textbox := &UITextbox{
		_textPre: text,
		_hint: hint,
		_face: face,
		_primaryColour: primaryColour,
		_suggestionColour: suggestionColour,
		_hintColour: hintColour,
		_label: NewUILabelPtr(integration.NewCompoundTypeChunk([]integration.TypeChunk{}), 0xFFFFFFFF, backgroundColour, align),
	}
	InitUILayoutProxy(textbox, textbox._label)
	textbox.rebuild()
	return textbox
}

// Text gets the text in the textbox.
func (tb *UITextbox) Text() string {
	return tb._textPre + tb._textPost
}

// SetText sets the text in the textbox.
func (tb *UITextbox) SetText(text string) {
	tb._textPre = text
	tb._textPost = ""
	tb.rebuild()
}

// FyETick overrides UILayoutProxy.FyETick
func (tb *UITextbox) FyETick(delta float64) {
	tb._stallTimer += delta
	if tb._stallTimer > 0.25 {
		tb._stallTimer = 0
		if tb.OnStall != nil {
			tb.OnStall()
		}
	}
	if tb._open {
		oldCB := tb._caretBlinker
		tb._caretBlinker += delta
		if oldCB < 1 && tb._caretBlinker >= 1 {
			tb.rebuild()
		} else if tb._caretBlinker >= 2 {
			tb._caretBlinker = 0
			tb.rebuild()
		}
	}
	tb.UILayoutProxy.FyETick(delta)
}

func (tb *UITextbox) rebuild() {
	tb._stallTimer = 0
	if !tb._open {
		if tb._textPre == "" && tb._textPost == "" {
			tb._label.SetText(integration.NewCompoundTypeChunk([]integration.TypeChunk{
				integration.NewColouredTextTypeChunk(tb._hint, tb._face, tb._hintColour),
			}))
		} else {
			tb._label.SetText(integration.NewCompoundTypeChunk([]integration.TypeChunk{
				integration.NewColouredTextTypeChunk(tb._textPre + tb._textPost, tb._face, tb._primaryColour),
			}))
		}
	} else if !tb._suggesting {
		caretColour := tb._primaryColour
		if tb._caretBlinker > 1 {
			caretColour = 0
		}
		tb._label.SetText(integration.NewCompoundTypeChunk([]integration.TypeChunk{
			integration.NewColouredTextTypeChunk(tb._textPre, tb._face, tb._primaryColour),
			integration.NewColouredTextTypeChunk("_", tb._face, caretColour),
			integration.NewColouredTextTypeChunk(tb._textPost, tb._face, tb._primaryColour),
		}))
	} else {
		tb._label.SetText(integration.NewCompoundTypeChunk([]integration.TypeChunk{
			integration.NewColouredTextTypeChunk(tb._textPre, tb._face, tb._primaryColour),
			integration.NewColouredTextTypeChunk(tb._textSuggestion1, tb._face, tb._suggestionColour),
			integration.NewUnderlineTypeChunk(integration.NewColouredTextTypeChunk(tb._textSuggestion2, tb._face, tb._suggestionColour), tb._suggestionColour),
			integration.NewColouredTextTypeChunk(tb._textSuggestion3, tb._face, tb._suggestionColour),
			integration.NewColouredTextTypeChunk(tb._textPost, tb._face, tb._primaryColour),
		}))
	}
	if tb.OnRebuild != nil {
		tb.OnRebuild()
	}
}

// FyEDraw overrides UILayoutProxy.FyEDraw
func (tb *UITextbox) FyEDraw(target frenyard.Renderer, under bool) {
	tb._translation = target.Translation()
	tb.UILayoutProxy.FyEDraw(target, under)
}

// FyENormalEvent overrides UILayoutProxy.FyENormalEvent
func (tb *UITextbox) FyENormalEvent(ne frenyard.NormalEvent) {
	switch ev := ne.(type) {
		case EnterWindowEvent:
			if tb._window != nil && tb._window != ev.Window {
				if tb._window.TextInput() == tb {
					tb._window.SetTextInput(nil)
				}
			}
			tb._window = ev.Window
		case FocusEvent:
			if tb._window != nil {
				if ev.Focused {
					tb._window.SetTextInput(tb)
				} else {
					if tb._window.TextInput() == tb {
						tb._window.SetTextInput(nil)
					}
				}
			}
		case frenyard.KeyEvent:
			if ev.Pressed && !tb._suggesting {
				if (ev.Keycode == 13) || (ev.Keycode == 10) {
					// confirm
					if tb.OnConfirm != nil {
						tb.OnConfirm()
					}
				} else if (ev.Keycode == 1073741898) || (ev.Keycode == 1073741919) {
					// Home
					tb._textPost = tb._textPre + tb._textPost
					tb._textPre = ""
					tb.rebuild()
				} else if (ev.Keycode == 1073741901) || (ev.Keycode == 1073741913) {
					// End
					tb._textPre = tb._textPre + tb._textPost
					tb._textPost = ""
					tb.rebuild()
				} else if (ev.Keycode == 1073741903) || (ev.Keycode == 127) {
					// Right Arrow / Delete
					moveRange := 1
					runes := []rune(tb._textPost)
					if len(runes) >= moveRange {
						transfer := string(runes[:moveRange])
						tb._textPost = string(runes[moveRange:])
						if ev.Keycode == 1073741903 {
							tb._textPre += transfer
						}
						tb.rebuild()
					}
				} else if (ev.Keycode == 1073741904) || (ev.Keycode == 8) {
					// Left Arrow / Backspace
					moveRange := 1
					runes := []rune(tb._textPre)
					if len(runes) >= moveRange {
						transfer := string(runes[len(runes) - moveRange:])
						tb._textPre = string(runes[:len(runes) - moveRange])
						if ev.Keycode == 1073741904 {
							tb._textPost = transfer + tb._textPost
						}
						tb.rebuild()
					}
				}
			}
	}
	tb.UILayoutProxy.FyENormalEvent(ne)
}

// FyTOpen implements TextInput.FyTOpen
func (tb *UITextbox) FyTOpen() {
	tb._open = true
	tb._textSuggestion1 = ""
	tb._textSuggestion2 = ""
	tb._textSuggestion3 = ""
	tb._suggesting = false
	tb.rebuild()
}
// FyTClose implements TextInput.FyTClose
func (tb *UITextbox) FyTClose() {
	tb._open = false
	tb._textSuggestion1 = ""
	tb._textSuggestion2 = ""
	tb._textSuggestion3 = ""
	tb._suggesting = false
	tb.rebuild()
}
// FyTEditing implements TextInput.FyTEditing
func (tb *UITextbox) FyTEditing(text string, start int, length int) {
	// Ok, so as it turns out, it's entirely possible for start/length to do stupid things.
	// Stupid things that cause runtime errors in this code.
	// If anyone REALLY REALLY needs this working, figure it out yourself...
	tb._textSuggestion1 = ""
	tb._textSuggestion2 = text
	tb._textSuggestion3 = ""
	tb._suggesting = text != ""
	tb.rebuild()
}
// FyTInput implements TextInput.FyTInput
func (tb *UITextbox) FyTInput(text string) {
	tb._textPre += text
	tb.FyTOpen()
}

// FyTArea implements TextInput.FyTArea
func (tb *UITextbox) FyTArea() frenyard.Area2i {
	return frenyard.Area2iFromVecs(tb._translation, tb.FyESize())
}
