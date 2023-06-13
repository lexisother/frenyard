package design

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/integration"
)

var circleTexReductionsMask []frenyard.Texture
var circleTexReductionsRaw []frenyard.Texture

func deSetupCircle() {
	circleTexReductionsMask = make([]frenyard.Texture, 7)
	circleTexReductionsRaw = make([]frenyard.Texture, 7)
	img := integration.CreateHardcodedPNGImage(circle192)
	for k := range circleTexReductionsMask {
		circleTexReductionsMask[k] = integration.GoImageToTexture(img, []integration.ColourTransform{
			integration.ColourTransformBlueToStencil,
		})
		circleTexReductionsRaw[k] = integration.GoImageToTexture(img, []integration.ColourTransform{})
		if k != len(circleTexReductionsMask)-1 {
			img = integration.ScaleImageToHalfSize(img)
		}
	}
}

func deGetCircleTextureIndex(size int32) int {
	for k, v := range circleTexReductionsMask {
		if v.Size().X <= size {
			return k
		}
	}
	return len(circleTexReductionsMask) - 1
}

func deEncodeCircleCmd(textureSet []frenyard.Texture, position frenyard.Vec2i, size int32, baseCmd frenyard.DrawRectCommand) frenyard.DrawRectCommand {
	baseCmd.Tex = textureSet[deGetCircleTextureIndex(size)]
	baseCmd.TexSprite = frenyard.Area2iOfSize(baseCmd.Tex.Size())
	baseCmd.Target = frenyard.Area2iFromVecs(position, frenyard.Vec2i{}).Align(frenyard.Vec2i{X: size, Y: size}, frenyard.Alignment2i{})
	return baseCmd

}
