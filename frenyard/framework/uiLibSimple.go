package framework

import (
	"golang.org/x/image/font"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

// UIRect is a 'filler' background element.
type UIRect struct {
	UIElementComponent
	// May be nil.
	Tex frenyard.Texture
	// If Texture is nil, this is ignored.
	Sprite frenyard.Area2i
	// Colour (either modulates Texture if present, or is filled as-is)
	Colour uint32
	// Alignment
	Alignment frenyard.Alignment2i
}

// NewColouredRectPtr creates a UIRect given a colour and size.
func NewColouredRectPtr(colour uint32, size frenyard.Vec2i) *UIRect {
	return &UIRect{NewUIElementComponent(size), nil, frenyard.Area2i{}, colour, frenyard.Alignment2i{}}
}

// NewTextureRectPtr creates a UIRect given a colour, texture, sprite area, and size.
func NewTextureRectPtr(colour uint32, tex frenyard.Texture, sprite frenyard.Area2i, size frenyard.Vec2i, alignment frenyard.Alignment2i) *UIRect {
	return &UIRect{NewUIElementComponent(size), tex, sprite, colour, alignment}
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (cr *UIRect) FyENormalEvent(ev frenyard.NormalEvent) {
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (cr *UIRect) FyEMouseEvent(ev frenyard.MouseEvent) {
}

// FyETick implements UIElement.FyETick
func (cr *UIRect) FyETick(deltaTime float64) {
}

// FyEDraw implements UIElement.FyEDraw
func (cr *UIRect) FyEDraw(target frenyard.Renderer, under bool) {
	if !under {
		target.DrawRect(frenyard.DrawRectCommand{
			Tex: cr.Tex,
			Colour: cr.Colour,
			TexSprite: cr.Sprite,
			Target: frenyard.Area2iOfSize(cr.FyESize()).Align(cr.Sprite.Size(), frenyard.Alignment2i{}),
		})
	}
}

type fyTextLayoutCacheEntry struct {
	Layout    integration.TextLayouterResult
	LimitsMin frenyard.Vec2i
	LimitsMax frenyard.Vec2i
}

func (cacheEntry fyTextLayoutCacheEntry) Matches(limits frenyard.Vec2i) bool {
	return limits.Ge(cacheEntry.LimitsMin) && cacheEntry.LimitsMax.Ge(limits)
}

// UILabel displays text.
type UILabel struct {
	UIElementComponent
	UILayoutElementComponent
	LayoutElementNoSubelementsComponent
	_text          integration.TypeChunk
	_font          font.Face
	_colour        uint32
	_background    uint32 // very useful for debugging
	_alignment     frenyard.Alignment2i
	_didInit       bool
	_texture       frenyard.Texture // Changes based on size!
	_textureLimits fyTextLayoutCacheEntry
	_layoutCache   []fyTextLayoutCacheEntry
	_preferredSize frenyard.Vec2i
}

// NewUILabelPtr creates a new UILabel from the various visual details about it.
func NewUILabelPtr(text integration.TypeChunk, colour uint32, back uint32, align frenyard.Alignment2i) *UILabel {
	base := &UILabel{}
	InitUILayoutElementComponent(base)
	base.SetText(text)
	base._colour = colour
	base._background = back
	base._alignment = align
	return base
}

// Text gets the document.
func (cr *UILabel) Text() integration.TypeChunk {
	return cr._text
}

// SetText sets the document.
func (cr *UILabel) SetText(text integration.TypeChunk) {
	cr._text = text
	cr._layoutCache = []fyTextLayoutCacheEntry{}
	baseLayout := cr.fyLayoutCacheGet(frenyard.Vec2iUnlimited())
	cr._texture = baseLayout.Layout.Draw()
	cr._textureLimits = baseLayout
	cr._preferredSize = cr._texture.Size()
	if !cr._didInit {
		cr._didInit = true
		cr.FyEResize(cr._preferredSize)
	} else {
		cr.ThisUILayoutElementComponentDetails.ContentChanged()
	}
}

// Colour gets the colour.
func (cr *UILabel) Colour() uint32 {
	return cr._colour
}

// SetColour sets the colour.
func (cr *UILabel) SetColour(colour uint32) {
	// This doesn't require a reset.
	cr._colour = colour
}

// Background gets the background colour
func (cr *UILabel) Background() uint32 {
	return cr._background
}

// SetBackground sets the background colour.
func (cr *UILabel) SetBackground(colour uint32) {
	// This doesn't require a reset.
	cr._background = colour
}

// Alignment gets the alignment.
func (cr *UILabel) Alignment() frenyard.Alignment2i {
	return cr._alignment
}

// SetAlignment sets the alignment.
func (cr *UILabel) SetAlignment(align frenyard.Alignment2i) {
	// This doesn't require a reset.
	cr._alignment = align
}

func (cr *UILabel) fyLayoutCacheGet(limits frenyard.Vec2i) fyTextLayoutCacheEntry {
	for _, cacheEntry := range cr._layoutCache {
		if cacheEntry.Matches(limits) {
			return cacheEntry
		}
	}
	entry := fyTextLayoutCacheEntry{}
	entry.Layout = integration.TheOneTextLayouterToRuleThemAll(integration.TextLayouterOptions{
		Text:   cr._text,
		Limits: limits,
	})
	areaSize := entry.Layout.Area.Size()
	entry.LimitsMin = limits.Min(areaSize)
	entry.LimitsMax = limits.Max(areaSize)
	cr._layoutCache = append(cr._layoutCache, entry)
	return entry
}

// FyEResize overrides UIElementComponent.FyEResize
func (cr *UILabel) FyEResize(s frenyard.Vec2i) {
	cr.UIElementComponent.FyEResize(s)
	// Check to see if we can just use the current texture. This is extremely efficient in comparison.
	if cr._textureLimits.Matches(s) {
		// We can, so we have nothing to do.
		return
	}
	result := cr.fyLayoutCacheGet(s)
	cr._texture = result.Layout.Draw()
	cr._textureLimits = result
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (cr *UILabel) FyENormalEvent(ev frenyard.NormalEvent) {
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (cr *UILabel) FyEMouseEvent(ev frenyard.MouseEvent) {
}

// FyETick implements UIElement.FyETick
func (cr *UILabel) FyETick(deltaTime float64) {
}

// FyEDraw implements UIElement.FyEDraw
func (cr *UILabel) FyEDraw(target frenyard.Renderer, under bool) {
	if !under {
		labelArea := frenyard.Area2iOfSize(cr.FyESize())
		if cr._background != 0 {
			target.DrawRect(frenyard.DrawRectCommand{
				Colour: cr._background,
				Target: labelArea,
			})
		}
		texSize := cr._texture.Size()
		target.DrawRect(frenyard.DrawRectCommand{
			Tex: cr._texture,
			TexSprite: frenyard.Area2iOfSize(texSize),
			Colour: cr._colour,
			Target: labelArea.Align(texSize, cr._alignment),
		})
	}
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (cr *UILabel) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	if limits.Ge(cr._preferredSize) {
		return cr._preferredSize
	}
	return cr.fyLayoutCacheGet(limits).Layout.Area.Size()
}
