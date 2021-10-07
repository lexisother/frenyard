package framework

import (
	"fmt"
	"github.com/yellowsink/frenyard"
	"sort"
)

// Implements a highly limited subset of flexbox to be extended to full support as-needed.

// FlexboxWrapMode describes a type of wrapping mode for Flexbox containers.
type FlexboxWrapMode uint8

// FlexboxWrapModeNone disallows wrapping for items, they are all on one line.
const FlexboxWrapModeNone FlexboxWrapMode = 0

// FlexboxWrapModeWrap allows items to wrap between lines.
const FlexboxWrapModeWrap FlexboxWrapMode = 1

// FlexboxContainer describes a UIFlexboxContainer's contents.
type FlexboxContainer struct {
	DirVertical bool
	WrapMode    FlexboxWrapMode
	// Ignored when used by the line solver; it uses fyFlexboxSlotlike instead
	Slots []FlexboxSlot

	// This is NOT part of the public API. It is provided for debugging use only.
	// Do not even mention this field in releases.
	Debug bool
}

// FlexboxSlot describes an element within a Flexbox container.
type FlexboxSlot struct {
	// Can be nil.
	Element UILayoutElement
	// If there is a surplus, these are used to distribute it.
	Grow int32
	// If there is a deficit, these are used to distribute it (along with minimum sizes)
	// DEFAULTS TO 1 IN CSS
	Shrink int32
	// If *non-zero*, then this specifies the "initial share size" of this element.
	// Useful when Element is nil.
	Basis int32
	// Slightly non-standard extension (or is it?) for cases where Basis would be used to pad a problematic element
	MinBasis int32
	// Used to order the flexboxes visually. The Z-Order remains the index order.
	Order int
	// If the container should respect the minimum size of this slot.
	// DEFAULTS TO TRUE IN CSS (settable to false by overriding min w/h), BUT THIS IS REALLY STUPID B/C IT IGNORES SIZE CONSTRAINTS, SO LET'S NOT DO THAT
	RespectMinimumSize bool
}

type fyFlexboxSlotlike interface {
	fyMainCrossSizeForMainCrossLimits(limits frenyard.Vec2i, vertical bool, debug bool) frenyard.Vec2i
	fyGrowShrink() (int32, int32)
	fyCalcBasis(cross int32, vertical bool) int32
	fyGetOrder() int
	fyRespectMinimumSize() bool
}

func (slot FlexboxSlot) fyMainCrossSizeForMainCrossLimits(limits frenyard.Vec2i, vertical bool, debug bool) frenyard.Vec2i {
	if slot.Element == nil {
		return frenyard.Vec2i{}
	}
	if debug {
		fmt.Print("?")
	}
	mainCrossSize := slot.Element.FyLSizeForLimits(limits.ConditionalTranspose(vertical)).ConditionalTranspose(vertical)
	mainCrossSize.X = frenyard.Max(mainCrossSize.X, slot.MinBasis)
	return mainCrossSize
}
func (slot FlexboxSlot) fyGrowShrink() (int32, int32) {
	return slot.Grow, slot.Shrink
}
func (slot FlexboxSlot) fyCalcBasis(cross int32, vertical bool) int32 {
	if slot.Basis != 0 {
		return slot.Basis
	}
	return slot.fyMainCrossSizeForMainCrossLimits(frenyard.Vec2i{frenyard.SizeUnlimited, cross}, vertical, false).X
}
func (slot FlexboxSlot) fyGetOrder() int {
	return slot.Order
}
func (slot FlexboxSlot) fyRespectMinimumSize() bool {
	return slot.RespectMinimumSize
}

// -- Solver --

func fyFlexboxGetPreferredSize(details FlexboxContainer) frenyard.Vec2i {
	// Do note, this is in main/cross format.
	mainCrossSize := frenyard.Vec2i{}
	for _, v := range details.Slots {
		sz := v.fyMainCrossSizeForMainCrossLimits(frenyard.Vec2iUnlimited(), details.DirVertical, false)
		mainCrossSize.X += sz.X
		mainCrossSize.Y = frenyard.Max(mainCrossSize.Y, sz.Y)
	}
	return mainCrossSize.ConditionalTranspose(details.DirVertical)
}

