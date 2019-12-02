package frenyard

// DrawMode represents some form of drawing mode for the primitive. (Documentation for these should include a GLSLish equal for those versed in graphics, ideally translated from the SDL2 documentation.)
type DrawMode uint8

// Do be aware that DrawModeNormal is the zero value. This is "off" from SDL2 order but is done on purpose as this is the default that most people will want.

// DrawModeNormal is a normal blended draw mode; typically "what you want". `vec4(mix(dst.rgb, src.rgb, src.a), src.a + (dstA * (1 - src.a)))`
const DrawModeNormal = 0
// DrawModeNoBlending disables blending. `src`
const DrawModeNoBlending = 1
// DrawModeAdd is an "additive" blend mode. `vec4(dst.rgb + src.rgb, dst.a)`
const DrawModeAdd = 2
// DrawModeModulate is a "modulate" blend mode. `vec4(dst.rgb * src.rgb, dst.a)`
const DrawModeModulate = 3

// DrawRectCommand represents a rectangle-based drawing command.
type DrawRectCommand struct {
	// Where to draw
	Target Area2i
	// Texture
	Tex Texture
	// Area within texture to get pixels from.
	TexSprite Area2i
	// Colour [modulation if Texture is given]
	Colour uint32
	// Blending Mode/etc.
	Mode DrawMode
}

// Renderer is an abstract rendering interface.
type Renderer interface {
	/* Draw */

	// Draws a DrawRectCommand.
	DrawRect(DrawRectCommand)

	/* Control */

	// Translates the renderer's target. Use with defer to undo later.
	Translate(vec Vec2i)
	// Sets the clip area, relative to the current translation.
	SetClip(clip Area2i)
	// Gets the clip area, relative to the current translation.
	Clip() Area2i

	/* Outer Control */

	// Gets the size of the drawing area.
	Size() Vec2i
	// Clears & resets clip/translate/etc. You should do this at the start of  frame.
	Reset(colour uint32)
}

// Texture interface. This is automatically deleted on finalization.
type Texture interface {
	Size() Vec2i
}
