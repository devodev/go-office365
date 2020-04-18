<p align="center">
  <a href="https://github.com/devodev/go-office365">
    <img alt="go-office365" src="assets/go-office365-logo.png" width="500">
  </a>
</p>
<p align="center">
  A CLI as well as a library to interact with the Microsoft Office365 Management Activity API.
</p>
<p align="center">
    <a href="https://github.com/golang/go/wiki/Modules">
        <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/devodev/go-office365">
    </a>
    <a href="https://pkg.go.dev/mod/github.com/devodev/go-office365">
        <img alt="go.dev reference" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
    <a href="https://goreportcard.com/report/github.com/devodev/go-office365">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/devodev/go-office365">
    </a>
</p>
<p align="center">
    <a href="https://github.com/devodev/go-office365/releases">
        <img alt="Release" src="https://github.com/devodev/go-office365/workflows/Release/badge.svg">
    </a>
    <a href="https://github.com/devodev/go-office365/tags">
        <img alt="GitHub tag (latest SemVer)" src="https://img.shields.io/github/v/tag/devodev/go-office365?sort=semver">
    </a>
    <a href="https://github.com/devodev/go-office365/blob/master/LICENSE.txt">
        <img alt="GitHub license" src="https://img.shields.io/github/license/devodev/go-office365?style=flat">
    </a>
</p>

## Overview
`go-office365` provides a client library for the `Microsoft Office365 Management Activity API` written in [Go](https://golang.org/). It follows the Microsoft API Reference available [here](https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference).

`go-office365` is also a CLI application with everything you need to interact with the API on the command line.

Currently, **`go-office365` requires Go version 1.13 or greater**.

#### Supported Architectures
We provide pre-built go-office365 binaries for Windows, Linux and macOS (Darwin) architectures, in both 386/amd64 flavors.</br>
Please see the release section [here](https://github.com/devodev/go-office365/releases).

## Table of Contents

- [Overview](#overview)
  - [Supported Architectures](#supported-architectures)
- [Get Started](#get-started)
  - [Build](#build)
- [CLI](#cli)
  - [Usage](#usage)
  - [Configuration file](#configuration-file)
  - [Interval flags](#interval-flags)
  - [Watcher](#watcher)
    - [How it works](#how-it-works)
  - [Extended Schemas](#extended-schemas)
- [Roadmap](#roadmap)
  - [CLI Commands](#cli-commands)
- [Contributing](#contributing)
- [License](#license)

## Get Started
`go-office365` uses Go Modules introduced in Go 1.11 for dependency management.

### Build
Build the CLI for a target platform (Go cross-compiling feature), for example linux, by executing:
```
$ mkdir $HOME/src
$ cd $HOME/src
$ git clone https://github.com/devodev/go-office365.git
$ cd go-office365
$ env GOOS=linux go build -o go_office365_linux ./cmd/go-office365
```
If you are a Windows user, substitute the $HOME environment variable above with %USERPROFILE%.

## CLI
### Usage
> Auto-generated documentation for each command can be found [here](./docs/go-office365.md).
```
Interact with the Microsoft Office365 Management Activity API.

Usage:
  go-office365 [command]

Available Commands:
  audit         Query audit records for the provided audit-id.
  content       Query content for the provided content-type.
  content-types List content types accepted by the Microsoft API.
  fetch         Query audit records for the provided content-type.
  gendoc        Generate markdown documentation for the go-office365 CLI.
  help          Help about any command
  start-sub     Start a subscription for the provided Content Type.
  stop-sub      Stop a subscription for the provided Content Type.
  subscriptions List current subscriptions.
  watch         Query audit records at regular intervals.

Flags:
  -h, --help   help for go-office365

Use "go-office365 [command] --help" for more information about a command.
```

### Configuration file
Commands that need to interact with the API require credentials to be provided using a YAML configuration file.</br>
The following locations are looked into if the --config flag is not provided:
```
$HOME/.go-office365.yaml
$CWD/.go-office365.yaml
```

The following is the current schema used.

> Identifier is provided on all queries to the Microsoft API as the `PublisherIdentifier` query param and is(was?) used to compute quotas. When empty, the param is not sent.

>Credentials can be found in `Azure Active Directory`, under: `Installed apps`.</br>

```
---
Global:
  Identifier: some-id
Credentials:
  ClientID: 00000000-0000-0000-0000-000000000000
  ClientSecret: 00000000000000000000000000000000
  TenantID: 00000000-0000-0000-0000-000000000000
  TenantDomain: some-company.onmicrosoft.com
```

### Interval flags
Commands that need to use a fixed interval will offer flags to set the start and end times.</br>
Here are the guidelines to follow when providing those flags.

```
- Both or neither of start/end time must be provided.
- When not provided, a 24 hour interval is used.
- Start and end time interval must be between 1 minute and 24 hours.
- Start time must not be earlier than 7 days behind the current time.
- Time format must match one of: 2006-01-02, 2006-01-02T15:04, 2006-01-02T15:04:05
```

### Watcher
The `watch` command provides a daemon like process for retrieving audit records at regular intervals.</br>
It uses a minimum amount of resources and a lot of useful flags can be provided.</br>
> For more details on what flags can be used, see the command documentation [here](./docs/go-office365_watch.md).

#### How it works
- Upon starting, a data structure is initialized to retain the last request time and the last content creation time. A statefile location can be provided for persisting state between restarts.</br>
- Following, a resource handler is spawned. It is responsible for receiving, formatting and sending records to the selected output.</br>
- Then, for each each Microsoft content type, a data pipeline is spawned. When triggered, it will query and relay audit records to the resource handler.</br>
- At fixed intervals, a subscription worker is spawned. It will query the content subscriptions currently enabled and will trigger the appropriate data pipelines.</br>

### Extended Schemas
By default, audit events are retrieved and stored using the AuditRecord type. An option is available to
add remaining fields, when present, depending on the RecordType provided in the Record.</br>
Whenever an extended schema assigned to a RecordType fails to parse the remaining fields, the base AuditRecord is returned.

## Roadmap
### CLI Commands
- `start-sub`: Add flag to provide a webhook object definition

## Contributing
> This is my first contribution to the open source community, so please feel free to open issues and discuss how you would improve the current code. I am eager to read you and learn from the community. Thanks!</br>`@devodev`

This repository is under heavy development and is subject to change in the near future.</br>
Versioning will be locked and a proper contributing section will be created in a timely manner, when code is stabilized.</br>

## License
`go-office365` is released under the MIT license. See [LICENSE.txt](LICENSE.txt)