type fyFlexboxRow struct {
	elem     []fyFlexboxSlotlike
	area     []frenyard.Area2i
	fullArea frenyard.Area2i
}

func (slot fyFlexboxRow) fyGrowShrink() (int32, int32) {
	return 1, 1
}

// Critical to the whole thing and it's full of guesswork due to the vertical flags and axis juggling.
func (slot fyFlexboxRow) fyMainCrossSizeForMainCrossLimits(limits frenyard.Vec2i, vertical bool, debug bool) frenyard.Vec2i {
	if debug {
		fmt.Print("R{")
	}
	// Main & Cross in here refer to in the row flexbox, not the outer flexbox.
	maximumMain := int32(0)
	presentAreaCross := slot.fullArea.Size().ConditionalTranspose(vertical).Y
	for _, v := range slot.elem {
		lim := frenyard.Vec2i{X: limits.X, Y: presentAreaCross}
		if debug {
			fmt.Print(" ", limits.X, "x", presentAreaCross)
		}
		rcs := v.fyMainCrossSizeForMainCrossLimits(lim, vertical, false)
		maximumMain = frenyard.Max(maximumMain, rcs.X)
		if debug {
			fmt.Print(":", rcs.X, "x", rcs.Y)
		}
	}
	if debug {
		fmt.Print(" }")
	}
	return frenyard.Vec2i{X: maximumMain, Y: presentAreaCross}
}
func (slot fyFlexboxRow) fyCalcBasis(cross int32, vertical bool) int32 {
	return slot.fyMainCrossSizeForMainCrossLimits(frenyard.Vec2i{X: frenyard.SizeUnlimited, Y: cross}, vertical, false).X
}
func (slot fyFlexboxRow) fyGetOrder() int {
	return 0
}
func (slot fyFlexboxRow) fyRespectMinimumSize() bool {
	return false
}

// Do be aware, this only handles the one relevant axis.
func (slot *fyFlexboxRow) Fill(area frenyard.Area2i, vertical bool) {
	for k := range slot.area {
		if !vertical {
			// Rows perpendicular to X
			slot.area[k].X = area.X
		} else {
			// Rows perpendicular to Y
			slot.area[k].Y = area.Y
		}
	}
	slot.fullArea = area
}

type fyFlexboxSortingCollection struct {
	// The collection being sorted.
	slots []fyFlexboxSlotlike
	// Given a SOURCE slot index, what is the RESULTING slot index?
	originalToDisplayIndices []int
	// Given a RESULTING slot index, what is the SOURCE slot index?
	displayToOriginalIndices []int
}

func (sc fyFlexboxSortingCollection) Len() int {
	return len(sc.slots)
}
func (sc fyFlexboxSortingCollection) Less(i int, j int) bool {
	order1 := sc.slots[i].fyGetOrder()
	order2 := sc.slots[j].fyGetOrder()
	// Order1 != order2?
	if order1 < order2 {
		return true
	}
	if order1 > order2 {
		return false
	}
	// No, they're equal. Sort by original index.
	if sc.displayToOriginalIndices[i] < sc.displayToOriginalIndices[j] {
		return true
	}
	return false
}
func (sc fyFlexboxSortingCollection) Swap(i int, j int) {
	backup := sc.slots[i]
	backup2 := sc.originalToDisplayIndices[i]
	backup3 := sc.displayToOriginalIndices[i]

	sc.slots[i] = sc.slots[j]
	sc.originalToDisplayIndices[i] = sc.originalToDisplayIndices[j]
	sc.displayToOriginalIndices[i] = sc.displayToOriginalIndices[j]

	sc.slots[j] = backup
	sc.originalToDisplayIndices[j] = backup2
	sc.displayToOriginalIndices[j] = backup3
}

