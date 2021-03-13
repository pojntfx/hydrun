# hydrun ("Hydra Run")

Execute a command for the current directory on multiple architectures and operating systems.

![hydrun CI](https://github.com/pojntfx/hydrun/workflows/hydrun%20CI/badge.svg)

## Installation

Binaries are available on [GitHub releases](https://github.com/pojntfx/hydrun/releases).

You can install them like so:

```shell
$ curl -L -o /tmp/hydrun https://github.com/pojntfx/hydrun/releases/download/latest/hydrun.linux-$(uname -m)
$ sudo install /tmp/hydrun /usr/local/bin
```

## Usage

```shell
$ hydrun --help
Execute a command for the current directory on multiple architectures and operating systems.

See https://github.com/pojntfx/hydrun for more information.

Usage: hydrun [OPTION...] "<COMMAND...>"
  -a, --arch string   Comma-separated list of architectures to run on (default "amd64")
  -i, --it            Attach stdin and setup a TTY
  -j, --jobs int      Maximum amount of parallel jobs (default 1)
  -o, --os string     Comma-separated list of operating systems to run on (default "debian")
```

## License

hydrun (c) 2021 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
