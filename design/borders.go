package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

/*
 * Before continuing, further details on how this module specifically works must be provided.
 * 
 * Firstly, it is to be understood that SDL2 Renderer DOES NOT support mipmapping.
 * It may supposedly support linear upscaling but in practice I have not found this to be working.
 * 
 * Furthermore, it'd be nice to be able to upscale the assets in future.
 * As such.
 * 
 * Ninepatches for this module are typically drawn by usage of a browser showing a special HTML page.
 * This is because a lot of elements from MD, which we're mimicking for continuity with the website, are designed assuming CSS box-shadow support.
 * Which we do not have.
 * 
 * The informal standard I've made up for this uses a "reference size" of a 32x32 image, with the component within a 16x16 centred square.
 * The units here are expected to match browser "px" units in a 1:1 px-to-real-pixel environment.
 * These are also known as "dp" in Material Design specifications.
 * Images are then postfixed with a multiplier; one dp is that many pixels in the image.
 * "X4" is presently the 'standard' here.
 * The images are then automatically *downscaled* as necessary.
 */

var borderImageScale int32
var borderEffectiveScale float64
var borderGenTextureMask frenyard.Texture
var borderGenTextureShadow frenyard.Texture
var borderGenTextureRaw frenyard.Texture

var borderButton framework.NinePatch
var borderButtonRaw framework.NinePatch
var borderButtonShadow framework.NinePatch
var borderButtonShadowFocus framework.NinePatch

// Used for searchboxes
var searchboxTheme framework.NinePatchFrame

// ScrollboxExterior should be wrapped around scrollboxen.
var ScrollboxExterior framework.NinePatchFrame

// ScrollbarThemeH is the ScrollbarTheme for horizontal scrollbars.
var ScrollbarThemeH framework.ScrollbarTheme
// ScrollbarThemeV is the ScrollbarTheme for vertical scrollbars.
var ScrollbarThemeV framework.ScrollbarTheme

