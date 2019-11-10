package frenyard
import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// TextLayouterOptions contains the options for text layout.
type TextLayouterOptions struct {
	Text string
	Font font.Face
	// SIZE_UNLIMITED should be used if an axis should be unbounded.
	Limits Vec2i
}

// TextLayouterResult The results from text layouting.
type TextLayouterResult struct {
	Size Vec2i
	_formattedText []string
	_font font.Face
}

// TextLayouterRenderable The results in a renderable form.
type TextLayouterRenderable struct {
	Size Vec2i
	_lines []Texture
	_interline int32
}

func (tlr *TextLayouterResult) fyAppendLine(line string) {
	tlr._formattedText = append(tlr._formattedText, line)
}
func (tlr *TextLayouterResult) fyCalcSize() {
	size := Vec2i{}
	interLine := FontInterline(tlr._font)
	for k, v := range tlr._formattedText {
		size = size.Max(FontSize(tlr._font, v).Add(Vec2i{0, int32(k) * interLine}))
	}
	tlr.Size = size
}

// Draw draws the laid-out text to textures and creates a TextLayouterRenderable.
func (tlr *TextLayouterResult) Draw() TextLayouterRenderable {
	tex := make([]Texture, len(tlr._formattedText))
	for k, v := range tlr._formattedText {
		tex[k] = FontDraw(tlr._font, v)
	}
	rdr := TextLayouterRenderable{
		Size: tlr.Size,
		_lines: tex,
		_interline: FontInterline(tlr._font),
	}
	return rdr
}

// Draw actually draws the TextLayouterRenderable to the screen.
func (tlr *TextLayouterRenderable) Draw(r Renderer, pos Vec2i, colour uint32) {
	interLine := tlr._interline
	for _, v := range tlr._lines {
		r.TexRect(v, colour, Area2iOfSize(v.Size()), Area2iFromVecs(pos, v.Size()))
		pos.Y += interLine
	}
}

func fyTextLayouterBreakerNormal(text string, fnt font.Face, xLimit int32, wordwrap bool) TextLayouterResult {
	committedBuffer := TextLayouterResult{
		_font: fnt,
		_formattedText: []string{},
	}
	// To lower the CPU usage, this function has to engage directly with the font drawing interface.
	lineDrawer := font.Drawer{
		Face: fnt,
	}
	lineBuffer := ""
	positionAfterLastSpace := -1
	// Do be alerted! This retrieves runes, but slices use byte indexes.
	// Used to get the start of the iterator to reset lineDrawer
	carriageReturn := false
	for len(text) > 0 {
		currentBody := text
		text = ""
		for charStart, char := range currentBody {
			if char == '\n' {
				// Newline, commit buffer
				committedBuffer.fyAppendLine(lineBuffer)
				lineBuffer = ""
				positionAfterLastSpace = -1
				carriageReturn = true
				continue
			}
			if carriageReturn {
				lineDrawer.Dot = fixed.Point26_6{}
				_, advance := lineDrawer.BoundString(lineBuffer)
				lineDrawer.Dot = lineDrawer.Dot.Add(fixed.Point26_6{X: advance})
			}
			// Insert character
			bound, advance := lineDrawer.BoundString(string(char))
			lineDrawer.Dot = lineDrawer.Dot.Add(fixed.Point26_6{X: advance})
			if int32(bound.Max.X.Ceil()) >= xLimit {
				// Check for cut position. Disabling wordwrap prevents these from ever being found.
				// Notably if the cut position is 0, that means it's probably a wordwrap start space.
				if positionAfterLastSpace > 0 {
					// Stitch together things. This gets very complicated because of the range iterator.
					text = lineBuffer[positionAfterLastSpace:] + currentBody[charStart:]
					committedBuffer.fyAppendLine(lineBuffer[:positionAfterLastSpace])
					lineBuffer = " "
					carriageReturn = true
					positionAfterLastSpace = -1
					// So with that done, use a break to reset the iterator.
					break
				} else {
					// Character break. Much simpler, doesn't even require interrupting the stream.
					// If the character actually failed here we'd be doomed anyway, so place it now, too.
					if lineBuffer == " " {
						committedBuffer.fyAppendLine(string(char))
						lineBuffer = ""
						carriageReturn = true
					} else {
						committedBuffer.fyAppendLine(lineBuffer)
						lineBuffer = string(char)
						carriageReturn = true
					}
				}
			} else {
				lineBuffer += string(char)
				if char == ' ' && wordwrap {
					// Successfully inserting space, let's note that
					positionAfterLastSpace = len(lineBuffer)
				}
			}
		}
	}
	if lineBuffer != "" {
		committedBuffer.fyAppendLine(lineBuffer)
	}
	committedBuffer.fyCalcSize()
	return committedBuffer
}

// TheOneTextLayouterToRuleThemAll lays out text with wrapping limits and other such constraints.
func TheOneTextLayouterToRuleThemAll(opts TextLayouterOptions) TextLayouterResult {
	brokenText := fyTextLayouterBreakerNormal(opts.Text, opts.Font, opts.Limits.X, true)
	if brokenText.Size.Y >= opts.Limits.Y {
		brokenText = fyTextLayouterBreakerNormal(opts.Text, opts.Font, opts.Limits.X, false)
	}
	return brokenText
}
