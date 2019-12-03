package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

// IconID represents a specific icon.
type IconID int32
// NullIconID should be used when no icon should be placed at all.
const NullIconID IconID = 0
// RunIconID is a right-facing triangle.
const RunIconID IconID = 1
// WarningIconID is a rounded triangle with '!' cut into it.
const WarningIconID IconID = 2
// GameIconID is an icon to represent the game.
const GameIconID IconID = 3
// DirectoryIconID is an icon with the 'folder' style.
const DirectoryIconID IconID = 4
// ModIconID is a gear.
const ModIconID IconID = 5
// ToolIconID is a claw hammer.
const ToolIconID IconID = 6
// BackIconID is a back arrow.
const BackIconID IconID = 7
// DriveIconID is a drive.
const DriveIconID IconID = 8

// Icons at 18DP
var icons18dp frenyard.Texture
var iconHeight18dp int32
// Icons at 24DP
var icons24dp frenyard.Texture
var iconHeight24dp int32
// Icons at 36DP
var icons36dp frenyard.Texture
var iconHeight36dp int32

func deSetupIcons() {
	icons18dpImg := integration.CreateHardcodedPNGImage(icon72B64)
	icons36dpImg := icons18dpImg
	iconHeight18dp = 72
	icons24dpImg := integration.CreateHardcodedPNGImage(icon96B64)
	iconHeight24dp = 96
	// Pass 1: 18 and 24dp (x4)
	iconEffectiveScale := DesignScale / 4
	for (iconEffectiveScale <= 0.5) && (iconHeight18dp > 9) {
		icons18dpImg = integration.ScaleImageToHalfSize(icons18dpImg)
		iconHeight18dp /= 2
		icons24dpImg = integration.ScaleImageToHalfSize(icons24dpImg)
		iconHeight24dp /= 2
		iconEffectiveScale *= 2
	}
	// Pass 2: 36dp (x2)
	iconHeight36dp = 72
	iconEffectiveScale = DesignScale / 2
	for (iconEffectiveScale <= 0.5) && (iconHeight36dp > 9) {
		icons36dpImg = integration.ScaleImageToHalfSize(icons36dpImg)
		iconHeight36dp /= 2
		iconEffectiveScale *= 2
	}
	icons18dp = integration.GoImageToTexture(icons18dpImg, []integration.ColourTransform{integration.ColourTransformBlueToStencil});
	icons24dp = integration.GoImageToTexture(icons24dpImg, []integration.ColourTransform{integration.ColourTransformBlueToStencil});
	icons36dp = integration.GoImageToTexture(icons36dpImg, []integration.ColourTransform{integration.ColourTransformBlueToStencil});
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
	} else if iconSizeDP == 36 {
		spriteSize = iconHeight36dp
		tex = icons36dp
	} else {
		panic("Unsupported icon height (check NewIconPtr calls)")
	}
	spriteSize2 := frenyard.Vec2i{X: spriteSize, Y: spriteSize}
	index := int32(iconID) - 1
	sprite := frenyard.Area2iFromVecs(frenyard.Vec2i{X: spriteSize * (index % 16), Y: spriteSize * (index / 16)}, spriteSize2)
	return framework.ConvertElementToLayout(framework.NewTextureRectPtr(colour, tex, sprite, spriteSize2, frenyard.Alignment2i{}))
}
