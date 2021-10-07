package design

import "github.com/yellowsink/frenyard"
import "github.com/yellowsink/frenyard/framework"

var borderGenSquareMaskX4Shadow framework.NinePatch
var borderGenSquareMaskX4Mask framework.NinePatch
var borderGenSquareMaskX4Raw framework.NinePatch
var borderGenRounded4dpMaskX4Shadow framework.NinePatch
var borderGenRounded4dpMaskX4Mask framework.NinePatch
var borderGenRounded4dpMaskX4Raw framework.NinePatch
var borderGenSquareShadow2dpX4Shadow framework.NinePatch
var borderGenSquareShadow2dpX4Mask framework.NinePatch
var borderGenSquareShadow2dpX4Raw framework.NinePatch
var borderGenSquareShadow4dpX4Shadow framework.NinePatch
var borderGenSquareShadow4dpX4Mask framework.NinePatch
var borderGenSquareShadow4dpX4Raw framework.NinePatch
var borderGenRounded4dpShadow2dpX4Shadow framework.NinePatch
var borderGenRounded4dpShadow2dpX4Mask framework.NinePatch
var borderGenRounded4dpShadow2dpX4Raw framework.NinePatch
var borderGenRounded4dpShadow4dpX4Shadow framework.NinePatch
var borderGenRounded4dpShadow4dpX4Mask framework.NinePatch
var borderGenRounded4dpShadow4dpX4Raw framework.NinePatch
var borderGenSquareHeadshadow1dpX4Shadow framework.NinePatch
var borderGenSquareHeadshadow1dpX4Mask framework.NinePatch
var borderGenSquareHeadshadow1dpX4Raw framework.NinePatch
var borderGenSquareHeadshadow2dpX4Shadow framework.NinePatch
var borderGenSquareHeadshadow2dpX4Mask framework.NinePatch
var borderGenSquareHeadshadow2dpX4Raw framework.NinePatch
var borderGenSquareHeadshadowInset1dpX4Shadow framework.NinePatch
var borderGenSquareHeadshadowInset1dpX4Mask framework.NinePatch
var borderGenSquareHeadshadowInset1dpX4Raw framework.NinePatch
var borderGenSquareHeadshadowInset2dpX4Shadow framework.NinePatch
var borderGenSquareHeadshadowInset2dpX4Mask framework.NinePatch
var borderGenSquareHeadshadowInset2dpX4Raw framework.NinePatch
var borderGenRounded4dpHeadshadow1dpX4Shadow framework.NinePatch
var borderGenRounded4dpHeadshadow1dpX4Mask framework.NinePatch
var borderGenRounded4dpHeadshadow1dpX4Raw framework.NinePatch
var borderGenRounded4dpHeadshadow2dpX4Shadow framework.NinePatch
var borderGenRounded4dpHeadshadow2dpX4Mask framework.NinePatch
var borderGenRounded4dpHeadshadow2dpX4Raw framework.NinePatch
var borderGenRounded4dpHeadshadowInset1dpX4Shadow framework.NinePatch
var borderGenRounded4dpHeadshadowInset1dpX4Mask framework.NinePatch
var borderGenRounded4dpHeadshadowInset1dpX4Raw framework.NinePatch
var borderGenRounded4dpHeadshadowInset2dpX4Shadow framework.NinePatch
var borderGenRounded4dpHeadshadowInset2dpX4Mask framework.NinePatch
var borderGenRounded4dpHeadshadowInset2dpX4Raw framework.NinePatch

func deBorderGenInit() {
	borderGenSquareMaskX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareMaskX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareMaskX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpMaskX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 32 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 40 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 44 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpMaskX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 32 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 40 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 44 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpMaskX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 32 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 40 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 44 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareShadow2dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 64 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 72 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 76 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareShadow2dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 64 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 72 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 76 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareShadow2dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 64 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 72 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 76 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareShadow4dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 96 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 104 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 108 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareShadow4dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 96 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 104 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 108 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareShadow4dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 96 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 104 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 108 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpShadow2dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 128 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 136 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 140 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpShadow2dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 128 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 136 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 140 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpShadow2dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 128 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 136 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 140 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpShadow4dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 160 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 168 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 172 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpShadow4dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 160 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 168 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 172 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpShadow4dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 160 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 168 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 172 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadow1dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 192 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 200 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 204 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadow1dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 192 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 200 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 204 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadow1dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 192 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 200 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 204 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadow2dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 224 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 232 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 236 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadow2dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 224 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 232 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 236 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadow2dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 224 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 232 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 236 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadowInset1dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 256 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 264 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 268 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadowInset1dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 256 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 264 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 268 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadowInset1dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 256 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 264 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 268 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadowInset2dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 288 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 296 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 300 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadowInset2dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 288 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 296 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 300 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenSquareHeadshadowInset2dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 288 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 296 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 300 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadow1dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 320 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 328 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 332 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadow1dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 320 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 328 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 332 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadow1dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 320 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 328 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 332 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadow2dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 352 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 360 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 364 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadow2dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 352 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 360 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 364 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadow2dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 352 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 360 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 364 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadowInset1dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 384 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 392 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 396 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadowInset1dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 384 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 392 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 396 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadowInset1dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 384 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 392 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 396 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadowInset2dpX4Shadow = framework.NinePatch{
		Tex:    borderGenTextureShadow,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 416 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 424 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 428 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadowInset2dpX4Mask = framework.NinePatch{
		Tex:    borderGenTextureMask,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 416 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 424 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 428 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
	borderGenRounded4dpHeadshadowInset2dpX4Raw = framework.NinePatch{
		Tex:    borderGenTextureRaw,
		Sprite: frenyard.Area2i{X: frenyard.Area1i{Pos: 416 * borderImageScale, Size: 32 * borderImageScale}, Y: frenyard.Area1i{Pos: 0 * borderImageScale, Size: 32 * borderImageScale}},
		Bounds: frenyard.Area2i{X: frenyard.Area1i{Pos: 424 * borderImageScale, Size: 16 * borderImageScale}, Y: frenyard.Area1i{Pos: 8 * borderImageScale, Size: 16 * borderImageScale}},
		Centre: frenyard.Area2i{X: frenyard.Area1i{Pos: 428 * borderImageScale, Size: 8 * borderImageScale}, Y: frenyard.Area1i{Pos: 12 * borderImageScale, Size: 8 * borderImageScale}},
	}
}

// 14 items (1792px)
