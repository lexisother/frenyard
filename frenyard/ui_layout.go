package frenyard

/*
 *  THE FULL REASONING BEHIND THE DESIGN
 * Let's assume a label changes in goodness knows where.
 * That label runs ContentsChanged. This in turn calls the parent FyLSubelementLimitsChange.
 * The parent has to assume the parent's limits have changed because the label's limits have changed.
 * This continues upwards until it hits either the root or a non-UILayoutElement parent.
 * With no parent to pass the buck to, it goes downwards.
 */

/*
 * Wrapping structure around Element.
 * Please note: When implementing this, you are expected to:
 * 1. React to FyEResize with a re-layout.
 *  This applies even if the FyEResize is redundant.
 * 2. When re-layouting, always FyEResize every sub-element.
 * 3. DO NOT signal you've changed content unless the content changed in a way that would affect limits!!!
 *    Otherwise, instead update your component some other way.
 *    Your limits ARE NOT, EVER, size-dependent!!!
 *    I've learned the hard way that this destabilizes in all sorts of bad ways.
 *    Instead the limits are flexible enough to specify the intended meaning directly.
 *    If you're something like a label, your limits change with your.
 *    If you're a panel layout whose limits are dependent on subelements, your limits have changed when the limits of your sub-elements have changed.
 *    As such update any cached limits from FyLSubelementChanged and signal ContentsChanged().
 *    The complexity of this is more or less equivalent to the tree depth.
 *    (There is room for optimization here by delaying ContentsChanged in the right way to cause many changes to execute in a single pass, but this requires an ElementHost interface to coordinate properly.)
 */

// UILayoutElement is the base type for all UIElements that support layout. All implementations must use UILayoutElementComponent or UILayoutProxy.
type UILayoutElement interface {
	UIElement
	/*
	 * Used to allow Details components to communicate with each other.
	 * Do not override from the Component, do not access.
	 * If and only if a proxy is in use, this may actually be the details of the "wrong" element,
	 *  because the point of proxies is to lower resource usage over a full panel.
	 */
	_fyGetUILayoutElementComponent() *UILayoutElementComponent
	/*
	 * Propagation of subelement limits change.
	 */
	FyLSubelementChanged()
	/*
	 * Gets the optimum size given limitations.
	 * The limitations use SIZE_UNLIMITED where unlimited.
	 * The result may exceed the limitations given.
	 * Note that this MUST NOT CHANGE given the same content.
	 */
	FyLSizeForLimits(limits Vec2i) Vec2i
}

// "Details", not "Private". Only access from the component's user.
// Additionally, note that the empty struct needs to be just about valid for a second.

// UILayoutElementComponentDetails is a holder for methods meant to be accessed by the owner struct of UILayoutElementComponent.
type UILayoutElementComponentDetails struct {
	_parent UILayoutElement
	_self   UILayoutElement
}

// UILayoutElementComponent implements layout notification sending and receiving.
type UILayoutElementComponent struct {
	// "This", not "Private". Only access from the component holder.
	ThisUILayoutElementComponentDetails UILayoutElementComponentDetails
}

// InitUILayoutElementComponent initializes a UILayoutElementComponent given the containing UILayoutElement. Do not call on UILayoutProxy or feed after midnight.
func InitUILayoutElementComponent(self UILayoutElement) {
	ulec := self._fyGetUILayoutElementComponent()
	if ulec.ThisUILayoutElementComponentDetails._self != nil {
		panic("UILayoutElementComponent was initialized twice. It should only be initialized by the initializer of the struct that includes it.")
	} else {
		ulec.ThisUILayoutElementComponentDetails._self = self
	}
}
func (cmp *UILayoutElementComponent) _fyGetUILayoutElementComponent() *UILayoutElementComponent {
	return cmp
}

// LayoutElementNoSubelementsComponent should be used when your element has no subelements.
type LayoutElementNoSubelementsComponent struct {
}

// FyLSubelementChanged implements UILayoutElement.FyLSubelementChanged
func (px *LayoutElementNoSubelementsComponent) FyLSubelementChanged() {
	panic("This mustn't actually be called since the element never gets any attachments. If it is called something went very, very wrong.")
}

// ContentChanged alerts the UI layouter that the limits of this component changed.
func (details *UILayoutElementComponentDetails) ContentChanged() {
	if details._parent != nil {
		details._parent.FyLSubelementChanged()
	} else {
		if details._self == nil {
			panic("UILayoutElementComponent was not properly initialized. Please use InitUILayoutElementComponent(self) on your structure.")
		}
		details._self.FyEResize(details._self.FyESize())
	}
}

// ConvertElementToLayout converts a UIElement to a UILayoutElement.
func ConvertElementToLayout(subelement UIElement) UILayoutElement {
	var e UILayoutElement
	switch e2 := subelement.(type) {
	case UILayoutElement:
		e = e2
	default:
		ev2 := &fyAdaptedLayoutElement{UILayoutElementComponent{}, LayoutElementNoSubelementsComponent{}, UIProxy{subelement}, subelement.FyESize()}
		InitUILayoutElementComponent(ev2)
		e = ev2
	}
	return e
}

// Attach attaches an element, so that if it changes, the parent will be notified.
func (details *UILayoutElementComponentDetails) Attach(subelement UILayoutElement) UILayoutElement {
	otherDetails := subelement._fyGetUILayoutElementComponent()
	if otherDetails.ThisUILayoutElementComponentDetails._parent != nil {
		panic("Double-attachment is a logical error that will completely blow up the application")
	}
	otherDetails.ThisUILayoutElementComponentDetails._parent = details._self
	return subelement
}

// Detach detaches a previously attached element.
func (details *UILayoutElementComponentDetails) Detach(subelement UILayoutElement) {
	otherDetails := subelement._fyGetUILayoutElementComponent()
	if otherDetails.ThisUILayoutElementComponentDetails._parent != details._self {
		panic("Tried to detach element that wasn't attached here in the first place")
	}
	otherDetails.ThisUILayoutElementComponentDetails._parent = nil
}

// Used to convert things
type fyAdaptedLayoutElement struct {
	UILayoutElementComponent
	LayoutElementNoSubelementsComponent
	UIProxy
	CacheSize Vec2i
}

func (ale *fyAdaptedLayoutElement) FyLSizeForLimits(limits Vec2i) Vec2i {
	return ale.CacheSize
}

// UILayoutProxy is UIProxy with layout support.
type UILayoutProxy struct {
	UIProxy
	LayoutElementNoSubelementsComponent
}

// InitUILayoutProxy initializes a UILayoutProxy, setting the target.
func InitUILayoutProxy(proxy UIProxyHost, target UILayoutElement) {
	InitUIProxy(proxy, target)
}

// _fyGetUILayoutElementComponent implements UILayoutElement._fyGetUILayoutElementComponent
func (px *UILayoutProxy) _fyGetUILayoutElementComponent() *UILayoutElementComponent {
	return px.fyUIProxyTarget.(UILayoutElement)._fyGetUILayoutElementComponent()
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (px *UILayoutProxy) FyLSizeForLimits(limits Vec2i) Vec2i {
	return px.fyUIProxyTarget.(UILayoutElement).FyLSizeForLimits(limits)
}
