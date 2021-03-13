# hydrun

Run a command for the current directory on multiple processor architectures and operating systems.

![hydrun CI](https://github.com/pojntfx/hydrun/workflows/hydrun%20CI/badge.svg)

## Installation

Binaries are available on [GitHub releases](https://github.com/pojntfx/hydrun/releases).

## Usage

```shell
$ hydrun --help
Run a command for the current directory on multiple processor architectures and operating systems.

Usage: hydrun [options...] <commands...>
  -a, --arch string   Processor architecture(s) to run on. Separate multiple values with commas. (default "amd64,arm64v8")
  -o, --os string     Operating system(s) to run on. Separate multiple values with commas. (default "debian")

See https://github.com/pojntfx/hydrun for more information.
```

## License

hydrun (c) 2021 Felix Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
