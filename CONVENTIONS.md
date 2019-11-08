# Not really conventional

`<submodule>` is `fy` for Frenyard

Fields that are so private even inherited structs shouldn't poke at them must be placed in one of two places:
Single fields: `_<submodule>_<Class>_<Field>`
Multiple fields: `_<submodule>_<Class>` of type `<submodule><Class>Private` immediately after the type definition of the main struct with package-private fields.

AS GO IS A FULL-RECOMPILATION LANGUAGE, IT IS NOT CONSIDERED A BREAKING CHANGE TO ALTER THESE SECTIONS.

Functions & Types that are private to the file: `<submodule><Type>`

If a struct has a `New<Struct>()` function with it, ALWAYS use functions to construct it!
However, if it merely has `<Struct>FromXYZ` or such, these are simply helpers.

All UI elements (implements and/or is UIElement) must start with `UI`. Things that are not UI elements must not.

Sometimes a struct is included inside another struct, and only the outer struct should have access to something.
In this case, create an inner structure.
This should only need to happen if the structure is intended to be inherited.

Reserved names for this are here:
```
UIThis : UILayoutElementComponentDetails
UIPanelDetails : UIPanelDetails
```

For the SDL2 backend, functions called within the OS-thread must start with `os_`.
