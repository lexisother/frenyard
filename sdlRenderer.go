package frenyard

import (
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
)

var fySDL2CRTCRegistry *crtcRegistry = newCRTCRegistryPtr()

// sdl2RendererCore is the crtcContext.
type sdl2RendererCore struct {
	base *sdl.Renderer
}

func (r *sdl2RendererCore) osDelete() {
	if r.base == nil {
		panic("Renderer was already destroyed!")
	}
	fySDL2CRTCRegistry.osRemoveRenderer(crtcContext(r))
	r.base.Destroy()
	r.base = nil
}

// sdl2Renderer is the Renderer; meaning is dependent on render target.
type sdl2Renderer struct {
	base      *sdl2RendererCore
	clip      Area2i
	translate Vec2i
}

func newSDL2Renderer(base *sdl2RendererCore) sdl2Renderer {
	return sdl2Renderer{
		base,
		Area2i{},
		Vec2i{},
	}
}

func fySDL2AreaToRect(a Area2i) sdl.Rect {
	return sdl.Rect{
		X: a.X.Pos,
		Y: a.Y.Pos,
		W: a.X.Size,
		H: a.Y.Size,
	}
}
func (r *sdl2Renderer) osFySDL2DrawColour(colour uint32) {
	red := (uint8)((colour >> 16) & 0xFF)
	green := (uint8)((colour >> 8) & 0xFF)
	blue := (uint8)((colour >> 0) & 0xFF)
	alpha := (uint8)((colour >> 24) & 0xFF)
	r.base.base.SetDrawColor(red, green, blue, alpha)
}
func (r *sdl2Renderer) osFyExtractSDL2Texture(tex Texture) *sdl.Texture {
	switch sheetActual := tex.(type) {
	case *crtcTextureExternal:
		// Explicit cast so you can see what's going on with the contexts
		sheetLocal := sheetActual.osGetLocalTexture(crtcContext(r.base)).(*fySDL2LocalTexture)
		return sheetLocal.base
	}
	panic("Unknown texture type forwarded into engine core.")
}
func (r *sdl2Renderer) DrawRect(drc DrawRectCommand) {
	{
		z := sdl2Os()
		defer z.End()
	}
	sRect := fySDL2AreaToRect(drc.TexSprite)
	tRect := fySDL2AreaToRect(drc.Target.Translate(r.translate))

	blendMode := sdl.BlendMode(sdl.BLENDMODE_BLEND)
	if drc.Mode == DrawModeNoBlending {
		blendMode = sdl.BlendMode(sdl.BLENDMODE_NONE)
	} else if drc.Mode == DrawModeAdd {
		blendMode = sdl.BlendMode(sdl.BLENDMODE_ADD)
	} else if drc.Mode == DrawModeModulate {
		blendMode = sdl.BlendMode(sdl.BLENDMODE_MOD)
	}

	if drc.Tex != nil {
		// If the image has zero size, it doesn't exist. Anyway, osGetLocalTexture will crash
		size := drc.Tex.Size()
		if size.X == 0 || size.Y == 0 {
			return
		}
		sheetLocal := r.osFyExtractSDL2Texture(drc.Tex)
		red := (uint8)((drc.Colour >> 16) & 0xFF)
		green := (uint8)((drc.Colour >> 8) & 0xFF)
		blue := (uint8)((drc.Colour >> 0) & 0xFF)
		alpha := (uint8)((drc.Colour >> 24) & 0xFF)
		sheetLocal.SetColorMod(red, green, blue)
		sheetLocal.SetAlphaMod(alpha)
		sheetLocal.SetBlendMode(blendMode)
		r.base.base.Copy(sheetLocal, &sRect, &tRect)
	} else {
		r.osFySDL2DrawColour(drc.Colour)
		r.base.base.SetDrawBlendMode(blendMode)
		r.base.base.FillRect(&tRect)
	}
}
func (r *sdl2Renderer) Clip() Area2i {
	{
		z := sdl2Os()
		defer z.End()
	}
	return r.clip.Translate(r.translate.Negate())
}
func (r *sdl2Renderer) SetClip(val Area2i) {
	{
		z := sdl2Os()
		defer z.End()
	}
	r.clip = val.Translate(r.translate)
	rect := fySDL2AreaToRect(r.clip)
	r.base.base.SetClipRect(&rect)
}
func (r *sdl2Renderer) Translate(val Vec2i) {
	{
		z := sdl2Os()
		defer z.End()
	}
	r.translate = r.translate.Add(val)
}
func (r *sdl2Renderer) Translation() Vec2i {
	{
		z := sdl2Os()
		defer z.End()
	}
	return r.translate
}
func (r *sdl2Renderer) Size() Vec2i {
	{
		z := sdl2Os()
		defer z.End()
	}
	wsX, wsY, err := r.base.base.GetOutputSize()
	if err != nil {
		// SDL2 will actually assertion-fail before ever erroring this way.
		// Why it even has the capability of erroring I do not know.
		panic(err)
	}
	return Vec2i{wsX, wsY}
}
func (r *sdl2Renderer) Reset(colour uint32) {
	{
		z := sdl2Os()
		defer z.End()
	}
	r.translate = Vec2i{}
	r.SetClip(Area2iOfSize(r.Size()))
	r.osFySDL2DrawColour(colour)
	r.base.base.SetDrawBlendMode(sdl.BLENDMODE_NONE)
	r.base.base.Clear()
}

