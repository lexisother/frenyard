package frenyard

// A NinePatch is a resizable rectangular border and background to fit a given container.
type NinePatch struct {
	// If nil, this NinePatch is disabled.
	Tex Texture
	// The area of the image containing the under nine-patch.
	Sprite Area2i
	// The area within that area, absolute, where the container bounds sit.
	Bounds Area2i
	// The area within that area, absolute, where the nine-patch centre bounds sit.
	Centre Area2i
	// Alpha/Colour to modulate by.
	ColourMod uint32
}

// Draw draws the NinePatch on the given renderer with the given container bounds.
func (np NinePatch) Draw(r Renderer, where Area2i) {
	if np.Tex == nil {
		return
	}
	expansionAreas := SplitArea2iGrid3x3(np.Sprite, np.Bounds)
	spriteAreas := SplitArea2iGrid3x3(np.Sprite, np.Centre)
	intrusionAreas := SplitArea2iGrid3x3(np.Bounds, np.Centre)
	expansionMargin := expansionAreas.AsMargin()
	intrusionMargin := intrusionAreas.AsMargin()
	
	whereOuter := where.Expand(expansionMargin)
	whereInner := where.Contract(intrusionMargin)
	drawAreas := SplitArea2iGrid3x3(whereOuter, whereInner)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.A, drawAreas.A)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.B, drawAreas.B)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.C, drawAreas.C)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.D, drawAreas.D)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.E, drawAreas.E)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.F, drawAreas.F)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.G, drawAreas.G)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.H, drawAreas.H)
	r.TexRect(np.Tex, np.ColourMod, spriteAreas.I, drawAreas.I)
}
