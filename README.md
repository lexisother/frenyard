# Frenyard

A GUI library for Golang originally built for CCUpdaterUI (a mod manager for CrossCode), and now separated out into a separate repo.

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` _will always pull the latest files from the master branch._

```sh
go get github.com/uwu/frenyard
```

### Usage

Import the package into your project.

```go
import "github.com/uwu/frenyard"
```

See Documentation and Examples below for more detailed information.

## Documentation

**NOTICE**: This library can be rather confusing.
Because of that it may be difficult to get into using the library.

The code is fairly well documented at this point and is currently
the only documentation available. Go reference (below) presents that information in a nice format.

- [![Go Reference](https://pkg.go.dev/badge/github.com/uwu/frenyard.svg)](https://pkg.go.dev/github.com/uwu/frenyard)
- Hand crafted documentation coming eventually.

## Examples

Below is a list of examples and other projects using Frenyard. Please submit
an issue if you would like your project added or removed from this list.

- [Frenyard Examples](https://github.com/uwu/frenyard/tree/master/examples) - A collection of example programs written using Frenyard
- [Cumcord Installer](https://github.com/Cumcord/Impregnate) - Installer for Cumcord, a discontinued and no-longer-functioning modification for Discord.
- [GitFren](https://github.com/lexisother/GitFren) - A program that was supposed to be a local frontend for Gitea.
- [CCUpdaterUI](https://github.com/dmitmel/CCUpdaterUI) - A mod manager and modloader installer for the video game CrossCode.

## Contributing

Contributions are very welcomed, however please follow the below guidelines.

- First open an issue describing the bug or enhancement so it can be
  discussed.
- Try to match current naming conventions as closely as possible.
- Create a Pull Request with your changes against the master branch.
