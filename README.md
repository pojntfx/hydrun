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

Before continuing, please ensure that you have both [Docker buildx](https://github.com/docker/buildx) and [qemu-user-static](https://ngithub.com/multiarch/qemu-user-static) installed on the host.

### To Get an Interactive Shell

As described in the [Reference](#Reference), you can get an interactive shell by using the `-i` flag. The `-a` parameter corresponds to an [architecture](https://www.docker.com/blog/multi-platform-docker-builds/) such as `amd64`, `arm64` or `ppc64le`; the `-o` flag corresponds to a [Docker image](https://hub.docker.com/search?q=&type=image) such as `debian`, `alpine` or `fedora`. To for example run arm64 Debian on an amd64 host, you can do the following:

```shell
$ uname -a
Linux dev-tmp 4.19.0-10-cloud-amd64 #1 SMP Debian 4.19.132-1 (2020-07-24) x86_64 GNU/Linux
$ hydrun -a arm64 -o debian -i "bash"
root@81647bd6aa02:/data# uname -a
Linux 81647bd6aa02 4.19.0-10-cloud-amd64 #1 SMP Debian 4.19.132-1 (2020-07-24) aarch64 GNU/Linux
root@81647bd6aa02:/data# ldd $(which ls)
        libselinux.so.1 => /lib/aarch64-linux-gnu/libselinux.so.1 (0x0000005501868000)
        libc.so.6 => /lib/aarch64-linux-gnu/libc.so.6 (0x000000550189e000)
        /lib/ld-linux-aarch64.so.1 (0x0000005500000000)
        libpcre.so.3 => /lib/aarch64-linux-gnu/libpcre.so.3 (0x0000005501a10000)
        libdl.so.2 => /lib/aarch64-linux-gnu/libdl.so.2 (0x0000005501a83000)
        libpthread.so.0 => /lib/aarch64-linux-gnu/libpthread.so.0 (0x0000005501a97000)
```

### As an Alternative to Cross-Compilation

It is very easy to use hydrun to get binaries for many platforms, without having to set up cross-compilation. Consider the following C hello world program:

```c
/* main.c */
#include <stdio.h>

int main()
{
  printf("Hello, world!\n");

  return 0;
}
```

Using hydrun, we can now compile it for multiple architectures:

```shell
$ ls
ls
main.c
$ hydrun -a amd64,arm64 -o gcc 'gcc -static -o hello-world.linux-$(uname -m) main.c'
$ ls
hello-world.linux-aarch64  hello-world.linux-x86_64  main.c
$ file *
hello-world.linux-aarch64: ELF 64-bit LSB executable, ARM aarch64, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.7.0, with debug_info, not stripped
hello-world.linux-x86_64:  ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, for GNU/Linux 3.2.0, with debug_info, not stripped
main.c:                    C source, ASCII text
```

It is also possible to run/test the compiled binaries with it:

```shell
$ ls
hello-world.linux-aarch64  hello-world.linux-x86_64  main.c
$ hydrun -i -a arm64 "./hello-world.linux-aarch64"
Hello, world!
$ hydrun -i -a amd64 "./hello-world.linux-x86_64"
Hello, world!
```

When building larger projects or requiring dependencies, it is recommended to put the commands used into a shell script, conventionally named the `Hydrunfile`:

```
$ cat <<EOT>Hydrunfile
#!/bin/bash
apt update
apt install -y build-essentials

gcc -static -o hello-world.linux-$(uname -m) main.c
EOT
$ chmod +x ./Hydrunfile
```

You can now easily build like so:

```shell
$ hydrun -a amd64,arm64 "./Hydrunfile"
$ ls
hello-world.linux-aarch64  hello-world.linux-x86_64  Hydrunfile  main.c
```

Most of the time you'll probably want to use the `Hydrunfile` to install toolchains/dependencies and then call your `Makefile`, the Go compiler, `cargo` etc.

For an example, check out [pojntfx/liwasc](https://github.com/pojntfx/liwasc).

### Usage in GitHub Actions

It is also possible to use hydrun to build multi-architecture binaries in a CI/CD system such as GitHub actions. Continuing with the C example from above, we could automatically build and release amd64 and arm64 binaries to GitHub releases using GitHub actions with the following workflow:

```yaml
name: hydrun CI

on:
  push:
  pull_request:

jobs:
  build-linux:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v1
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v1
      - name: Set up hydrun
        run: |
          curl -L -o /tmp/hydrun https://github.com/pojntfx/hydrun/releases/download/latest/hydrun.linux-$(uname -m)
          sudo install /tmp/hydrun /usr/local/bin
      - name: Build with hydrun
        run: hydrun -a amd64,arm64 ./Hydrunfile
      - name: Publish to GitHub releases
        if: ${{ github.ref == 'refs/heads/main' }}
        uses: marvinpinto/action-automatic-releases@latest
        with:
          repo_token: "${{ secrets.GITHUB_TOKEN }}"
          automatic_release_tag: "latest"
          prerelease: false
          files: |
            *.linux*
```

For an example, check out [pojntfx/liwasc](https://github.com/pojntfx/liwasc).

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
