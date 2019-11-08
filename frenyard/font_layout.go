package frenyard
import (
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

/*
 * PLEASE use key-based options when using this!!! IT IS SUBJECT TO CHANGE!
 */
type TextLayouterOptions struct {
	Text string
	Font font.Face
	// SIZE_UNLIMITED should be used if an axis should be unbounded.
	Limits Vec2i
}

/*
 * The results, horrific as they may be.
 */
type TextLayouterResult struct {
	Size Vec2i
	_fy_TextLayouterResult_formattedText []string
	_fy_TextLayouterResult_font font.Face
}
type TextLayouterRenderable struct {
	Size Vec2i
	_fy_TextLayouterResult_lines []Texture
	_fy_TextLayouterResult_interline int32
}

func (tlr *TextLayouterResult) fyAppendLine(line string) {
	tlr._fy_TextLayouterResult_formattedText = append(tlr._fy_TextLayouterResult_formattedText, line)
}
func (tlr *TextLayouterResult) fyCalcSize() {
	size := Vec2i{}
	interLine := FontInterline(tlr._fy_TextLayouterResult_font)
	for k, v := range tlr._fy_TextLayouterResult_formattedText {
		size = size.Max(FontSize(tlr._fy_TextLayouterResult_font, v).Add(Vec2i{0, int32(k) * interLine}))
	}
	tlr.Size = size
}

func (tlr *TextLayouterResult) Draw() TextLayouterRenderable {
	tex := make([]Texture, len(tlr._fy_TextLayouterResult_formattedText))
	for k, v := range tlr._fy_TextLayouterResult_formattedText {
		tex[k] = FontDraw(tlr._fy_TextLayouterResult_font, v)
	}
	rdr := TextLayouterRenderable{
		Size: tlr.Size,
		_fy_TextLayouterResult_lines: tex,
		_fy_TextLayouterResult_interline: FontInterline(tlr._fy_TextLayouterResult_font),
	}
	return rdr
}

func (tlr *TextLayouterRenderable) Draw(r Renderer, pos Vec2i, colour uint32) {
	interLine := tlr._fy_TextLayouterResult_interline
	for _, v := range tlr._fy_TextLayouterResult_lines {
		r.TexRect(v, colour, Area2iOfSize(v.Size()), Area2iFromVecs(pos, v.Size()))
		pos.Y += interLine
	}
}

func fyTextLayouterBreakerNormal(text string, fnt font.Face, xLimit int32, wordwrap bool) TextLayouterResult {
	committedBuffer := TextLayouterResult{
		_fy_TextLayouterResult_font: fnt,
		_fy_TextLayouterResult_formattedText: []string{},
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

/*
 * Seven backends to bind them
 */
func TheOneTextLayouterToRuleThemAll(opts TextLayouterOptions) TextLayouterResult {
	brokenText := fyTextLayouterBreakerNormal(opts.Text, opts.Font, opts.Limits.X, true)
	if brokenText.Size.Y >= opts.Limits.Y {
		brokenText = fyTextLayouterBreakerNormal(opts.Text, opts.Font, opts.Limits.X, false)
	}
	return brokenText
}
