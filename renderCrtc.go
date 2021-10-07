package frenyard

import (
	"runtime"
)

type crtcContext interface{}

// Cross-Renderer Texture Cache.
// Not needed until we have to deal with multiple windows or render-targets.

type crtcRegistry struct {
	// A queue of textures to delete. Must be buffered or else.
	InternalTexturesToDelete chan *crtcTextureInternal
	// The full set of crtcLocalTexture entries. Has to be here for Renderer deletion.
	_fyCrtcRegistryAllLocalTextures map[crtcContext]map[*crtcTextureInternal]crtcLocalTexture
}

func newCRTCRegistryPtr() *crtcRegistry {
	return &crtcRegistry{
		make(chan *crtcTextureInternal, 1),
		map[crtcContext]map[*crtcTextureInternal]crtcLocalTexture{},
	}
}

// CreateTexture is at the very bottom as the creator of crtcTextureInternal

// Cleans up internal textures that are no longer referenced.
// Run this every so often.
func (cr *crtcRegistry) osFlush() {
	if len(cr.InternalTexturesToDelete) > 0 {
		tex := <-cr.InternalTexturesToDelete
		tex.osDelete()
	}
}

func (cr *crtcRegistry) osGetRendererEntry(r crtcContext) map[*crtcTextureInternal]crtcLocalTexture {
	alt, altPresent := cr._fyCrtcRegistryAllLocalTextures[r]
	if !altPresent {
		alt = map[*crtcTextureInternal]crtcLocalTexture{}
		cr._fyCrtcRegistryAllLocalTextures[r] = alt
	}
	return alt
}

// Use to notify CRTC that a Renderer is going away.
// This deletes all local textures for the given renderer.
func (cr *crtcRegistry) osRemoveRenderer(r crtcContext) {
	remnant, present := cr._fyCrtcRegistryAllLocalTextures[r]
	if present {
		for _, v := range remnant {
			v.osDelete()
		}
		delete(cr._fyCrtcRegistryAllLocalTextures, r)
	}
}

type crtcLocalTexture interface {
	osDelete()
}

/*
 * This is an abstraction over the backend's cross-renderer image format (such as sdl.Surface)
 * Inherits Texture because this is the real Texture implementation
 */
type crtcTextureData interface {
	Texture
	osMakeLocal(target crtcContext) crtcLocalTexture
	// This does not have to delete the local textures; that is automatic.
	osDelete()
}

// This is the 'cross-renderer texture' structure. This is where the cache is attached.
type crtcTextureInternal struct {
	// The registry hosting this. (Creates a reference loop, but it's all cleaned up later)
	Registry *crtcRegistry
	// The backend's data. May be nil if this is a local-only texture.
	Data crtcTextureData
	// The size of the texture.
	Size Vec2i
}

func (c *crtcTextureInternal) osDelete() {
	if c.Data != nil {
		c.Data.osDelete()
	}
	for _, rv := range c.Registry._fyCrtcRegistryAllLocalTextures {
		purgeMe := rv[c]
		if purgeMe != nil {
			purgeMe.osDelete()
		}
		delete(rv, c)
	}
}

// DO NOT HOLD IN crtcLocalTexture or fyCRTCTextureReal! This value is the outer wrapper and finalization is used as a sentinel.
type crtcTextureExternal struct {
	Internal *crtcTextureInternal
}

func (cte *crtcTextureExternal) Size() Vec2i {
	return cte.Internal.Size
}

func (cte *crtcTextureExternal) osGetLocalTexture(r crtcContext) crtcLocalTexture {
	alt := cte.Internal.Registry.osGetRendererEntry(r)

	localTexture := alt[cte.Internal]
	if localTexture == nil {
		if cte.Internal.Data == nil {
			panic("Attempted to use local texture outside of valid renderer")
		}
		localTexture = cte.Internal.Data.osMakeLocal(r)
		alt[cte.Internal] = localTexture
	}
	return localTexture
}

func fyCRTCTextureFinalizer2(t *crtcTextureInternal) {
	t.Registry.InternalTexturesToDelete <- t
}
func fyCRTCTextureFinalizer(ext *crtcTextureExternal) {
	go fyCRTCTextureFinalizer2(ext.Internal)
}

func (cr *crtcRegistry) osCreateTexture(data crtcTextureData) Texture {
	internal := &crtcTextureInternal{
		Registry: cr,
		Data:     data,
		Size:     data.Size(),
	}
	external := &crtcTextureExternal{internal}
	runtime.SetFinalizer(external, fyCRTCTextureFinalizer)
	return external
}

func (cr *crtcRegistry) osCreateLocalTexture(renderer crtcContext, data crtcLocalTexture, size Vec2i) Texture {
	internal := &crtcTextureInternal{cr, nil, size}
	external := &crtcTextureExternal{internal}
	cr.osGetRendererEntry(renderer)[internal] = data
	runtime.SetFinalizer(external, fyCRTCTextureFinalizer)
	return external
}