func fyFlexboxSolveLayout(details FlexboxContainer, limits frenyard.Vec2i) []frenyard.Area2i {
	// Stage 1. Element order pre-processing (DirReverse)
	slots := make([]fyFlexboxSlotlike, len(details.Slots))
	originalToDisplayIndices := make([]int, len(details.Slots))
	displayToOriginalIndices := make([]int, len(details.Slots))
	for k, v := range details.Slots {
		originalToDisplayIndices[k] = k
		displayToOriginalIndices[k] = k
		slots[k] = v
	}
	sort.Sort(fyFlexboxSortingCollection{
		slots:                    slots,
		originalToDisplayIndices: originalToDisplayIndices,
		displayToOriginalIndices: displayToOriginalIndices,
	})
	// Stage 2. Wrapping (if relevant)
	out := make([]frenyard.Area2i, len(slots))
	mainCrossLimits := limits.ConditionalTranspose(details.DirVertical)
	shouldWrap := fyFlexboxSolveLine(details, slots, out, mainCrossLimits, details.Debug)
	// One row, so this is simple
	rows := []fyFlexboxRow{{slots, out, frenyard.UnionArea2i(out)}}
	if shouldWrap && details.WrapMode != FlexboxWrapModeNone {
		// Wrapping has to start. Oh no...
		// Do note, lines is implicitly limited because of the "one slot cannot wrap" rule.
		lines := int32(2)
		for {
			rows = make([]fyFlexboxRow, lines)
			lineStartSlot := 0
			consumedSlots := 0

			currentLine := int32(0)
			for consumedSlots < len(slots) {
				// If it wraps...
				if fyFlexboxSolveLine(details, slots[lineStartSlot:consumedSlots+1], out[lineStartSlot:consumedSlots+1], mainCrossLimits, false) {
					// Revert it & finish the line.
					rows[currentLine] = fyFlexboxRow{
						slots[lineStartSlot:consumedSlots],
						out[lineStartSlot:consumedSlots],
						frenyard.UnionArea2i(out),
					}
					fyFlexboxSolveLine(details, rows[currentLine].elem, rows[currentLine].area, mainCrossLimits, false)
					// Now setup the new line.
					currentLine++
					lineStartSlot = consumedSlots
					if currentLine == lines {
						// Out of range, cancel before rows[currentLine] brings it to a halt
						break
					}
					// Retry the same slot (slot not consumed)
				} else {
					// Success! Advance.
					consumedSlots++
				}
			}
			if currentLine < lines {
				// Finish last line
				rows[currentLine] = fyFlexboxRow{
					slots[lineStartSlot:consumedSlots],
					out[lineStartSlot:consumedSlots],
					frenyard.UnionArea2i(out),
				}
				break
			}
			lines++
		}
	}
	if details.WrapMode != FlexboxWrapModeNone {
		// Stage 3. Row compression
		rowAreas := make([]frenyard.Area2i, len(rows))
		rowSlots := make([]fyFlexboxSlotlike, len(rows))
		for rk, row := range rows {
			rowSlots[rk] = row
		}
		fyFlexboxSolveLine(FlexboxContainer{
			DirVertical: !details.DirVertical,
			WrapMode:    FlexboxWrapModeNone,
		}, rowSlots, rowAreas, frenyard.Vec2i{mainCrossLimits.Y, mainCrossLimits.X}, false)
		for rk, row := range rows {
			row.Fill(rowAreas[rk], !details.DirVertical)
		}
	} else {
		// Stage 3. Row setup
		if mainCrossLimits.Y != frenyard.SizeUnlimited {
			rows[0].Fill(frenyard.Area2iOfSize(mainCrossLimits.ConditionalTranspose(details.DirVertical)), !details.DirVertical)
		}
	}
	// Stage 4. Element order post-processing (DirReverse)
	realOutput := make([]frenyard.Area2i, len(out))
	for k, v := range originalToDisplayIndices {
		realOutput[k] = out[v]
	}
	return realOutput
}