func (r *sdl2Renderer) RenderToTexture(size Vec2i, drawer func(), reserved bool) Texture {
	{
		z := sdl2Os()
		defer z.End()
	}
	if reserved {
		panic("reserved must be kept false for future expansion")
	}

	if size.X <= 0 || size.Y <= 0 {
		// Empty texture
		return GlobalBackend.CreateTexture(Vec2i{}, []uint32{})
	}
	tex, err := r.base.base.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, size.X, size.Y)
	if err != nil {
		// No RTT, so ignore the whole thing and return an empty texture
		return GlobalBackend.CreateTexture(Vec2i{}, []uint32{})
	}
	oldClip := r.Clip()
	oldTranslate := r.translate
	oldRT := r.base.base.GetRenderTarget()

	r.translate = Vec2i{}
	r.SetClip(Area2iOfSize(size))
	r.base.base.SetRenderTarget(tex)

	drawer()

	r.base.base.SetRenderTarget(oldRT)
	r.translate = oldTranslate
	r.SetClip(oldClip)

	localTexture := &fySDL2LocalTexture{
		base: tex,
	}
	return fySDL2CRTCRegistry.osCreateLocalTexture(crtcContext(r.base), localTexture, size)
}
func (r *sdl2Renderer) Present() {
	{
		z := sdl2Os()
		defer z.End()
	}
	r.base.base.Present()
}

type fySDL2LocalTexture struct {
	base *sdl.Texture
}

func (r *fySDL2LocalTexture) osDelete() {
	r.base.Destroy()
}

type fySDL2TextureData struct {
	base *sdl.Surface
}

func (r *fySDL2TextureData) osMakeLocal(render crtcContext) crtcLocalTexture {
	renderActual := render.(*sdl2RendererCore)
	result, err := renderActual.base.CreateTextureFromSurface(r.base)
	if err != nil {
		panic(err)
	}
	return &fySDL2LocalTexture{
		result,
	}
}

func (r *fySDL2TextureData) Size() Vec2i {
	return Vec2i{X: r.base.W, Y: r.base.H}
}
func (r *fySDL2TextureData) osDelete() {
	r.base.Free()
}

// This API is package-local so that sdl_ttf & sdl can get at it
func osSdl2SurfaceToFyTexture(surface *sdl.Surface) Texture {
	return fySDL2CRTCRegistry.osCreateTexture(&fySDL2TextureData{
		base: surface,
	})
}

// { z := sdl2Os(); defer z.End() }
type sdl2BlockOS struct{}

func sdl2Os() sdl2BlockOS { runtime.LockOSThread(); return sdl2BlockOS{} }
func (*sdl2BlockOS) End() { runtime.UnlockOSThread() }
