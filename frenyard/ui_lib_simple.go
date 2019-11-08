package frenyard
import "golang.org/x/image/font"

/*
 * UIRect is a 'filler' background element.
 */
type UIRect struct {
	UIElementComponent
	// May be nil.
	tex Texture
	// May be entirely ignored.
	sprite Area2i
	colour uint32
}
func NewColouredRect(colour uint32, size Vec2i) UIRect {
	return UIRect{NewUIElementComponent(size), nil, Area2i{}, colour}
}
func NewTextureRect(colour uint32, tex Texture, sprite Area2i, size Vec2i) UIRect {
	return UIRect{NewUIElementComponent(size), tex, sprite, colour}
}
func (cr *UIRect) FyENormalEvent(ev NormalEvent) {
}
func (cr *UIRect) FyEMouseEvent(ev MouseEvent) {
}
func (cr *UIRect) FyETick(deltaTime float64) {
}
func (cr *UIRect) FyEDraw(target Renderer, under bool) {
	if (!under) {
		if (cr.tex != nil) {
			target.TexRect(cr.tex, cr.colour, cr.sprite, Area2iOfSize(cr.FyESize()))
		} else {
			target.FillRect(cr.colour, Area2iOfSize(cr.FyESize()))
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

/*
 * UILabel displays text.
 */
type UILabel struct {
	UIElementComponent
	UILayoutElementComponent
	LayoutElementNoSubelementsComponent
	_fy_UILabel_Text string
	_fy_UILabel_Font font.Face
	_fy_UILabel_Colour uint32
	_fy_UILabel_Background uint32 // very useful for debugging
	_fy_UILabel_AlignX int8
	_fy_UILabel_AlignY int8
	_fy_UILabel_DidInit bool
	_fy_UILabel_Texture TextLayouterRenderable // Changes based on size!
	_fy_UILabel_TextureLimits fyTextLayoutCacheEntry
	_fy_UILabel_LayoutCache []fyTextLayoutCacheEntry
	_fy_UILabel_PreferredSize Vec2i
}
func NewUILabelPtr(text string, font font.Face, colour uint32, back uint32, alignX int8, alignY int8) *UILabel {
	base := &UILabel{}
	InitUILayoutElement(base)
	base.SetContent(text, font, colour, back, alignX, alignY)
	return base
}

func (cr *UILabel) Text() string {
	return cr._fy_UILabel_Text
}
func (cr *UILabel) SetText(text string) {
	cr.SetContent(text, cr.Font(), cr.Colour(), cr.Background(), cr.AlignX(), cr.AlignY())
}
func (cr *UILabel) Font() font.Face {
	return cr._fy_UILabel_Font
}
func (cr *UILabel) SetFont(font font.Face) {
	cr.SetContent(cr.Text(), font, cr.Colour(), cr.Background(), cr.AlignX(), cr.AlignY())
}
func (cr *UILabel) Colour() uint32 {
	return cr._fy_UILabel_Colour
}
func (cr *UILabel) SetColour(colour uint32) {
	// This doesn't require a reset.
	cr._fy_UILabel_Colour = colour
}
func (cr *UILabel) Background() uint32 {
	return cr._fy_UILabel_Background
}
func (cr *UILabel) SetBackground(colour uint32) {
	// This doesn't require a reset.
	cr._fy_UILabel_Background = colour
}
func (cr *UILabel) AlignX() int8 {
	return cr._fy_UILabel_AlignX
}
func (cr *UILabel) SetAlignX(alignX int8) {
	// This doesn't require a reset.
	cr._fy_UILabel_AlignX = alignX
}
func (cr *UILabel) AlignY() int8 {
	return cr._fy_UILabel_AlignY
}
func (cr *UILabel) SetAlignY(alignY int8) {
	// This doesn't require a reset.
	cr._fy_UILabel_AlignY = alignY
}

func (cr *UILabel) fyLayoutCacheGet(limits Vec2i) fyTextLayoutCacheEntry {
	for _, cacheEntry := range cr._fy_UILabel_LayoutCache {
		if cacheEntry.Matches(limits) {
			return cacheEntry
		}
	}
	entry := fyTextLayoutCacheEntry{}
	entry.Layout = TheOneTextLayouterToRuleThemAll(TextLayouterOptions{
		Text: cr._fy_UILabel_Text,
		Font: cr._fy_UILabel_Font,
		Limits: limits,
	})
	entry.LimitsMin = limits.Min(entry.Layout.Size)
	entry.LimitsMax = limits.Max(entry.Layout.Size)
	cr._fy_UILabel_LayoutCache = append(cr._fy_UILabel_LayoutCache, entry)
	return entry
}

func (cr *UILabel) SetContent(text string, font font.Face, colour uint32, background uint32, alignX int8, alignY int8) {
	cr._fy_UILabel_Text = text
	cr._fy_UILabel_Font = font
	cr._fy_UILabel_Colour = colour
	cr._fy_UILabel_Background = background
	cr._fy_UILabel_AlignX = alignX
	cr._fy_UILabel_AlignY = alignY
	cr._fy_UILabel_LayoutCache = []fyTextLayoutCacheEntry{}
	baseLayout := cr.fyLayoutCacheGet(Vec2iUnlimited())
	cr._fy_UILabel_Texture = baseLayout.Layout.Draw()
	cr._fy_UILabel_TextureLimits = baseLayout
	cr._fy_UILabel_PreferredSize = cr._fy_UILabel_Texture.Size
	if !cr._fy_UILabel_DidInit {
		cr._fy_UILabel_DidInit = true
		cr.FyEResize(cr._fy_UILabel_PreferredSize)
	} else {
		cr.UIThis.ContentChanged()
	}
}

func (cr *UILabel) FyEResize(s Vec2i) {
	cr.UIElementComponent.FyEResize(s)
	// Check to see if we can just use the current texture. This is extremely efficient in comparison.
	if cr._fy_UILabel_TextureLimits.Matches(s) {
		// We can, so we have nothing to do.
		return
	}
	result := cr.fyLayoutCacheGet(s)
	cr._fy_UILabel_Texture = result.Layout.Draw()
	cr._fy_UILabel_TextureLimits = result
}

func (cr *UILabel) FyENormalEvent(ev NormalEvent) {
}
func (cr *UILabel) FyEMouseEvent(ev MouseEvent) {
}
func (cr *UILabel) FyETick(deltaTime float64) {
}
func (cr *UILabel) FyEDraw(target Renderer, under bool) {
	if (!under) {
		labelArea := Area2iOfSize(cr.FyESize())
		if cr._fy_UILabel_Background != 0 {
			target.FillRect(cr._fy_UILabel_Background, labelArea)
		}
		cr._fy_UILabel_Texture.Draw(target, labelArea.Align(cr._fy_UILabel_Texture.Size, cr._fy_UILabel_AlignX, cr._fy_UILabel_AlignY).Pos(), cr._fy_UILabel_Colour)
	}
}

func (cr *UILabel) FyLSizeForLimits(limits Vec2i) Vec2i {
	if limits.Ge(cr._fy_UILabel_PreferredSize) {
		return cr._fy_UILabel_PreferredSize
	}
	return cr.fyLayoutCacheGet(limits).Layout.Size
}