// Returns true if should wrap. Will not return true ever for only one slot as this cannot wrap.
func fyFlexboxSolveLine(details FlexboxContainer, slots []fyFlexboxSlotlike, out []frenyard.Area2i, mainCrossLimits frenyard.Vec2i, debug bool) bool {
	if len(slots) == 0 {
		// Nowhere to output. Also, some calculations rely on at least one slot existing.
		return false
	}
	if debug {
		if details.DirVertical {
			fmt.Print("VERTICAL ")
		}
		fmt.Println("AREA", mainCrossLimits.X, "x", mainCrossLimits.Y)
	}
	// Substage 1. Input basis values & create total
	shares := make([]int32, len(slots))
	totalMainAccumulator := int32(0)
	totalGrowAccumulator := int32(0)
	totalShrinkAccumulator := int32(0)
	for idx, slot := range slots {
		shares[idx] = slot.fyCalcBasis(mainCrossLimits.Y, details.DirVertical)
		totalMainAccumulator += shares[idx]
		slotGrow, slotShrink := slot.fyGrowShrink()
		totalGrowAccumulator += slotGrow
		totalShrinkAccumulator += slotShrink
	}
	// Notably, totalMainAccumulator must not change after this point.
	// It's the 'reference' for if we ought to wrap.
	// Substage 2. Determine expansion or contraction
	if mainCrossLimits.X != frenyard.SizeUnlimited && totalMainAccumulator != mainCrossLimits.X {
		additionalSpaceAvailable := mainCrossLimits.X - totalMainAccumulator
		if debug {
			fmt.Println("COMPRESSOR II: ", additionalSpaceAvailable)
		}
		// Determine which accumulator to use.
		totalFactorAccumulator := totalGrowAccumulator
		if additionalSpaceAvailable < 0 {
			totalFactorAccumulator = totalShrinkAccumulator
		}
		// Actually redistribute space. This may require multiple passes as a factor may not always be fully appliable.
		// This is because the Flexbox system respects minimum size.
		// When set to true, the relevant factor must be subtracted from the Accumulator.
		slotsHitMinimumSize := make([]bool, len(slots))
		needAnotherPass := true
		for needAnotherPass && totalFactorAccumulator != 0 {
			needAnotherPass = false
			totalAlloc := int32(0)
			for idx, slot := range slots {
				if slotsHitMinimumSize[idx] {
					continue
				}
				grow, shrink := slot.fyGrowShrink()
				factor := grow
				smallestAlloc := -shares[idx] // Cannot shrink below 0
				largestAlloc := frenyard.SizeUnlimited
				// There is no 'largest alloc'; if the element is told to grow, that is what it will do
				if additionalSpaceAvailable < 0 {
					factor = shrink
				}
				if factor == 0 {
					// has no effect, and means totalFactorAccumulator could be 0
					continue
				}
				if additionalSpaceAvailable < 0 && shrink > 0 && slot.fyRespectMinimumSize() {
					// Smallest possible alloc: maximum amount that can be shrunk
					smallestAlloc = slot.fyMainCrossSizeForMainCrossLimits(frenyard.Vec2i{X: 0, Y: mainCrossLimits.Y}, details.DirVertical, false).X - shares[idx]
				}
				alloc := (additionalSpaceAvailable * factor) / totalFactorAccumulator
				// Limit allocation.
				clamped := false
				if alloc <= smallestAlloc {
					alloc = smallestAlloc
					clamped = true
				}
				if alloc >= largestAlloc {
					alloc = largestAlloc
					clamped = true
				}
				// If the limit is hit, remove from processing for the next loop.
				if clamped {
					slotsHitMinimumSize[idx] = true
					needAnotherPass = true
					totalFactorAccumulator -= factor
				}
				// Confirm allocation
				shares[idx] += alloc
				totalAlloc += alloc
			}
			additionalSpaceAvailable -= totalAlloc
		}
		// additionalSpaceAvailable non-zero: justify-content implementation goes here
	}
	// Substage 3. With horizontal sizes established, calculate crossLimit
	crossLimit := int32(0)
	for idx := 0; idx < len(slots); idx++ {
		crossLimit = frenyard.Max(crossLimit, slots[idx].fyMainCrossSizeForMainCrossLimits(frenyard.Vec2i{X: shares[idx], Y: mainCrossLimits.Y}, details.DirVertical, false).Y)
	}
	// -- Actual layout! For real this time! --
	mainPosition := int32(0)
	if debug {
		fmt.Println(" CROSS ", crossLimit)
	}
	for idx := 0; idx < len(slots); idx++ {
		out[idx] = frenyard.Area2iOfSize(frenyard.Vec2i{shares[idx], crossLimit}.ConditionalTranspose(details.DirVertical)).Translate(frenyard.Vec2i{mainPosition, 0}.ConditionalTranspose(details.DirVertical))
		if debug {
			fmt.Println(" SHARE ", shares[idx])
		}
		mainPosition += shares[idx]
	}
	if debug {
		fmt.Println("END AREA")
	}
	// If len(slots) <= 1 then wrapping would inf. loop, so only wrap for >1.
	return (len(slots) > 1) && (totalMainAccumulator > mainCrossLimits.X)
}

