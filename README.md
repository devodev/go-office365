# go-office365
`go-office365` provides both a client library as well as a CLI tool to interact with the `Microsoft Office365 Management Activity API`.

## Table of Contents

- [Overview](#overview)
- [Configuration file](#configuration-file)
- [Usage](#usage)
- [Commands](#commands)
  - [Audit](#audit)
  - [Content](#content)
  - [Fetch](#fetch)
  - [Start-sub](#start-sub)
  - [Stop-sub](#stop-sub)
  - [Subscriptions](#subscriptions)
  - [Watch](#watch)
- [Roadmap](#roadmap)
- [License](#license)

## Overview
`go-office365` provides a client library for the `Microsoft Office365 Management Activity API` written in Go. It follows the Microsoft API Reference available [here](https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-management-activity-api-reference).

`go-office365` is also a CLI application with a ton of useful features.

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

## Commands
### audit
### content
### fetch
### start-sub
### stop-sub
### subscriptions
### watch

## Roadmap
Overall
- Better README
- More documentation
- More tests

Office365 library
- Add logger

CLI
- Add tests
- Add flag to turn On/Off the verbosity of the logger
- Add flag to provide file to log to
- Start-sub command
  - Add option to provide a webhook object definition
    - `Maybe:` Use viper to load a file containing a webhook schema
- Watch command
  - Create JsonHandler
  - Let user choose which resource handler to use using command line arg
  - Let user provide additional config based on resource handler, like output target(network, file, etc.)
  - `Maybe:` Use a sub-sub-command to handle handler logic

## Contributing
> This is my first Go project, and also in conjunction, my first contribution to the open source community, so please feel free to open issues and discuss how you would improve the current code. I am eager to read you and learn from the community. Thanks!</br>`@devodev`

This repository is under heavy development and is subject to change in the near future.</br>
Versioning will be locked and a proper contributing section will be created in a timely manner, when code is stabilized.</br>

## License
`go-office365` is released under the MIT license. See [LICENSE.txt](LICENSE.txt)
