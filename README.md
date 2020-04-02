# go-office365
`go-office365` provides both a client library as well as a CLI tool to interact with the `Microsoft Office365 Management Activity API`.

## Table of Contents

- [Overview](#overview)
- [Configuration file](#configuration-file)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [License](#license)

## Overview
`go-office365` provides a client library for the `Microsoft Office365 Management Activity API` written in [Go](https://golang.org/). It follows the Microsoft API Reference available [here](https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference).

`go-office365` is also a CLI application with everything you need to interact with the API remotely.

## Get Started
`go-office365` uses Go Modules introduced in Go 1.11 for dependency management.

### Build
Build the CLI for a target platform (Go cross-compiling feature), for example linux, by executing:
```
$ mkdir $HOME/src
$ cd $HOME/src
$ git clone https://github.com/devodev/go-office365.git
$ cd go-office365
$ env GOOS=linux go build -o go_office365_linux ./go-office365
```
If you are a Windows user, substitute the $HOME environment variable above with %USERPROFILE%.

## Configuration file
> An alternative location can be provided using the `--config` flag.

The CLI will look into the root directory of the executable using the name `.go-office365.yaml`.

The credential attributes can be found in `Azure Active Directory`, under: `Installed apps`.</br>
The Identifier attribute is provided on all queries to the Microsoft API as the `PublisherIdentifier` query param and is(was?) used to compute quotas. When empty, the param is not sent.

```
// .go-office365.yaml
---
Global:
  Identifier: some-id
Credentials:
  ClientID: 00000000-0000-0000-0000-000000000000
  ClientSecret: 00000000000000000000000000000000
  TenantID: 00000000-0000-0000-0000-000000000000
  TenantDomain: some-company.onmicrosoft.com
```

## Usage
Auto-generated documentation of commands can be found [here](./go-office365/docs/go-office365.md).
```
Query the Microsoft Office365 Management Activity API.

Usage:
  go-office365 [command]

Available Commands:
  audit         Retrieve events and/or actions for the provided audit-id.
  content       List content that is available to be fetched for the provided content-type.
  content-types List content types accepted by the Microsoft API.
  fetch         Combination of content and audit commands.
  help          Help about any command
  start-sub     Start a subscription for the provided Content Type.
  stop-sub      Stop a subscription for the provided Content Type.
  subscriptions List current subscriptions.
  watch         Fetch audit events at regular intervals.

Flags:
      --config string   config file
  -h, --help            help for go-office365
  -v, --version         version for go-office365

Use "go-office365 [command] --help" for more information about a command.
```

## Roadmap
Overall
- <strike>Better README</strike>
- More documentation
- More tests

<strike>Office365 library</strike>
- <strike>Add logger</strike>

CLI
- Add tests
- `Wont do:` <strike>Add flag to turn On/Off the verbosity of the logger</strike>
- Add flag to provide file to log to
- Start-sub command
  - Add option to provide a webhook object definition
    - `Maybe:` Use viper to load a file containing a webhook schema
- <strike>Watch command</strike>
  - <strike>Create JsonHandler</strike>
  - <strike>Let user choose which resource handler to use using command line arg</strike>
  - <strike>Let user provide additional config based on resource handler, like output target(network, file, etc.)</strike>
  - <strike>`Maybe:` Use a sub-sub-command to handle handler logic</strike>

## Contributing
> This is my first Go project, and also in conjunction, my first contribution to the open source community, so please feel free to open issues and discuss how you would improve the current code. I am eager to read you and learn from the community. Thanks!</br>`@devodev`

This repository is under heavy development and is subject to change in the near future.</br>
Versioning will be locked and a proper contributing section will be created in a timely manner, when code is stabilized.</br>

## License
`go-office365` is released under the MIT license. See [LICENSE.txt](LICENSE.txt)