// -- UI element --

// UIFlexboxContainer lays out UILayoutElements using a partial implementation of Flexbox.
type UIFlexboxContainer struct {
	UIPanel
	UILayoutElementComponent
	_state         FlexboxContainer
	_preferredSize frenyard.Vec2i
}

// NewUIFlexboxContainerPtr creates a UIFlexboxContainer from the FlexboxContainer details
func NewUIFlexboxContainerPtr(setup FlexboxContainer) *UIFlexboxContainer {
	container := &UIFlexboxContainer{
		UIPanel: NewPanel(frenyard.Vec2i{}),
	}
	InitUILayoutElementComponent(container)
	container.SetContent(setup)
	container.FyEResize(container._preferredSize)
	return container
}

// FyLSubelementChanged implements UILayoutElement.FyLSubelementChanged
func (ufc *UIFlexboxContainer) FyLSubelementChanged() {
	ufc._preferredSize = fyFlexboxGetPreferredSize(ufc._state)
	ufc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (ufc *UIFlexboxContainer) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	if limits.Ge(ufc._preferredSize) {
		return ufc._preferredSize
	}
	solved := fyFlexboxSolveLayout(ufc._state, limits)
	max := frenyard.Vec2i{}
	for _, v := range solved {
		max = max.Max(v.Pos().Add(v.Size()))
	}
	return max
}

// SetContent changes the contents of the UIFlexboxContainer.
func (ufc *UIFlexboxContainer) SetContent(setup FlexboxContainer) {
	if ufc._state.Slots != nil {
		for _, v := range ufc._state.Slots {
			if v.Element != nil {
				ufc.ThisUILayoutElementComponentDetails.Detach(v.Element)
			}
		}
	}
	ufc._state = setup
	for _, v := range setup.Slots {
		if v.Element != nil {
			ufc.ThisUILayoutElementComponentDetails.Attach(v.Element)
		}
	}
	ufc.FyLSubelementChanged()
}

// FyEResize overrides UIPanel.FyEResize
func (ufc *UIFlexboxContainer) FyEResize(size frenyard.Vec2i) {
	ufc.UIPanel.FyEResize(size)
	areas := fyFlexboxSolveLayout(ufc._state, size)
	fixes := make([]PanelFixedElement, len(areas))
	fixesCount := 0
	for idx, slot := range ufc._state.Slots {
		if slot.Element != nil {
			fixes[fixesCount].Pos = areas[idx].Pos()
			fixes[fixesCount].Visible = true
			fixes[fixesCount].Element = slot.Element
			slot.Element.FyEResize(areas[idx].Size())
			fixesCount++
		}
	}
	ufc.ThisUIPanelDetails.SetContent(fixes[:fixesCount])
}
