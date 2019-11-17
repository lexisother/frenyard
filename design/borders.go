package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
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

var borderButton frenyard.NinePatch
var borderButtonShadow frenyard.NinePatch
var borderButtonShadowFocus frenyard.NinePatch

// ScrollboxExterior should be wrapped around scrollboxen.
var ScrollboxExterior frenyard.NinePatchFrame

// ScrollbarThemeH is the ScrollbarTheme for horizontal scrollbars.
var ScrollbarThemeH frenyard.ScrollbarTheme
// ScrollbarThemeV is the ScrollbarTheme for vertical scrollbars.
var ScrollbarThemeV frenyard.ScrollbarTheme

func deSetupBorders() {
	// This must all be kept in sync!
	
	borderImageScale = 4
	generationImage := frenyard.CreateHardcodedPNGImage(generationX4B64)
	borderEffectiveScale = DesignScale / float64(borderImageScale)
	
	for borderImageScale > 1 && borderEffectiveScale <= 0.5 {
		borderEffectiveScale *= 2
		borderImageScale /= 2
		
		generationImage = frenyard.ScaleImageToHalfSize(generationImage)
	}

	borderGenTextureMask = frenyard.GoImageToTexture(generationImage, []frenyard.ColourTransform{frenyard.ColourTransformBlueToStencil})
	borderGenTextureShadow = frenyard.GoImageToTexture(generationImage, []frenyard.ColourTransform{frenyard.ColourTransformInvert, frenyard.ColourTransformBlueToStencil, frenyard.ColourTransformInvert})
	
	deBorderGenInit()
	
	// -- Standard ninepatches --
	borderButton = borderGenRounded4dpMaskX4
	borderButtonShadow = borderGenRounded4dpShadow2dpX4
	borderButtonShadowFocus = borderGenRounded4dpShadow4dpX4
	
	// -- Scrollbar Theme --
	sbBaseInframe := sizeScale(2)
	sbBase := sizeScale(4)
	sbMovement := sizeScale(4)
	sbMovementLong := sizeScale(8)
	ScrollbarThemeH.Base = frenyard.NinePatchFrame{
		Layers: []frenyard.NinePatchFrameLayer{
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenSquareHeadshadow2dpX4,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF101010,
			},
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenSquareMaskX4,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF181818,
			},
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenSquareMaskX4.Inset(frenyard.Area2iMargin(sbBaseInframe, sbBaseInframe, sbBaseInframe, sbBaseInframe)),
				Scale: borderEffectiveScale,
				ColourMod: 0xFF101010,
			},
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenSquareHeadshadowInset1dpX4.Inset(frenyard.Area2iMargin(sbBaseInframe, sbBaseInframe, sbBaseInframe, sbBaseInframe)),
				Scale: borderEffectiveScale,
				ColourMod: 0xFFFFFFFF,
			},
		},
		Padding: frenyard.Area2iMargin(sbBase, sbBase, sbBase, sbBase),
	}
	movementFrameH := frenyard.NinePatchFrame{
		Layers: []frenyard.NinePatchFrameLayer{
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenRounded4dpHeadshadow1dpX4,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF606060,
			},
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenRounded4dpMaskX4,
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
	
	ScrollboxExterior = frenyard.NinePatchFrame{
		Layers: []frenyard.NinePatchFrameLayer{
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassUnderBefore,
				NinePatch: borderGenSquareMaskX4,
				Scale: borderEffectiveScale,
				ColourMod: ThemeBackgroundUnderlayer,
			},
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverAfter,
				NinePatch: borderGenSquareHeadshadowInset2dpX4,
				Scale: borderEffectiveScale,
				ColourMod: 0xFF000000,
			},
		},
	}
}

// Returns a "standard ninepatch" for the autogenerated stuff.
// 'n' represents the bounds <-> centre margin in prescaled units.
func _borderStandardBorderNinepatch(base frenyard.NinePatch, n int32) frenyard.NinePatch {
	base.Sprite = frenyard.Area2iFromVecs(frenyard.Vec2i{X: 0, Y: 0}, frenyard.Vec2i{X: 32 * borderImageScale, Y: 32 * borderImageScale})
	base.Bounds = frenyard.Area2iFromVecs(frenyard.Vec2i{X: 8 * borderImageScale, Y: 8 * borderImageScale}, frenyard.Vec2i{X: 16 * borderImageScale, Y: 16 * borderImageScale})
	base.Centre = frenyard.Area2iFromVecs(frenyard.Vec2i{X: (8 + n) * borderImageScale, Y: (8 + n) * borderImageScale}, frenyard.Vec2i{X: (16 - (n * 2)) * borderImageScale, Y: (16 - (n * 2)) * borderImageScale})
	return base
}

// This border deliberately "hangs over" on the OVER layer.
// With correct manipulation of Z (ensure title is LAST in flexbox & use Order to correct)
// this produces a shadow effect that even hangs over other UI.

// BorderTitle produces a border for the shadowing effect under a title.
func BorderTitle(colour uint32) frenyard.Frame {
	addedBorderX := sizeScale(8)
	addedBorderY := sizeScale(8)
	return frenyard.NinePatchFrame{
		Layers: []frenyard.NinePatchFrameLayer{
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenSquareHeadshadow2dpX4,
				Scale: borderEffectiveScale,
				ColourMod: colour,
			},
			frenyard.NinePatchFrameLayer{
				Pass: frenyard.FramePassOverBefore,
				NinePatch: borderGenSquareMaskX4,
				Scale: borderEffectiveScale,
				ColourMod: colour,
			},
		},
		Padding:  frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2}),
		Clipping: true,
	}
}
