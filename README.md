# go-office365
`go-office365` provides both a client library as well as a CLI tool to interact with the `Microsoft Office365 Management Activity API`.

## Table of Contents

- [Overview](#overview)
- [Configuration file](#configuration-file)
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

## License
go-office365 is released under the MIT license. See [LICENSE.txt](LICENSE.txt)