func deSetupBorders() {
	// This must all be kept in sync!
	
	borderImageScale = 4
	generationImage := integration.CreateHardcodedPNGImage(generationX4B64)
	borderEffectiveScale = DesignScale / float64(borderImageScale)
	
	for borderImageScale > 1 && borderEffectiveScale <= 0.5 {
		borderEffectiveScale *= 2
		borderImageScale /= 2
		
		generationImage = integration.ScaleImageToHalfSize(generationImage)
	}

	borderGenTextureMask = integration.GoImageToTexture(generationImage, []integration.ColourTransform{integration.ColourTransformBlueToStencil})
	borderGenTextureShadow = integration.GoImageToTexture(generationImage, []integration.ColourTransform{integration.ColourTransformInvert, integration.ColourTransformBlueToStencil, integration.ColourTransformInvert})
	borderGenTextureRaw = integration.GoImageToTexture(generationImage, []integration.ColourTransform{})
	
	deBorderGenInit()
	
	// -- Standard ninepatches --
	borderButton = borderGenRounded4dpMaskX4Mask
	borderButtonRaw = borderGenRounded4dpMaskX4Raw
	borderButtonShadow = borderGenRounded4dpShadow2dpX4Shadow
	borderButtonShadowFocus = borderGenRounded4dpShadow4dpX4Shadow
	
	// -- Scrollbar Theme --
	sbBaseInframe := sizeScale(2)
	sbBase := sizeScale(4)
	sbMovement := sizeScale(4)
	sbMovementLong := sizeScale(8)
	ScrollbarThemeH.Base = framework.NinePatchFrame{
		Layers: []framework.NinePatchFrameLayer{
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenSquareHeadshadow2dpX4Shadow,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF101010,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenSquareMaskX4Mask,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF181818,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenSquareMaskX4Mask.Inset(frenyard.Area2iMargin(sbBaseInframe, sbBaseInframe, sbBaseInframe, sbBaseInframe)),
				Scale: borderEffectiveScale,
				ColourMod: 0xFF101010,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenSquareHeadshadowInset1dpX4Shadow.Inset(frenyard.Area2iMargin(sbBaseInframe, sbBaseInframe, sbBaseInframe, sbBaseInframe)),
				Scale: borderEffectiveScale,
				ColourMod: 0xFFFFFFFF,
			},
		},
		Padding: frenyard.Area2iMargin(sbBase, sbBase, sbBase, sbBase),
	}
	movementFrameH := framework.NinePatchFrame{
		Layers: []framework.NinePatchFrameLayer{
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenRounded4dpHeadshadow1dpX4Shadow,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF606060,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenRounded4dpMaskX4Mask,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF606060,
			},
		},
		Padding:  frenyard.Area2iMargin(sbMovementLong, sbMovement, sbMovementLong, sbMovement),
	}
	ScrollbarThemeH.Movement = movementFrameH
	ScrollbarThemeV = ScrollbarThemeH
	movementFrameV := movementFrameH
	movementFrameV.Padding = frenyard.Area2i{X: movementFrameH.Padding.Y, Y: movementFrameH.Padding.X}
	ScrollbarThemeV.Movement = movementFrameV
	
	ScrollboxExterior = framework.NinePatchFrame{
		Layers: []framework.NinePatchFrameLayer{
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassUnderBefore,
				NinePatch: borderGenSquareMaskX4Mask,
				Scale: borderEffectiveScale,
				ColourMod: ThemeBackgroundUnderlayer,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverAfter,
				NinePatch: borderGenSquareHeadshadowInset2dpX4Shadow,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF000000,
			},
		},
	}
	
	// -- Searchboxes --
	searchboxMargin := sizeScale(4)
	searchboxTheme = framework.NinePatchFrame{
		Layers: []framework.NinePatchFrameLayer{
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassUnderBefore,
				NinePatch: borderGenRounded4dpShadow2dpX4Shadow,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF000000,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenRounded4dpMaskX4Mask,
				Scale: borderEffectiveScale,
				ColourMod: ThemeBackgroundSearch,
			},
		},
		Padding: frenyard.Area2iMargin(searchboxMargin, searchboxMargin, searchboxMargin, searchboxMargin),
	}
}

// Returns a "standard ninepatch" for the autogenerated stuff.
// 'n' represents the bounds <-> centre margin in prescaled units.
func _borderStandardBorderNinepatch(base framework.NinePatch, n int32) framework.NinePatch {
	base.Sprite = frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 32 * borderImageScale, Y: 32 * borderImageScale})
	base.Bounds = frenyard.Area2iFromVecs(frenyard.Vec2i{X: 8 * borderImageScale, Y: 8 * borderImageScale}, frenyard.Vec2i{X: 16 * borderImageScale, Y: 16 * borderImageScale})
	base.Centre = frenyard.Area2iFromVecs(frenyard.Vec2i{X: (8 + n) * borderImageScale, Y: (8 + n) * borderImageScale}, frenyard.Vec2i{X: (16 - (n * 2)) * borderImageScale, Y: (16 - (n * 2)) * borderImageScale})
	return base
}

// This border deliberately "hangs over" on the OVER layer.
// With correct manipulation of Z (ensure title is LAST in flexbox & use Order to correct)
// this produces a shadow effect that even hangs over other UI.

// BorderTitle produces a border for the shadowing effect under a title.
func BorderTitle(colour uint32) framework.Frame {
	addedBorderX := sizeScale(8)
	addedBorderY := sizeScale(8)
	return framework.NinePatchFrame{
		Layers: []framework.NinePatchFrameLayer{
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenSquareHeadshadow2dpX4Shadow,
				Scale: borderEffectiveScale,
				ColourMod: colour,
			},
			framework.NinePatchFrameLayer{
				Pass: framework.FramePassOverBefore,
				NinePatch: borderGenSquareMaskX4Mask,
				Scale: borderEffectiveScale,
				ColourMod: colour,
			},
		},
		Padding:  frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2}),
		Clipping: true,
	}
}
