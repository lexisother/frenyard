## Submodules

Submodules are 2-letter prefixes meant to further discourage interference.

`<submodule>` is `fy` for `frenyard`.

`<submodule>` is `de` for `CCUpdaterUI/design`.

`<submodule>` is `up` for `CCUpdaterUI` itself.

## Terminology

A "mathematical struct" is any struct intended solely as a mathematical grouping.

It is thus not subject to these rules because it does not concern itself with inheritance, instead being annotated with helper methods for ease of creation in whatever form is simplest for the user at the time, and structured for ease of processing.

An "embeddable struct" is a struct that may be, or is designed to be, placed anonymously into another struct.

An "unembeddable struct" is any struct that is explicitly not an embeddable struct.

A "formal initializer" is a function intended as the structure's primary initializer. Helper methods do not count.

The "outer struct" is the struct which embeds an embeddable struct.

## Initialization

Note: `<Variant>` here may be omitted entirely. If present, it describes the 'variant' of constructor, which may be important in cases like UIRect.

In any non-mathematical structure, if a formal initializer is present, it must always be used instead of the zero value, except for `Init<Variant><Struct>`.

There are three kinds of formal initializer:

`New<Variant><Struct>` : Function that returns `Struct`. Makes the structure potentially embeddable.

`New<Variant><Struct>Ptr` : Function that returns `*Struct`. Makes the structure explicitly unembeddable.

`Init<Variant><Struct>` : The (explicitly embedded) structure starts with the zero value, and is then initialized in-place using a pointer to the outer struct 'drilling down' using a method implemented by the embedded struct.

They may have `<Struct>From<Variant>` or such, these are simply helpers and do not count as the initializer.

## Private Fields / Methods

Fields that are so private only the struct itself (and very specifically authorized related functions) can poke at them must be in one of two places:

1. Within embeddable structs, `_<submodule><Class><Field>` (such as `_fyEventDebuggerMousePos`)
2. Within unembeddable structs, `_<field>`.

In practice, there is another option for embedded structs described later in this document that follows these two rules.

As Go does not have struct layout versioning issues, *IT IS NOT CONSIDERED A BREAKING CHANGE TO ALTER PRIVATE FIELDS*.

Note that manually initializing or specifying the field values on creation of private fields counts as private access and is thus undefined behavior for the purposes of if a change is breaking.

## File-Private Functions/Types/Fields

NOTE: THIS IS NOT FINALIZED YET! There's plenty of stuff here that requires smoothing over and understanding.

Functions & Types that are private to the file: `<submodule><Name>`

To indicate it is not private to the file, do not use the submodule name.

Do not confuse with `_<submodule><Name>` (struct-local)

TODO: Put all of this into a "Safety" supersection?

## Interfaces & Public Interface Prefixes

Public Interface Prefixes are used to help prevent collisions, and follow the scheme:

`<Submodule><I>` where `<I>` is some set of interface initial(s).

Where a function exists to implement a function from an interface, the comment to satiate golint must be:

`// <Name> implements <Interface>.<Name>`

The interface name chosen must be from the base interface, not from interfaces that extend it.

## Embeddable Structs, Compositional Inheritance, Supercalls

Sometimes a struct is intended to be included inside another struct to achieve composition-based inheritance.

The naming for these is not standard because the details of a given struct are specific to the situation.

To prevent breakage, it is recommended that these structures prefer the `Init` pattern over `New` or zero-values.

Where a function exists to override a function implemented elsewhere, the comment to satiate golint must be:

`// <Name> overrides <Struct>.<Name>`

The structure that is chosen to represent the overridden function, and thus the structure chosen for supercalls, must be the *immediately* embedded struct. This is in order to prevent issues if an override is introduced later.

It is NOT considered breaking to introduce an override; if the rules are properly followed, this override will act "normally" as per Java-style rules.

## ThisStructDetails Structs

Sometimes, when writing an embeddable struct, sometimes only the outer struct should have access to something.

In this case, create an unembeddable struct within the embeddable struct.

All private fields should be placed within the unembeddable struct (using the naming rules for unembeddable structs) for code brevity and to lower the likelihood the unembeddable struct will need a pointer to the embeddable struct.

The field name should be `This<Struct>Details` - the struct name should be `<Struct>Details`.

This should only need to happen if the structure is intended to be inherited.

## Assorted Notes

All UI elements (implements and/or is UIElement) must start with `UI`. Things that are not UI elements must not.
(Note: In conflict with the 'Details' postfix rule, the 'Details' postfix rule wins.)

To prevent accidents, all structures intended for direct use without embedding must be unembeddable (thus "final").

To then 'extend' these structures anyway, a form of proxy should be used instead.

This is particularly important for layout classes such as UIFlexboxContainer; the choice of layout is arbitrary but would end up influencing the required "super-calls" all over the struct using it.

For the SDL2 backend, functions called within the OS-thread must start with `os`.

## Putting This Together

An example of some of these concepts:

```go
package main

import "fmt"

// -- actions
type Actionable interface {
	// Performs the action for this object.
	ExAAction()
}

// -- names
type Namable interface {
	// Returns the noun for the object.
	ExNName() string
}

// -- verbactions

// Note: This part of the structure is open to adaptation (UILayoutElement is merged with what would be UILayoutElementComponentHost, which helps UILayoutProxy 'skip over' redundant elements)
// Implements Actionable on a Namable by outputting an action of the form "*Verb <name>*" to console.
type VerbActionableComponent struct {
	// Accessable by the owner of this component.
	ThisVerbActionableComponent VerbActionableComponentDetails
}

// Accessable by the owner of this component, this contains the Verb.
type VerbActionableComponentDetails struct {
	// Verb that is written when the action is performed. Defaults to "pokes".
	Verb string
	// Private field, do not touch outside of VerbActionableComponent[Details]
	_host VerbActionableComponentHost
}

// Interface for a struct that contains VerbActionableComponent and implements Namable; required to initialize VerbActionableComponent.
type VerbActionableComponentHost interface {
	Namable
	_exGetVerbActionableComponent() *VerbActionableComponent
}

// Private function used when initializing to retrieve the VerbActionableComponent being initialized.
func (fac *VerbActionableComponent) _exGetVerbActionableComponent() *VerbActionableComponent {
	return fac
}

// Implements Actionable.ExAAction
func (fac *VerbActionableComponent) ExAAction() {
	fmt.Printf("*%v %v*\n", fac.ThisVerbActionableComponent.Verb, fac.ThisVerbActionableComponent._host.ExNName())
}

// Initializes the VerbActionableComponent
func InitVerbActionableComponent(f VerbActionableComponentHost) {
	fac := f._exGetVerbActionableComponent()
	// This binds the component to it's host so that it can access Namable.
	// This is the reason behind _exGetVerbActionableComponent and VerbActionableComponentHost.
	fac.ThisVerbActionableComponent._host = f
	fac.ThisVerbActionableComponent.Verb = "pokes"
}

// -- orange

// Orange (Uses VerbActionableComponent to implement Actionable, which requires that the Orange implement Namable)
type Orange struct {
	VerbActionableComponent
}

func (o *Orange) ExNName() string {
	return "orange"
}

func NewOrangePtr() *Orange {
	base := &Orange{}
	// Initialize the component. Note that Orange implements Namable and the VerbActionableComponent implements the _exGetVerbActionableComponent function.
	InitVerbActionableComponent(base)
	// With that done, use the ThisVerbActionableComponent field to alter the verb.
	base.ThisVerbActionableComponent.Verb = "eats"
	return base
}

// -- main

func main() {
	NewOrangePtr().ExAAction() // "*eats orange*"
}
```
