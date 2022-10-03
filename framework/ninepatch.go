package framework

import "github.com/uwu/frenyard"

// A NinePatch is a resizable rectangular border and background to fit a given container.
type NinePatch struct {
	// If nil, this NinePatch is disabled.
	Tex frenyard.Texture
	// The area of the image containing the under nine-patch.
	Sprite frenyard.Area2i
	// The area within that area, absolute, where the container bounds sit.
	Bounds frenyard.Area2i
	// The area within that area, absolute, where the nine-patch centre bounds sit.
	Centre frenyard.Area2i
}

// Inset insets the NinePatch by expanding the container bounds.
func (np NinePatch) Inset(margin frenyard.Area2i) NinePatch {
	np.Bounds = np.Bounds.Expand(margin)
	return np
}

// Draw draws the NinePatch on the given renderer with the given container bounds.
func (np NinePatch) Draw(r frenyard.Renderer, where frenyard.Area2i, scale float64, drawBase frenyard.DrawRectCommand) {
	if np.Tex == nil {
		return
	}
	expansionAreas := frenyard.SplitArea2iGrid3x3(np.Sprite, np.Bounds)
	spriteAreas := frenyard.SplitArea2iGrid3x3(np.Sprite, np.Centre)
	intrusionAreas := frenyard.SplitArea2iGrid3x3(np.Bounds, np.Centre)
	expansionMargin := frenyard.ScaleMargin2(scale, expansionAreas.AsMargin())
	intrusionMargin := frenyard.ScaleMargin2(scale, intrusionAreas.AsMargin())

	whereOuter := where.Expand(expansionMargin)
	whereInner := where.Contract(intrusionMargin)
	drawAreas := frenyard.SplitArea2iGrid3x3(whereOuter, whereInner)
	drawBase.Tex = np.Tex

	drawBase.TexSprite = spriteAreas.A
	drawBase.Target = drawAreas.A
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.B
	drawBase.Target = drawAreas.B
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.C
	drawBase.Target = drawAreas.C
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.D
	drawBase.Target = drawAreas.D
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.E
	drawBase.Target = drawAreas.E
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.F
	drawBase.Target = drawAreas.F
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.G
	drawBase.Target = drawAreas.G
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.H
	drawBase.Target = drawAreas.H
	r.DrawRect(drawBase)
	drawBase.TexSprite = spriteAreas.I
	drawBase.Target = drawAreas.I
	r.DrawRect(drawBase)
}
