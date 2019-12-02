package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

// IconID represents a specific icon.
type IconID int32
// NullIconID should be used when no icon should be placed at all.
const NullIconID IconID = -1
// RunIconID is a right-facing triangle.
const RunIconID IconID = 0
// WarningIconID is a rounded triangle with '!' cut into it.
const WarningIconID IconID = 1
// GameIconID is an icon to represent the game.
const GameIconID IconID = 2
// DirectoryIconID is an icon with the 'folder' style.
const DirectoryIconID IconID = 3
// ModIconID is a gear.
const ModIconID IconID = 4
// ToolIconID is a claw hammer.
const ToolIconID IconID = 5

// Icons at 18DP
var icons18dp frenyard.Texture
var iconHeight18dp int32
// Icons at 24DP
var icons24dp frenyard.Texture
var iconHeight24dp int32

func deSetupIcons() {
	icons18dpImg := integration.CreateHardcodedPNGImage(icon72B64)
	iconHeight18dp = 72
	icons24dpImg := integration.CreateHardcodedPNGImage(icon96B64)
	iconHeight24dp = 96
	iconEffectiveScale := DesignScale / 4
	for (iconEffectiveScale <= 0.5) && (iconHeight18dp > 9) {
		icons18dpImg = integration.ScaleImageToHalfSize(icons18dpImg)
		iconHeight18dp /= 2
		icons24dpImg = integration.ScaleImageToHalfSize(icons24dpImg)
		iconHeight24dp /= 2
		iconEffectiveScale *= 2
	}
	icons18dp = integration.GoImageToTexture(icons18dpImg, []integration.ColourTransform{integration.ColourTransformBlueToStencil});
	icons24dp = integration.GoImageToTexture(icons24dpImg, []integration.ColourTransform{integration.ColourTransformBlueToStencil});
}

// NewIconPtr returns a UILayoutElement for an icon at the given DP size. (Note: Only select values are supported.)
func NewIconPtr(colour uint32, iconID IconID, iconSizeDP int32) framework.UILayoutElement {
	var spriteSize int32
	var tex frenyard.Texture
	if iconSizeDP == 18 {
		spriteSize = iconHeight18dp
		tex = icons18dp
	} else if iconSizeDP == 24 {
		spriteSize = iconHeight24dp
		tex = icons24dp
	} else {
		panic("Unsupported icon height (check NewIconPtr calls)")
	}
	spriteSize2 := frenyard.Vec2i{X: spriteSize, Y: spriteSize}
	sprite := frenyard.Area2iFromVecs(frenyard.Vec2i{X: spriteSize * (int32(iconID) % 16), Y: spriteSize * (int32(iconID) / 16)}, spriteSize2)
	return framework.ConvertElementToLayout(framework.NewTextureRectPtr(colour, tex, sprite, spriteSize2, frenyard.Alignment2i{}))
}
