# go-office365
`go-office365` provides both a client library as well as a CLI tool to interact with the `Microsoft Office365 Management Activity API`.

## Table of Contents

- [Overview](#overview)
- [Configuration file](#configuration-file)
- [Usage](#usage)
  - [Commands](#commands)
    - [Audit](#audit)
    - [Content](#content)
    - [Content-types](#content-types)
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

### Commands
#### audit
```
Retrieve events and/or actions for the provided audit-id.

Usage:
  go-office365 audit [audit-id] [flags]

Flags:
  -h, --help   help for audit

Global Flags:
      --config string   config file
```
#### content
```
List content that is available to be fetched for the provided content-type.

Usage:
  go-office365 content [content-type] [flags]

Flags:
      --end string     End time
  -h, --help           help for content
      --start string   Start time

Global Flags:
      --config string   config file
```
#### content-types
```
List content types accepted by the Microsoft API.

Usage:
  go-office365 content-types [flags]

Flags:
  -h, --help   help for content-types

Global Flags:
      --config string   config file
```
#### fetch
```
Combination of content and audit commands.

Usage:
  go-office365 fetch [content-type] [flags]

Flags:
      --end string     End time
  -h, --help           help for fetch
      --start string   Start time

Global Flags:
      --config string   config file
```
#### start-sub
```
Start a subscription for the provided Content Type.

Usage:
  go-office365 start-sub [content-type] [flags]

Flags:
  -h, --help   help for start-sub

Global Flags:
      --config string   config file
```
#### stop-sub
```
Stop a subscription for the provided Content Type.

Usage:
  go-office365 stop-sub [content-type] [flags]

Flags:
  -h, --help   help for stop-sub

Global Flags:
      --config string   config file
```
#### subscriptions
```
List current subscriptions.

Usage:
  go-office365 subscriptions [flags]

Flags:
  -h, --help   help for subscriptions

Global Flags:
      --config string   config file
```
#### watch
```
Fetch audit events at regular intervals.

Usage:
  go-office365 watch [flags]

Flags:
  -h, --help               help for watch
      --human-readable     Human readable output format.
      --interval int       TickerIntervalSeconds (default 5)
      --lookbehind int     Number of minutes from request time used when fetching
                           available content. (default 1)
      --output string      Target where to send audit records.
                           Available schemes:
                           file://path/to/file
                           udp://1.2.3.4:1234
                           tcp://1.2.3.4:1234
      --statefile string   File used to read/save state on start/exit.

Global Flags:
      --config string   config file
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
