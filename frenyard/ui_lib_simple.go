package frenyard
import "golang.org/x/image/font"

// UIRect is a 'filler' background element.
type UIRect struct {
	UIElementComponent
	// May be nil.
	Texture Texture
	// If Texture is nil, this is ignored.
	Sprite Area2i
	// Colour (either modulates Texture if present, or is filled as-is)
	Colour uint32
}
// NewColouredRectPtr creates a UIRect given a colour and size.
func NewColouredRectPtr(colour uint32, size Vec2i) *UIRect {
	return &UIRect{NewUIElementComponent(size), nil, Area2i{}, colour}
}
// NewTextureRectPtr creates a UIRect given a colour, texture, sprite area, and size.
func NewTextureRectPtr(colour uint32, tex Texture, sprite Area2i, size Vec2i) *UIRect {
	return &UIRect{NewUIElementComponent(size), tex, sprite, colour}
}
// FyENormalEvent implements UIElement.FyENormalEvent
func (cr *UIRect) FyENormalEvent(ev NormalEvent) {
}
// FyEMouseEvent implements UIElement.FyEMouseEvent
func (cr *UIRect) FyEMouseEvent(ev MouseEvent) {
}
// FyETick implements UIElement.FyETick
func (cr *UIRect) FyETick(deltaTime float64) {
}
// FyEDraw implements UIElement.FyEDraw
func (cr *UIRect) FyEDraw(target Renderer, under bool) {
	if (!under) {
		if (cr.Texture != nil) {
			target.TexRect(cr.Texture, cr.Colour, cr.Sprite, Area2iOfSize(cr.FyESize()))
		} else {
			target.FillRect(cr.Colour, Area2iOfSize(cr.FyESize()))
		}
	}
}

type fyTextLayoutCacheEntry struct {
	Layout TextLayouterResult
	LimitsMin Vec2i
	LimitsMax Vec2i
}
func (cacheEntry fyTextLayoutCacheEntry) Matches(limits Vec2i) bool {
	return limits.Ge(cacheEntry.LimitsMin) && cacheEntry.LimitsMax.Ge(limits)
}

// UILabel displays text.
type UILabel struct {
	UIElementComponent
	UILayoutElementComponent
	LayoutElementNoSubelementsComponent
	_text string
	_font font.Face
	_colour uint32
	_background uint32 // very useful for debugging
	_alignment Alignment2i
	_didInit bool
	_texture TextLayouterRenderable // Changes based on size!
	_textureLimits fyTextLayoutCacheEntry
	_layoutCache []fyTextLayoutCacheEntry
	_preferredSize Vec2i
}
// NewUILabelPtr creates a new UILabel from the various visual details about it.
func NewUILabelPtr(text string, font font.Face, colour uint32, back uint32, align Alignment2i) *UILabel {
	base := &UILabel{}
	InitUILayoutElementComponent(base)
	base.SetTextAndFont(text, font)
	base._colour = colour
	base._background = back
	base._alignment = align
	return base
}

// Text gets the text.
func (cr *UILabel) Text() string {
	return cr._text
}

// SetText sets the text.
func (cr *UILabel) SetText(text string) {
	cr.SetTextAndFont(text, cr.Font())
}

// Font gets the font.
func (cr *UILabel) Font() font.Face {
	return cr._font
}

// SetFont sets the font.
func (cr *UILabel) SetFont(font font.Face) {
	cr.SetTextAndFont(cr.Text(), font)
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
func (cr *UILabel) Alignment() Alignment2i {
	return cr._alignment
}

// SetAlignment sets the alignment.
func (cr *UILabel) SetAlignment(align Alignment2i) {
	// This doesn't require a reset.
	cr._alignment = align
}

func (cr *UILabel) fyLayoutCacheGet(limits Vec2i) fyTextLayoutCacheEntry {
	for _, cacheEntry := range cr._layoutCache {
		if cacheEntry.Matches(limits) {
			return cacheEntry
		}
	}
	entry := fyTextLayoutCacheEntry{}
	entry.Layout = TheOneTextLayouterToRuleThemAll(TextLayouterOptions{
		Text: cr._text,
		Font: cr._font,
		Limits: limits,
	})
	entry.LimitsMin = limits.Min(entry.Layout.Size)
	entry.LimitsMax = limits.Max(entry.Layout.Size)
	cr._layoutCache = append(cr._layoutCache, entry)
	return entry
}

// SetTextAndFont sets the layout-affecting content of the label atomically, which helps with performance over individual changes.
func (cr *UILabel) SetTextAndFont(text string, font font.Face) {
	cr._text = text
	cr._font = font
	cr._layoutCache = []fyTextLayoutCacheEntry{}
	baseLayout := cr.fyLayoutCacheGet(Vec2iUnlimited())
	cr._texture = baseLayout.Layout.Draw()
	cr._textureLimits = baseLayout
	cr._preferredSize = cr._texture.Size
	if !cr._didInit {
		cr._didInit = true
		cr.FyEResize(cr._preferredSize)
	} else {
		cr.ThisUILayoutElementComponentDetails.ContentChanged()
	}
}

// FyEResize overrides UIElementComponent.FyEResize
func (cr *UILabel) FyEResize(s Vec2i) {
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
func (cr *UILabel) FyENormalEvent(ev NormalEvent) {
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (cr *UILabel) FyEMouseEvent(ev MouseEvent) {
}

// FyETick implements UIElement.FyETick
func (cr *UILabel) FyETick(deltaTime float64) {
}

// FyEDraw implements UIElement.FyEDraw
func (cr *UILabel) FyEDraw(target Renderer, under bool) {
	if (!under) {
		labelArea := Area2iOfSize(cr.FyESize())
		if cr._background != 0 {
			target.FillRect(cr._background, labelArea)
		}
		cr._texture.Draw(target, labelArea.Align(cr._texture.Size, cr._alignment).Pos(), cr._colour)
	}
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (cr *UILabel) FyLSizeForLimits(limits Vec2i) Vec2i {
	if limits.Ge(cr._preferredSize) {
		return cr._preferredSize
	}
	return cr.fyLayoutCacheGet(limits).Layout.Size
}
