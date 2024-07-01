# VN Workday - Config

This is a Go project that be widely used in the VN Workday ecosystem. It provides a simple way to manage
configuration files in Go projects.

## Prerequisites installation

- [x] Install [Node.js (v.20.13.1+)](https://nodejs.org/en/download/) or via `nvm`
- [x] Install [Go 1.22.3+](https://golang.org/doc/install)
- [x] (For Windows users) Install [WSL2](https://docs.microsoft.com/en-us/windows/wsl/install)
- [x] (For Windows users) Install [Chocolatey](https://chocolatey.org/install) and
  then run `choco install make` to install `make` command
- [x] Install [golangci-lint](https://golangci-lint.run/welcome/install/) and
    - Set up IDE integration (see [instructions](https://golangci-lint.run/welcome/integrations/)). This is optional
      because it may cause performance issues in some IDEs. You are still can run `make check` to lint your code
      instead.

## Prepare the environment

1. Run `npm run install` to install the project dependencies
2. Run `npm run prepare` to make sure commit hooks are installed

## ⚠️ Pre-commit ⚠️

Make sure you have already run `make lint` before committing your code. This will ensure that your code is
properly formatted and linted.