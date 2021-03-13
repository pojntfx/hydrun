# hydrun ("Hydra Run")

Execute a command for the current directory on multiple architectures and operating systems.

![hydrun CI](https://github.com/pojntfx/hydrun/workflows/hydrun%20CI/badge.svg)

## Overview

Hydra Run, or hydrun, is a thin (<200 SLOC) layer atop [Docker buildx](https://github.com/docker/buildx) and [qemu-user-static](https://ngithub.com/multiarch/qemu-user-static). It allows one to easily execute a command on different processor architectures and operating systems than the host.

It can, for example, be used for ...

- **Cross-compilation that "just works"**, without having to set up a cross-compiler (at the cost of longer build times)
- **Multi-architecture testing**
- **Building arm64 binaries on GitHub actions**, which doesn't support arm64 runners or Linux distros other than Ubuntu
- Quickly getting an **interactive arm64 shell for the current directory on an amd64 host** or the other way round
- Running binaries built against **glibc on an Alpine Linux host**
- Making **CI release builds locally reproducable and testable**, without having to `git push` and wait

## Installation

Binaries are available on [GitHub releases](https://github.com/pojntfx/hydrun/releases).

You can install them like so:

```shell
$ curl -L -o /tmp/hydrun https://github.com/pojntfx/hydrun/releases/download/latest/hydrun.linux-$(uname -m)
$ sudo install /tmp/hydrun /usr/local/bin
```

## Usage

TODO: Add quick usage guide

## Reference

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
