# VN Gateway

This project is a gateway for the VN project. It is responsible for handling all the incoming requests and routing them
to the appropriate services. It also handles the authentication and authorization of the requests. The gateway is built
on top of the [Go](https://golang.org/) programming language and uses the [Fiber](https://gofiber.io/) framework.

## Project structure

The project follows the standard Go project layout. The structure of the project is as follows:

```
.
├── .github                     # GitHub actions workflows
├── .golangci.yml
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── gateway
│       └── main.go
├── docs
├── go.mod
├── go.sum
├── internal
│   ├── config
│   ├── middleware
│   ├── models
│   ├── routes
│   ├── service
│   ├── tools
│   └── utils
├── scripts
└── test
```

## Prerequisites installation

- [x] Install [Node.js (v.20.13.1+)](https://nodejs.org/en/download/) or via `nvm`
- [x] Install [Go 1.22.3+](https://golang.org/doc/install)
- [x] Install [Docker Desktop](https://docs.docker.com/get-docker/)
- [x] (For Windows users) Install [WSL2](https://docs.microsoft.com/en-us/windows/wsl/install)
- [x] (For Windows users) Install [Chocolatey](https://chocolatey.org/install) and
  then run `choco install make` to install `make` command
- [x] Install [golangci-lint](https://golangci-lint.run/welcome/install/) and
    - Set up IDE integration (see [instructions](https://golangci-lint.run/welcome/integrations/)). This is optional
      because it may cause performance issues in some IDEs. You are still can run `make check` to lint your code
      instead.

## Prepare the environment

1. Run `npm run install` to install the project dependencies

## ⚠️ Pre-commit ⚠️

Make sure you have already run `make lint` before committing your code. This will ensure that your code is
properly formatted and passes all the tests.