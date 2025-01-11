# go-fltk-investment-balancer

Balances your portfolio.

Features dark/light mode and portrait/landscape mode that is responsive.

The UI uses FLTK for extremely minimal memory usage.

**This application is not finished yet, but works for the most part.** This readme is incomplete.

## Screenshots

Coming soon

<!-- ![Light mode landscape](./docs/light-landscape.png) -->

<!-- ![Dark mode portrait](./docs/dark-portrait.png) -->

## Installation

### Prerequisites

On Fedora:

```bash
sudo dnf install go git-lfs fltk fltk-devel mesa-libGL-devel mesa-libGLU-devel
```

On Arch Linux:

```bash
# TODO: verify this on a fresh install
sudo pacman -S go git-lfs fltk
```

### Build

```bash
go get -v
go build -v
```

### Run

```bash
./go-fltk-investment-balancer
```

## Development setup

This repository makes use of `git lfs`. Please ensure you have it working.

To build, run

```bash
make build-prod
```

To install to `~/.local/bin/`, run

```bash
make install
```
