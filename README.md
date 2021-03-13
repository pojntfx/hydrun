# hydrun

Execute a command for multiple architectures and operating systems.

![hydrun CI](https://github.com/pojntfx/hydrun/workflows/hydrun%20CI/badge.svg)

## Installation

Binaries are available on [GitHub releases](https://github.com/pojntfx/hydrun/releases).

## Usage

```shell
$ hydrun --help
Execute a command for multiple architectures and operating systems.

See https://github.com/pojntfx/hydrun for more information.

Usage: hydrun [OPTION...] "<COMMAND...>"
  -a, --arch string   Comma-separated list of architectures to run on (default "amd64")
  -i, --it            Attach stdin and setup a TTY
  -j, --jobs int      Maximum amount of parallel jobs (default 1)
  -o, --os string     Comma-separated list of operating systems to run on (default "debian")
```

## License

hydrun (c) 2021 Felix Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
