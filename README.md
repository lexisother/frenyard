# Frenyard

A GUI library for Golang originally built for CCUpdaterUI (a mod manager for CrossCode), and now separated out into a separate repo.

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` _will always pull the latest files from the master branch._

```sh
go get github.com/lexisother/frenyard
```

### Usage

Import the package into your project.

```go
import "github.com/lexisother/frenyard"
```

See Documentation and Examples below for more detailed information.

## Documentation

**NOTICE**: This library can be rather confusing.
Because of that it may be difficult to get into using the library.

The code is fairly well documented at this point and is currently
the only documentation available. Go reference (below) presents that information in a nice format.

- [![Go Reference](https://pkg.go.dev/badge/github.com/lexisother/frenyard.svg)](https://pkg.go.dev/github.com/lexisother/frenyard)
- Hand crafted documentation coming eventually.

## Environment variables

Frenyard looks for certain environment variables as input to some internal functions.

### `FRENYARD_SCALE`

A float that directly decides the scale of the window, no calculation is done.

### `FRENYARD_EXPR_MACOS_FIX`

Enables some experimental fixes regarding scaling bugs on macOS. Can be used in tandem with `FRENYARD_SCALE` to debug
said scaling issues.

## Examples

Below is a list of examples and other projects using Frenyard. Please submit
an issue if you would like your project added or removed from this list.

- [Frenyard Examples](https://github.com/lexisother/frenyard/tree/master/examples) - A collection of example programs written using Frenyard
- [Cumcord Installer](https://github.com/Cumcord/Impregnate) - Installer for Cumcord, a discontinued and no-longer-functioning modification for Discord.
- [Replugged Installer](https://github.com/replugged-org/installer) - Installer for Replugged, a functioning modification for Discord.
- [GitFren](https://github.com/lexisother/GitFren) - A program that was supposed to be a local frontend for Gitea.
- [CCUpdaterUI](https://github.com/dmitmel/CCUpdaterUI) - A mod manager and modloader installer for the video game CrossCode.

## Contributing

Contributions are very welcomed, however please follow the below guidelines.

- First open an issue describing the bug or enhancement so it can be
  discussed.
- Try to match current naming conventions as closely as possible.
- Create a Pull Request with your changes against the master branch.
