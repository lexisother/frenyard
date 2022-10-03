package integration

import (
	"github.com/uwu/frenyard"
	"golang.org/x/image/math/fixed"
	"image"
)

// "The Annoyance" is stupid things like characters going behind their own start points.
// "Accounts for the Annoyance" is, essentially, awful workarounds HERE to keep the REST of the system sane.

// TextLayouterOptions contains the options for text layout.
type TextLayouterOptions struct {
	Text TypeChunk
	// SizeUnlimited should be used if an axis should be unbounded.
	Limits frenyard.Vec2i
}

// TextLayouterResult contains the results from text layouting.
type TextLayouterResult struct {
	Area  frenyard.Area2i
	Lines []TypeChunk
}

// Calculates the Area.
func (tlr *TextLayouterResult) fyCalcSize(xLimit int32) {
	bounds := fixed.Rectangle26_6{}
	dot := fixed.Point26_6{}
	for _, v := range tlr.Lines {
		_, lineBounds := v.FyCBounds(dot)
		if bounds.Empty() {
			bounds = lineBounds
		} else {
			bounds = bounds.Union(lineBounds)
		}
		dot = dot.Add(fixed.P(0, v.FyCHeight()))
	}
	tlr.Area = FontRectangleConverter(bounds)
	// Accounts for the Annoyance {
	if tlr.Area.X.Pos > -1 {
		// This tries to keep the X.Pos reasonably constant, which stops the text shifting left/right.
		tlr.Area.X.Size += tlr.Area.X.Pos + 1
		tlr.Area.X.Pos = -1
	}
	tlr.Area.X.Size++
	if tlr.Area.X.Size > xLimit {
		// This is an "at any cost" solution to FORCE things to be met.
		tlr.Area.X.Size = xLimit
	}
	// }
}

// Draw draws the laid-out text to a texture.
func (tlr *TextLayouterResult) Draw() frenyard.Texture {
	img := image.NewNRGBA(image.Rect(0, 0, int(tlr.Area.X.Size), int(tlr.Area.Y.Size)))
	dotFy := tlr.Area.Pos().Negate()
	dot := fixed.P(int(dotFy.X), int(dotFy.Y))
	for _, v := range tlr.Lines {
		v.FyCDraw(img, dot)
		dot = dot.Add(fixed.P(0, v.FyCHeight()))
	}
	return GoImageToTexture(img, []ColourTransform{})
}

func fyTextLayouterBreakerNormal(text TypeChunk, xLimit int32, wordwrap bool) TextLayouterResult {
	// Accounts for the Annoyance {
	xLimit--
	// }
	committedBuffer := TextLayouterResult{
		Lines: []TypeChunk{},
	}
	lineBufferStart := 0
	textCursor := 0
	// End of line buffer is text cursor
	lineBufferDot := fixed.Point26_6{}
	lastSpaceComponent := -1
	// Do be alerted! This retrieves runes, but slices use byte indexes.
	carriageReturn := false
	textLength := text.FyCComponentCount()
	for textCursor < textLength {
		breakStatus := text.FyCComponentBreakStatus(textCursor)
		if breakStatus == TypeChunkComponentBreakStatusNewline {
			// Newline, commit buffer
			committedBuffer.Lines = append(committedBuffer.Lines, text.FyCSection(lineBufferStart, textCursor))
			textCursor++
			lineBufferStart = textCursor
			lastSpaceComponent = -1
			carriageReturn = true
			continue
		}
		if carriageReturn {
			lineBufferDot = fixed.Point26_6{}
			advance, _ := text.FyCSection(lineBufferStart, textCursor).FyCBounds(lineBufferDot)
			lineBufferDot = lineBufferDot.Add(advance)
		}
		// Insert character
		advance := text.FyCComponentAdvance(textCursor, textCursor != lineBufferStart)
		lineBufferDot = lineBufferDot.Add(fixed.Point26_6{X: advance})
		if int32(lineBufferDot.X.Ceil()) >= xLimit {
			// Check for cut position. Disabling wordwrap prevents these from ever being found.
			if lastSpaceComponent >= 0 {
				// Append the confirmed text.
				committedBuffer.Lines = append(committedBuffer.Lines, text.FyCSection(lineBufferStart, lastSpaceComponent))
				// Reset the line buffer to being at the " "
				lineBufferStart = lastSpaceComponent
				textCursor = lastSpaceComponent + 1
				carriageReturn = true
				lastSpaceComponent = -1
			} else {
				// Character break. Much simpler, doesn't even require interrupting the stream.
				// If the character actually failed here we'd be doomed anyway, so place it now, too.
				if lineBufferStart == textCursor-1 && text.FyCComponentBreakStatus(lineBufferStart) == TypeChunkComponentBreakStatusSpace {
					// The line buffer just contains a space.
					// Draw this char and then CR immediately since it's clear there's pretty much no room.
					committedBuffer.Lines = append(committedBuffer.Lines, text.FyCSection(textCursor, textCursor+1))
					textCursor++
					lineBufferStart = textCursor
					carriageReturn = true
					// Line buffer is now empty.
				} else {
					committedBuffer.Lines = append(committedBuffer.Lines, text.FyCSection(lineBufferStart, textCursor))
					// This includes the character in the new line buffer.
					lineBufferStart = textCursor
					textCursor++
					carriageReturn = true
				}
			}
		} else {
			if breakStatus == TypeChunkComponentBreakStatusSpace && wordwrap {
				// Successfully inserting space, let's note that
				lastSpaceComponent = textCursor
			}
			textCursor++
		}
	}
	if lineBufferStart != textCursor {
		committedBuffer.Lines = append(committedBuffer.Lines, text.FyCSection(lineBufferStart, textCursor))
	}
	committedBuffer.fyCalcSize(xLimit)
	return committedBuffer
}

// TheOneTextLayouterToRuleThemAll lays out text with wrapping limits and other such constraints.
func TheOneTextLayouterToRuleThemAll(opts TextLayouterOptions) TextLayouterResult {
	brokenText := fyTextLayouterBreakerNormal(opts.Text, opts.Limits.X, true)
	if brokenText.Area.Y.Size >= opts.Limits.Y {
		brokenText = fyTextLayouterBreakerNormal(opts.Text, opts.Limits.X, false)
	}
	return brokenText
}
