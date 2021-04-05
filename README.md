# gslauncher

Repository for the Groovestats Launcher for the Simply Love SM5 theme.


## Building on Linux

Building the binary is straight-forward. Just run go build to create the
executable (this assumes you already have git and golang installed).

```sh
git clone git@github.com:GrooveStats/gslauncher.git
cd gslauncher
go build ./cmd/gslauncher/
```

The binary will be called `gslauncher`. It can be distributed to other linux
systems, no dependencies required.

For development you can also take a shortcut for building the binary and
running int in one command.

```sh
go run ./cmd/gslauncher
```


## Building on Windows

Some setup is required, but otherwise building works the same as on Linux.

1. Install git (get the latest version from https://git-scm.com/).
2. Install golang (https://golang.org/doc/install. The website has a guide)
3. Install msys2 (https://www.msys2.org/#installation).
   Run `pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-pkg-config` in
   the msys2 shell. The mingw toolchain is necessary for building cgo
   dependencies.
4. Open git bash and clone the repository. Before running `go build` the PATH
   environment variable has to be changed to have access to the mingw
   toolchain. Run `export PATH="/c/msys64/mingw64/bin/:$PATH"`. You will have
   to do that whenever you open git bash.
5. You are good to go now. Refer to the previous section for general buliding
   instructions.


## Debug Build

The launcher can be built with the `debug` build tag. It adds additional
settings to the menu.

```sh
go run -tags debug ./cmd/gslauncher
```
